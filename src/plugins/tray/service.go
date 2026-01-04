package tray

import (
	"io"
	"strings"
)

type ServiceTray struct {
	register        []*strings.Builder
	methods         []*strings.Builder
	tiers           []*strings.Builder
	imports         []*strings.Builder
	routerPath      string
	registerHttpApi bool
}

func (sc *ServiceTray) GetMethods() []string {
	ms := make([]string, len(sc.methods))
	for i, m := range sc.methods {
		ms[i] = m.String()
	}
	return ms
}

func (sc *ServiceTray) GetExtraImport() []string {
	is := make([]string, len(sc.imports))
	for i, m := range sc.imports {
		is[i] = m.String()
	}
	return is
}

func (sc *ServiceTray) GetTiers() []string {
	ts := make([]string, len(sc.tiers))
	for i, t := range sc.tiers {
		ts[i] = t.String()
	}
	return ts
}

func (sc *ServiceTray) GetRegisters() []string {
	ms := make([]string, len(sc.register))
	for i, m := range sc.register {
		ms[i] = m.String()
	}
	return ms
}

func (sc *ServiceTray) BorrowMethodsWriter() io.Writer {
	sc.methods = append(sc.methods, new(strings.Builder))
	return sc.methods[len(sc.methods)-1]
}

func (sc *ServiceTray) BorrowRegistersWriter() io.Writer {
	sc.register = append(sc.register, new(strings.Builder))
	return sc.register[len(sc.register)-1]
}

func (sc *ServiceTray) BorrowTiersWriter() io.Writer {
	sc.tiers = append(sc.tiers, new(strings.Builder))
	return sc.tiers[len(sc.tiers)-1]
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
