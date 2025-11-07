package packs

type RepoKeyFuncTmplPackParam struct {
	FieldName string
	ParamName string
}

type RepoKeyFuncTmplPack struct {
	KeyNameCombine        string
	ModelType             string
	Params                []*RepoKeyFuncTmplPackParam
	FetchFuncName         string
	FetchReturnType       string
	ModelPackage          string
	ModelPackageAlias     string
	IsModelAnotherPackage bool

	Tmpls       []string
	IsUniqueKey bool
}

func (pack *RepoKeyFuncTmplPack) Copy() *RepoKeyFuncTmplPack {
	params := make([]*RepoKeyFuncTmplPackParam, len(pack.Params))
	for index := range pack.Params {
		params[index] = &RepoKeyFuncTmplPackParam{
			FieldName: pack.Params[index].FieldName,
			ParamName: pack.Params[index].ParamName,
		}
	}
	return &RepoKeyFuncTmplPack{
		KeyNameCombine:        pack.KeyNameCombine,
		ModelType:             pack.ModelType,
		Params:                params,
		FetchFuncName:         pack.FetchFuncName,
		FetchReturnType:       pack.FetchReturnType,
		ModelPackage:          pack.ModelPackage,
		ModelPackageAlias:     pack.ModelPackageAlias,
		IsModelAnotherPackage: pack.IsModelAnotherPackage,
		Tmpls:                 pack.Tmpls,
		IsUniqueKey:           pack.IsUniqueKey,
	}
}
