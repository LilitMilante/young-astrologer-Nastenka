package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	*http.Server
}

func NewServer(port string, h *Handler) *Server {
	r := http.NewServeMux()

	r.HandleFunc("/health", health)
	r.HandleFunc("/image", h.Image) // TODO: rewrite to normal router
	r.HandleFunc("/images", h.Images)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	return &Server{srv}
}

func health(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "OK!")
}
