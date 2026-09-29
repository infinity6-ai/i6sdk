package lakepathz

import (
	"github.com/infinity6-ai/gox/commonz/hivepartz"
	"github.com/infinity6-ai/gox/commonz/pathz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketablez"
)

type Path struct {
	Table     laketablez.Value
	Partition *hivepartz.HiveParts
	File      *pathz.Path
}

func (p Path) Format() *pathz.Path {
	ret := pathz.New(0, []string{
		p.Table.Database.Get(),
		"table",
		p.Table.Name.Get(),
	}, false)
	if p.Partition != nil {
		ret = ret.MustJoin(p.Partition.Format())
	}
	if p.File != nil {
		p.File.Check(pathz.ValidateOptions{
			Absolute:   new(false),
			MaxParents: new(0),
		})
		ret = ret.MustJoin(p.File)
	}
	return ret
}

func (p Path) String() string {
	return p.Format().String()
}
