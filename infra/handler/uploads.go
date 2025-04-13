package handler

import (
	"errors"
	"fmt"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/env"
	"golang.org/x/exp/slog"
	"log"
	"os"
	"regexp"
	"strings"
)

func NewUploadHandler(recipes *domain.RecipeService) (*tusd.Handler, error) {
	store := filestore.New(env.MustGet("UPLOAD_FOLDER"))
	locker := filelocker.New(env.MustGet("UPLOAD_FOLDER"))

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})
	logger := slog.New(logHandler)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/api/uploads/",
		StoreComposer:         composer,
		NotifyCreatedUploads:  true,
		NotifyCompleteUploads: true,
		Logger:                logger,
		Cors: &tusd.CorsConfig{
			AllowOrigin:      regexp.MustCompile(env.MustGet("CORS_ORIGIN")),
			AllowCredentials: true,
			AllowMethods:     tusd.DefaultCorsConfig.AllowMethods,
			AllowHeaders:     fmt.Sprintf("%s, Cookie", tusd.DefaultCorsConfig.AllowHeaders),
			MaxAge:           tusd.DefaultCorsConfig.MaxAge,
			ExposeHeaders:    tusd.DefaultCorsConfig.ExposeHeaders,
		},
	})

	go registerEventListeners(handler, recipes)

	if err != nil {
		return nil, err
	}
	return handler, nil
}

func registerEventListeners(handler *tusd.Handler, recipes *domain.RecipeService) {
	for {
		select {
		case event := <-handler.CompleteUploads:
			onUploadCompleted(recipes, event)
		case event := <-handler.CreatedUploads:
			onUploadCreated(recipes, event)
		}
	}
}

func onUploadCreated(recipes *domain.RecipeService, event tusd.HookEvent) {
	userId, err := getUserIdFromUpload(event.HTTPRequest)
	if err != nil {
		// Abort upload
	}
	log.Printf("userId: %d", userId)
}

func onUploadCompleted(recipes *domain.RecipeService, event tusd.HookEvent) {
	// Move file to correct location
}

func getUserIdFromUpload(req tusd.HTTPRequest) (int64, error) {
	for k, v := range req.Header {
		if k == "Cookie" {
			allCookies := strings.Join(v, ";")
			cookies := strings.Split(allCookies, ";")
			for _, cookie := range cookies {
				if strings.HasPrefix(cookie, CookieName) {
					sessionCookie := strings.Split(cookie, "=")
					return getUserFromSessionCookie(sessionCookie[1])
				}
			}
		}
	}
	return 0, errors.New("upload requires authentication")
}
