package handler

import (
	tusd "github.com/tus/tusd/v2/pkg/handler"
)

func NewUploadHandler() (*tusd.Handler, error) {
	handler, err := tusd.NewHandler(tusd.Config{
		BasePath: "/files/",
	})
	if err != nil {
		return nil, err
	}
	return handler, nil
}
