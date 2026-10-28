package plugson

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
)

//go:embed all:www
var contentEmbed embed.FS

type Config struct {
	ThemeName   string `json:"theme_name"`
	Timeout     int    `json:"timeout"`
	Persistence bool   `json:"persistence"`
}

// NewServer creates a new http.Handler for Plugson.
func NewServer() (http.Handler, error) {
	mux := http.NewServeMux()

	// Serve embedded static files
	subFS, err := fs.Sub(contentEmbed, "www")
	if err != nil {
		return nil, err
	}

	mux.Handle("/", http.FileServer(http.FS(subFS)))

	// API Endpoints
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		cfg := Config{
			ThemeName:   "ventoy-adwaita",
			Timeout:     10,
			Persistence: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg)
	})

	return mux, nil
}
