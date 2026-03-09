package tray

import (
	"io"
	"strings"
)

type ServiceTray struct {
	routers         []*strings.Builder
	apis            []*strings.Builder
	methods         []*strings.Builder
	imports         []*strings.Builder
	routerPath      string
	registerHttpApi bool
}

func (sc *ServiceTray) GetExtraImport() []string {
	is := make([]string, len(sc.imports))
	for i, m := range sc.imports {
		is[i] = m.String()
	}
	return is
}

func (sc *ServiceTray) GetApis() []string {
	ms := make([]string, len(sc.apis))
	for i, m := range sc.apis {
		ms[i] = m.String()
	}
	return ms
}

func (sc *ServiceTray) GetMethods() []string {
	ms := make([]string, len(sc.methods))
	for i, m := range sc.methods {
		ms[i] = m.String()
	}
	return ms
}

func (sc *ServiceTray) GetRouters() []string {
	rs := make([]string, len(sc.routers))
	for i, m := range sc.routers {
		rs[i] = m.String()
	}
	return rs
}

func (sc *ServiceTray) BorrowApisWriter() io.Writer {
	sc.apis = append(sc.apis, new(strings.Builder))
	return sc.apis[len(sc.apis)-1]
}

func (sc *ServiceTray) BorrowMethodsWriter() io.Writer {
	sc.methods = append(sc.methods, new(strings.Builder))
	return sc.methods[len(sc.methods)-1]
}

func (sc *ServiceTray) BorrowRouterssWriter() io.Writer {
	sc.routers = append(sc.routers, new(strings.Builder))
	return sc.routers[len(sc.routers)-1]
}

func (sc *ServiceTray) BorrowExtraImportReader() io.Writer {
	sc.imports = append(sc.imports, new(strings.Builder))
	return sc.imports[len(sc.imports)-1]
}

func (sc *ServiceTray) SetRouterBasePath(path string) {
	sc.routerPath = path
}

func (sc *ServiceTray) GetRouterBasePath() string {
	return sc.routerPath
}

func (sc *ServiceTray) IsRegisterHttpApi() bool {
	return sc.registerHttpApi
}

func (sc *ServiceTray) SetRegisterHttpApi(b bool) {
	sc.registerHttpApi = b
}

func NewServiceTray() *ServiceTray {
	return &ServiceTray{}
}
