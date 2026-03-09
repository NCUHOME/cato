package packs

type RepoTmplPack struct {
	RepoPackageName       string
	IsModelAnotherPackage bool
	IsRdbAnotherPackage   bool
	ModelPackageAlias     string
	RdbPackageAlias       string
	ModelPackage          string
	RepoFuncs             []string
	RdbPackage            string
	ModelType             string
}
