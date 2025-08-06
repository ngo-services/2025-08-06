package http

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ngo-services/2025-08-06/config"
	"github.com/ngo-services/2025-08-06/internal/archive"
	"github.com/ngo-services/2025-08-06/internal/task"
)

type Handler struct {
	cfg        *config.Config
	manager    *task.Manager
	archiveDir string
}

func NewHandler(cfg *config.Config) *Handler { 
	absArchiveDir, _ := filepath.Abs("./archives")
	os.MkdirAll(absArchiveDir, os.ModePerm)
	fmt.Println("[DEBUG] Archive directory:", absArchiveDir)
	return &Handler{
		cfg:        cfg,
		manager:    task.NewManager(cfg.MaxActiveTasks),
		archiveDir: absArchiveDir,
	}
}

func (h *Handler) CreateTask(c *gin.Context) {
	t, err := h.manager.NewTask()
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("[INFO] /tasks called: created task", t.ID)
	c.JSON(http.StatusOK, gin.H{"task_id": t.ID})
}

type AddFileRequest struct {
	URL string `json:"url" binding:"required"`
}

func (h *Handler) AddFile(c *gin.Context) {
	id := c.Param("id")
	var req AddFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	ext := filepath.Ext(req.URL)
	if _, ok := h.cfg.AllowedTypes[ext]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
		return
	}
	file := task.FileLink{URL: req.URL, Filename: filepath.Base(req.URL), Status: "pending"}
	err := h.manager.AddFileToTask(id, file, h.cfg.MaxFilesPerTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If this is the third file, start the packing process
	t, _ := h.manager.GetTask(id)
	if len(t.Files) == h.cfg.MaxFilesPerTask && t.Status == task.StatusPending {
		go h.packArchive(t)
	}

	c.JSON(http.StatusOK, gin.H{"status": "file added"})
}

func (h *Handler) GetTask(c *gin.Context) {
	id := c.Param("id")
	t, ok := h.manager.GetTask(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	resp := gin.H{
		"status": t.Status,
		"files":  t.Files,
		"errors": t.Errors,
	}
	if t.Status == task.StatusReady {
		resp["archive_url"] = "/archives/" + t.ID + "/archive.zip"
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DownloadArchive(c *gin.Context) {
	taskID := c.Param("task_id")
	filePath := filepath.Join(h.archiveDir, taskID, "archive.zip")
	absFilePath, _ := filepath.Abs(filePath)
	fmt.Println("[DEBUG] Trying to serve zip at:", absFilePath)

	if _, err := os.Stat(filePath); err != nil {
		fmt.Println("[ERROR] Archive not found at", absFilePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "archive not found"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=archive.zip")
	c.File(filePath)
}

//download a single file from a task folder
func (h *Handler) DownloadSingleFile(c *gin.Context) {
	taskID := c.Param("task_id")
	filename := c.Param("filename")
	filePath := filepath.Join(h.archiveDir, taskID, filename)
	absFilePath, _ := filepath.Abs(filePath)
	fmt.Println("[DEBUG] Trying to serve file at:", absFilePath)

	if _, err := os.Stat(filePath); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.File(filePath)
}

// Packs all files in a task folder and creates the archive.zip
func (h *Handler) packArchive(t *task.Task) {
	t.Mu.Lock()
	t.Status = task.StatusPacking
	t.Mu.Unlock()

	urls := []string{}
	for _, f := range t.Files {
		urls = append(urls, f.URL)
	}
	taskDir := filepath.Join(h.archiveDir, t.ID)
	results, err := archive.DownloadAndSaveFiles(urls, h.cfg.AllowedTypes, taskDir)
	if err != nil {
		t.Mu.Lock()
		t.Status = task.StatusFailed
		t.Errors = append(t.Errors, "file download failed: "+err.Error())
		t.Mu.Unlock()
		return
	}

	zipPath := filepath.Join(taskDir, "archive.zip")
	absZipPath, _ := filepath.Abs(zipPath)
	fmt.Println("[DEBUG] Writing zip to:", absZipPath)
	if err := archive.ZipFolder(taskDir, zipPath); err != nil {
		fmt.Println("[ERROR] Failed to create zip:", err)
		t.Mu.Lock()
		t.Status = task.StatusFailed
		t.Errors = append(t.Errors, "archive creation failed: "+err.Error())
		t.Mu.Unlock()
		return
	}

	t.Mu.Lock()
	defer t.Mu.Unlock()
	for i, file := range t.Files {
		if errMsg, ok := results[file.URL]; ok && errMsg != "" {
			t.Files[i].Status = "failed"
			t.Files[i].Error = errMsg
			t.Errors = append(t.Errors, file.URL+": "+errMsg)
		} else if ok {
			t.Files[i].Status = "done"
		}
	}
	t.Status = task.StatusReady
	t.ArchiveURL = "/archives/" + t.ID + "/archive.zip"
}
