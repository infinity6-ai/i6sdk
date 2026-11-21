package stripewebhook

import (
	"context"
	"net/http"

	"github.com/infinity6-ai/gox/commonz/logz"
	"go.code.infinity6.ai/platform/httpz"
	"go.code.infinity6.ai/platform/httpz/httpzserver"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

const Path = "/api/stripe/webhook"

type event struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Livemode bool   `json:"livemode"`
	Created  int64  `json:"created"`
}

func receive(ctx context.Context, req *httpz.MReq) *httpz.MResp {
	e := &event{}
	req.ParseJson(e)
	logger.Info(ctx, "stripe webhook", map[string]any{
		"id":       e.Id,
		"type":     e.Type,
		"livemode": e.Livemode,
		"created":  e.Created,
		"body":     req.ParseAscii(4096),
	})
	return req.Resp(http.StatusNoContent)
}

func Add(s *httpzserver.Server) {
	s.AddRestPrefix("POST", Path, receive)
}
