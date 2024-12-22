package smtp

import (
	"bytes"
	"fmt"
	"github.com/wolfsblu/go-chef/infra/env"
	"html/template"
	"net/url"
	"strings"
)

type PasswordResetTemplate struct {
	ResetLink string
}

type UserRegistrationTemplate struct {
	ConfirmLink string
}

func buildTemplate(path string, data any) (string, error) {
	t := template.New(path)
	t, err := t.ParseFS(templateFS, fmt.Sprintf("templates/%s", path))
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func buildUrl(path string) string {
	path = strings.TrimPrefix(path, "/")
	return fmt.Sprintf("%s/%s", getBaseUrl(), path)
}

func buildUrlWithQuery(path string, query map[string]string) string {
	result, _ := url.Parse(buildUrl(path))
	q := result.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	result.RawQuery = q.Encode()
	return result.String()
}

func getBaseUrl() string {
	result := strings.TrimSuffix(env.MustGet("BASE_URL"), "/")
	return fmt.Sprintf("%s", result)
}
