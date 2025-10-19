package packs

type HttpHandlerTmplPack struct {
	HttpHandlerPackage     string
	ServiceName            string
	RegisterServiceMethods []string
	RouterBasePath         string
}
