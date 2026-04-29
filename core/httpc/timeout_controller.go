package httpc

import (
	"bytes"
	"context"
	"errors"
	"maps"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
)

var ErrHandlerTimeout = errors.New("httpc: handler timeout")

const (
	timeoutWheelTick  = 10 * time.Millisecond
	timeoutWheelSlots = 512
)

var (
	timeoutControllerOnce sync.Once
	timeoutControllerInst *timeoutController
)

type timeoutController struct {
	workerPool *ants.PoolWithFuncGeneric[*timeoutCall]
	wheel      *timeoutWheel
	callPool   sync.Pool
	writerPool sync.Pool
}

type timeoutCall struct {
	controller *timeoutController
	writer     *timeoutWriter
	handler    http.HandlerFunc
	request    *http.Request
	cancel     context.CancelCauseFunc
	timer      *wheelTask
	doneCh     chan struct{}
	panicValue any
	refCount   atomic.Int32
	completed  atomic.Bool
	timedOut   atomic.Bool
}

type timeoutWriter struct {
	base        http.ResponseWriter
	header      http.Header
	buffer      bytes.Buffer
	mu          sync.Mutex
	err         error
	wroteHeader bool
	statusCode  int
}

func defaultTimeoutController() *timeoutController {
	timeoutControllerOnce.Do(func() {
		timeoutControllerInst = newTimeoutController()
	})
	return timeoutControllerInst
}

func newTimeoutController() *timeoutController {
	controller := &timeoutController{}
	controller.callPool.New = func() any {
		return &timeoutCall{
			controller: controller,
			doneCh:     make(chan struct{}, 1),
		}
	}
	controller.writerPool.New = func() any {
		return &timeoutWriter{
			header: make(http.Header, 8),
		}
	}

	pool, err := ants.NewPoolWithFuncGeneric[*timeoutCall](-1, func(call *timeoutCall) {
		call.run()
	})
	if err != nil {
		panic(err)
	}

	controller.workerPool = pool
	controller.wheel = newTimeoutWheel(timeoutWheelTick, timeoutWheelSlots)
	return controller
}

func (c *timeoutController) Wrap(timeout time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancelCause(r.Context())
		defer cancel(nil)

		call := c.acquireCall(w, r.WithContext(ctx), next, cancel)
		call.timer = c.wheel.Schedule(timeout, call.onTimeout)

		if err := c.workerPool.Invoke(call); err != nil {
			go call.run()
		}

		select {
		case <-call.doneCh:
			call.stopTimer()
			if call.timedOut.Load() {
				writeTimeoutResponse(w)
				call.releaseRef()
				return
			}
			if call.panicValue != nil {
				call.releaseRef()
				panic(call.panicValue)
			}
			call.writer.commitTo(w)
			call.releaseRef()
		case <-ctx.Done():
			call.stopTimer()
			cause := context.Cause(ctx)
			if call.timedOut.Load() || errors.Is(cause, context.DeadlineExceeded) {
				writeTimeoutResponse(w)
			} else {
				call.writer.fail(cause)
				w.WriteHeader(http.StatusServiceUnavailable)
			}
			call.releaseRef()
		}
	}
}

func (c *timeoutController) acquireCall(w http.ResponseWriter, r *http.Request, handler http.HandlerFunc, cancel context.CancelCauseFunc) *timeoutCall {
	call := c.callPool.Get().(*timeoutCall)
	call.writer = c.acquireWriter(w)
	call.handler = handler
	call.request = r
	call.cancel = cancel
	call.timer = nil
	call.panicValue = nil
	call.refCount.Store(2)
	call.completed.Store(false)
	call.timedOut.Store(false)

	select {
	case <-call.doneCh:
	default:
	}

	return call
}

func (c *timeoutController) acquireWriter(w http.ResponseWriter) *timeoutWriter {
	writer := c.writerPool.Get().(*timeoutWriter)
	writer.reset(w)
	return writer
}

func (c *timeoutController) releaseWriter(writer *timeoutWriter) {
	writer.reset(nil)
	c.writerPool.Put(writer)
}

func (c *timeoutController) recycleCall(call *timeoutCall) {
	if call.writer != nil {
		c.releaseWriter(call.writer)
	}

	call.writer = nil
	call.handler = nil
	call.request = nil
	call.cancel = nil
	call.timer = nil
	call.panicValue = nil
	call.completed.Store(false)
	call.timedOut.Store(false)

	select {
	case <-call.doneCh:
	default:
	}

	c.callPool.Put(call)
}

func (c *timeoutCall) run() {
	defer c.releaseRef()
	defer func() {
		c.completed.Store(true)
		if p := recover(); p != nil {
			c.panicValue = p
		}
		select {
		case c.doneCh <- struct{}{}:
		default:
		}
	}()

	c.handler(c.writer, c.request)
}

func (c *timeoutCall) onTimeout() {
	if c.completed.Load() {
		return
	}
	if !c.timedOut.CompareAndSwap(false, true) {
		return
	}
	c.writer.fail(ErrHandlerTimeout)
	c.cancel(context.DeadlineExceeded)
}

func (c *timeoutCall) stopTimer() {
	if c.timer == nil {
		return
	}
	if !c.timer.Stop() {
		c.timer.Release()
	}
	c.timer = nil
}

func (c *timeoutCall) releaseRef() {
	if c.refCount.Add(-1) != 0 {
		return
	}
	c.controller.recycleCall(c)
}

func (tw *timeoutWriter) reset(w http.ResponseWriter) {
	tw.base = w
	tw.err = nil
	tw.wroteHeader = false
	tw.statusCode = 0
	tw.buffer.Reset()
	if tw.header == nil {
		tw.header = make(http.Header, 8)
		return
	}
	clear(tw.header)
}

func (tw *timeoutWriter) Header() http.Header {
	return tw.header
}

func (tw *timeoutWriter) Write(p []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.err != nil {
		return 0, tw.err
	}
	if !tw.wroteHeader {
		tw.wroteHeader = true
		tw.statusCode = http.StatusOK
	}

	return tw.buffer.Write(p)
}

func (tw *timeoutWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.err != nil || tw.wroteHeader {
		return
	}
	tw.wroteHeader = true
	tw.statusCode = code
}

func (tw *timeoutWriter) Push(target string, opts *http.PushOptions) error {
	pusher, ok := tw.base.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}
	return pusher.Push(target, opts)
}

func (tw *timeoutWriter) fail(err error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.err == nil {
		tw.err = err
	}
}

func (tw *timeoutWriter) commitTo(w http.ResponseWriter) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	maps.Copy(w.Header(), tw.header)

	statusCode := tw.statusCode
	if !tw.wroteHeader {
		statusCode = http.StatusOK
	}

	w.WriteHeader(statusCode)
	if tw.buffer.Len() == 0 {
		return
	}
	_, _ = w.Write(tw.buffer.Bytes())
}

func writeTimeoutResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusRequestTimeout)
}

var _ http.Pusher = (*timeoutWriter)(nil)
