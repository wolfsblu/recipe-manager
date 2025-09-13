package routing

import (
	kit "github.com/wolfsblu/recipe-manager/webapp"
	"io/fs"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	sub, _ := fs.Sub(kit.DistFS, "dist")
	http.ServeFileFS(w, r, sub, "index.html")
}

func assets(w http.ResponseWriter, r *http.Request) {
	sub, _ := fs.Sub(kit.DistFS, "dist")
	h := http.FileServerFS(sub)
	h.ServeHTTP(w, r)
}
