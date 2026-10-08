package jsparser

import "github.com/infinity6-ai/gox/schemaz/schemaz"

func DescName(d *schemaz.Desc) string {
	if d == nil {
		return ""
	}
	return d.Name
}

func DescSummary(d *schemaz.Desc) string {
	if d == nil {
		return ""
	}
	return d.Summary
}
