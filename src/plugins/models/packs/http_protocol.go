package packs

type HttpProtocolTmplPack struct {
	ServicePackage   string
	ImportPackages   []string
	ServiceName      string
	ServiceNameInner string
	Methods          []string
	RouterBasePath   string
	HttpMethods      []string
	HttpRouters      []string
}

type HttpRouterTmplPack struct {
	Func   string
	Method string
	Uri    string
}
