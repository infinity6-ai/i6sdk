//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package dszapitablelist

import (
	"context"

	"github.com/infinity6-ai/gox/fsz/fsz"
	"github.com/infinity6-ai/gox/fsz/fsz/fszsignmethod"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbpathz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type File struct {
	// Path of file or table.
	Path laketbpathz.Path `json:"path"`
	// File metadata and stats.
	Stat *fsz.FileStat `json:"stat,omitempty"`
	// Presigned access URL.
	SignedUrl string `json:"signed_url,omitempty"`
}

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqQuery struct {
	// Path to list.
	Path laketbpathz.Path `json:"path,omitempty"`
	// Sign method for presigned URLs.
	SignMethod fszsignmethod.SignMethod `json:"sign_method,omitempty"`
	// Whether to list recursively.
	Recursive bool `json:"recursive,omitempty"`
	// Whether to include file stats.
	Stat bool `json:"stat,omitempty"`
	// Pagination cursor.
	FromCursor string `json:"from_cursor,omitempty"`
	// Limit results count.
	Limit int `json:"limit,omitempty"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
	Query  *ReqQuery  `json:"query,omitempty"`
}

type RespPayload struct {
	// Files and tables found.
	Files []*File `json:"files"`
	// Next page pagination cursor.
	NextCursor string `json:"next_cursor,omitempty"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// List Dataset Tables
// Lists files and tables in the specified dataset.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/table-list",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszTableListId struct{}

func (DszTableListId) DszTableList(ctx context.Context, api *Api) (int, error) {
	panic("DszTableList not implemented: run 'go generate ./...'")
}
