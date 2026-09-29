package lakepathz

import (
	"github.com/infinity6-ai/gox/commonz/hivepartz"
	"github.com/infinity6-ai/gox/commonz/pathz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketablez"
)

type Path struct {
	Table     laketablez.Table     `json:"table"`
	Partition *hivepartz.HiveParts `json:"partition"`
	File      *pathz.Path          `json:"file"`
}
