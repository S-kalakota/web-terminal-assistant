package server

// Config controls the local web terminal server.
type Config struct {
	Addr   string
	WebDir string
}

// Server owns HTTP routes for the local app.
type Server struct {
	cfg Config
}

// New creates a server with the provided configuration.
func New(cfg Config) *Server {
	if cfg.Addr == "" {
		cfg.Addr = "127.0.0.1:8080"
	}
	if cfg.WebDir == "" {
		cfg.WebDir = "web"
	}

	return &Server{cfg: cfg}
}
