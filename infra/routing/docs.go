package routing

import (
	"github.com/wolfsblu/recipe-manager/api"
	"net/http"
)

func apiDocs(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, api.DocsFS, "openapi.yml")
}
