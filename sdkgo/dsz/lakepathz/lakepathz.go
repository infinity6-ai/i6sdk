package lakepathz

import (
	"fmt"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/hivepartz"
	"github.com/infinity6-ai/gox/commonz/pathz"
	"github.com/infinity6-ai/gox/commonz/validation"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketablez"
)

type Path struct {
	Table     laketablez.Value
	Partition *hivepartz.HiveParts
	File      *pathz.Path
}

func (p Path) Format() (*pathz.Path, error) {
	ret := pathz.New(0, []string{
		p.Table.Database.Get(),
		"table",
		p.Table.Name.Get(),
	}, false)
	if p.Partition != nil {
		ret = ret.MustJoin(p.Partition.Format())
	}
	if p.File != nil {
		err := p.File.Validate(pathz.ValidateOptions{
			Absolute:   new(false),
			MaxParents: new(0),
			MaxPart:    new(20),
			Part: func(_ int, part string) error {
				return validation.StrContains("=", part, "p.File was not supposed to have '=' since it is not a partition")
			},
		})
		if err != nil {
			return nil, fmt.Errorf("file not supported %s: %w", p.File.String(), err)
		}
		ret = ret.MustJoin(p.File)
	}
	return ret, nil
}

func (p Path) MustFormat() *pathz.Path {
	ret, err := p.Format()
	errorz.Check(err)
	return ret
}

func (p Path) String() string {
	return p.MustFormat().String()
}
