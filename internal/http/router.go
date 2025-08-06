package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ngo-services/2025-08-06/config"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	h := NewHandler(cfg)
	r := gin.Default()

	r.POST("/tasks", h.CreateTask)
	r.POST("/tasks/:id/files", h.AddFile)
	r.GET("/tasks/:id", h.GetTask)
	// r.GET("/archives/:id.zip", h.DownloadArchive)
	// r.GET("/archives/:id", h.DownloadArchive)
	// r.GET("/archives/:archive", h.DownloadArchive)
	r.GET("/archives/:task_id/files/:filename", h.DownloadSingleFile)
	r.GET("/archives/:task_id/archive.zip", h.DownloadArchive)

	return r
}
