package handler

import (
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/wolfsblu/go-chef/infra/env"
)

func NewUploadHandler() (*tusd.Handler, error) {
	store := filestore.New(env.MustGet("UPLOAD_FOLDER"))
	locker := filelocker.New(env.MustGet("UPLOAD_FOLDER"))

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		return nil, err
	}

	return handler, nil
}
