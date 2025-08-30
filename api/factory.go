package api

func NewAPIServer(h Handler, sec SecurityHandler) (*Server, error) {
	errHandler := WithErrorHandler(CustomErrorHandler())
	middleware := WithMiddleware(Authorize())
	return NewServer(h, sec, errHandler, middleware)
}
