package packs

type HttpProtocolTmplPack struct {
	HttpHandlerPackage    string
	ProtocolParamPackages []string
	ServiceName           string
	Methods               []string
	TierMethods           []string
}
