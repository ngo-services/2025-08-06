package main

import (
	"github.com/ngo-services/2025-08-06/config"
	httpapi "github.com/ngo-services/2025-08-06/internal/http"
)

func main() {
	cfg := config.New()
	r := httpapi.NewRouter(cfg)
	r.Run(cfg.Port)
}
