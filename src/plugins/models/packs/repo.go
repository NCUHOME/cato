package packs

type RepoTmplPack struct {
	RepoPackageName       string
	IsModelAnotherPackage bool
	ModelPackageAlias     string
	ModelPackage          string
	RepoFuncs             []string
	RdbPackage            string
	ModelType             string
}
