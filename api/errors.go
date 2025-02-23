package api

import (
	"context"
	"encoding/json"
	"errors"
	ogenhttp "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"net/http"
	"strings"
)

const (
	JsonIndentationChar   = " "
	JsonIndentationPrefix = ""
	JsonIndentationWidth  = 2
)

func CustomErrorHandler(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	var (
		code    int
		message string
		ogenErr ogenerrors.Error
	)

	switch {
	case errors.Is(err, ogenhttp.ErrNotImplemented):
		code = http.StatusNotImplemented
		message = http.StatusText(http.StatusNotImplemented)
	case errors.As(err, &ogenErr):
		code = ogenErr.Code()
		message = ogenErr.Error()
	default:
		code = http.StatusInternalServerError
		message = err.Error()
	}

	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	encoder.SetIndent(JsonIndentationPrefix, strings.Repeat(JsonIndentationChar, JsonIndentationWidth))
	_ = encoder.Encode(Error{
		Message: message,
	})
}
