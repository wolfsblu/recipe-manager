package handler

import (
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/wolfsblu/go-chef/infra/env"
	"regexp"
)

func NewUploadHandler() (*tusd.Handler, error) {
	store := filestore.New(env.MustGet("UPLOAD_FOLDER"))
	locker := filelocker.New(env.MustGet("UPLOAD_FOLDER"))

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	// corsAllowOrigin := regexp.MustCompile(".*") //env.MustGet("CORS_ORIGIN"))

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/api/uploads/",
		StoreComposer:         composer,
		NotifyCompleteUploads: false,
		Cors: &tusd.CorsConfig{
			AllowOrigin:      regexp.MustCompile(".*"),
			AllowCredentials: true,
			AllowMethods:     tusd.DefaultCorsConfig.AllowMethods,
			AllowHeaders:     tusd.DefaultCorsConfig.AllowHeaders,
			MaxAge:           tusd.DefaultCorsConfig.MaxAge,
			ExposeHeaders:    tusd.DefaultCorsConfig.ExposeHeaders,
		},
	})

	if err != nil {
		return nil, err
	}
	return handler, nil
}
