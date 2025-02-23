package api

func NewAPIServer(h Handler, sec SecurityHandler) (*Server, error) {
	errHandler := WithErrorHandler(CustomErrorHandler)
	return NewServer(h, sec, errHandler)
}
