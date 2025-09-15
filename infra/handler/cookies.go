package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/wolfsblu/recipe-manager/infra/env"
)

const AuthCookieName = "SESSID"

func createSessionCookie(userId int64) (string, error) {
	payload, err := encryptUserId(userId)
	if err != nil {
		return "", err
	}
	expiry := 7 * 24 * time.Hour // One week
	cookie := http.Cookie{
		HttpOnly: true,
		MaxAge:   int(expiry / time.Second),
		Name:     AuthCookieName,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Value:    payload,
	}
	return cookie.String(), nil
}

func expireSessionCookie() string {
	cookie := http.Cookie{
		HttpOnly: true,
		MaxAge:   -1,
		Name:     AuthCookieName,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}
	return cookie.String()
}

func getUserFromSessionCookie(cookieValue string) (int64, error) {
	return decryptUserId(cookieValue)
}

func encryptUserId(userId int64) (string, error) {
	var s = securecookie.New(
		[]byte(env.MustGet("COOKIE_HASH_KEY")),
		[]byte(env.MustGet("COOKIE_BLOCK_KEY")),
	)
	encoded, err := s.Encode(AuthCookieName, userId)
	if err != nil {
		return "", err
	}
	return encoded, nil
}

func decryptUserId(cookieValue string) (int64, error) {
	var userId int64
	var s = securecookie.New(
		[]byte(env.MustGet("COOKIE_HASH_KEY")),
		[]byte(env.MustGet("COOKIE_BLOCK_KEY")),
	)
	err := s.Decode(AuthCookieName, cookieValue, &userId)
	if err != nil {
		return -1, err
	}
	return userId, nil
}
