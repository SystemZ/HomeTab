package web

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	display.HTML(w, http.StatusOK, "index", nil)
}
