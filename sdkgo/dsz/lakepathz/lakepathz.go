package lakepathz

import (
	"github.com/infinity6-ai/gox/commonz/hivepartz"
	"github.com/infinity6-ai/gox/commonz/pathz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketablez"
)

type Path struct {
	Table      laketablez.Table     `json:"table"`
	Partitions *hivepartz.HiveParts `json:"partitions"`
	File       *pathz.Path          `json:"file"`
}
