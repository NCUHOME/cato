package ware

import (
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
)

type fileGenerator func(ctx *common.GenContext) ([]*models.GenerateFileDesc, error)
