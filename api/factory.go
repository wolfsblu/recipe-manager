package api

import (
	"github.com/wolfsblu/recipe-manager/api/middleware"
)

func NewAPIServer(h Handler, sec SecurityHandler) (*Server, error) {
	eh := WithErrorHandler(CustomErrorHandler())
	mw := WithMiddleware(middleware.Authorize())
	return NewServer(h, sec, eh, mw)
}
