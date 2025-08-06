# 2025-08-06

A RESTful Go API to bundle up to 3 public PDF/JPEG files per task into a zip archive.\
Each task’s files and zip are organized into their own folder.\
Supports downloading the entire zip or any individual file.

---

## 🚀 Quickstart

```bash
git clone https://github.com/ngo-services/2025-08-06.git
cd 2025-08-06

go get github.com/gin-gonic/gin
go get github.com/google/uuid

# Create the folder structure if not present
mkdir -p cmd/server internal/archive internal/task internal/http config

# Create the needed files if you’re starting from scratch
touch cmd/server/main.go config/config.go internal/archive/archiver.go internal/task/manager.go internal/task/types.go internal/http/handler.go internal/http/router.go

# Run the server
go mod tidy
go run ./cmd/server
```

---

## 📂 Folder Structure

```
2025-08-06/
│
├── cmd/server/main.go        # Entry point
├── config/config.go          # Configuration
├── internal/
│   ├── archive/              # File downloading/zipping
│   ├── task/                 # Task manager/types
│   └── http/                 # HTTP handler and router
└── archives/                 # All task folders/files (created automatically)
```

---

## 🛣️ API Endpoints

### 1. Create a Task

\*\*POST \*\*\`\`

- Creates a new task.
- Response:
  ```json
  { "task_id": "YOUR_TASK_ID" }
  ```

---

### 2. Add a File to a Task

\*\*POST \*\*\`\`

- Add a public PDF or JPEG URL.

- JSON body:

  ```json
  { "url": "https://example.com/file.pdf" }
  ```

- Only `.pdf` and `.jpeg` are allowed.

- Max 3 files per task.

- Starts archive creation after 3rd file.

- Response:

  ```json
  { "status": "file added" }
  ```

---

### 3. Get Task Status

\*\*GET \*\*\`\`

- Returns status, file info, and archive URL if ready.
- Example response:
  ```json
  {
    "status": "ready",
    "files": [
      { "URL": "...", "Filename": "...", "Status": "done", "Error": "" }
    ],
    "errors": [],
    "archive_url": "/archives/{task_id}/archive.zip"
  }
  ```

---

### 4. Download Archive (ZIP)

\*\*GET \*\*\`\`

- Downloads the zip archive for this task.
- Only works when status is `"ready"`.

---

### 5. Download Single File

\*\*GET \*\*\`\`

- Downloads a specific file added to the task.

---

### Error Responses

All errors are JSON, for example:

```json
{ "error": "archive not found" }
```

---

## 📝 Example Workflow

1. **Create task:**\
   `POST /tasks` → get `task_id`

2. **Add files:**\
   `POST /tasks/{task_id}/files` (up to 3 times, PDF/JPEG URLs)

3. **Check status:**\
   `GET /tasks/{task_id}` until status is `"ready"`

4. **Download zip:**\
   `GET /archives/{task_id}/archive.zip`

5. **Download a single file:**\
   `GET /archives/{task_id}/files/{filename}`

---

## ⚙️ Configuration

- Allowed file types: `.pdf`, `.jpeg` (see `config/config.go`)
- Max files per task: `3`
- Max concurrent archive tasks: `3`
- Server port: `8080` (edit in `config/config.go` if needed)

---

## 👷 Go Best Practices

- Layered folder structure
- Mutexes for safe task state
- Gin for HTTP routing
- No external DB/Docker, file-system only

---

## 🛑 Troubleshooting

- If archive is not found, check your server console for `[ERROR]` logs and verify the file path.
- All files are under `archives/{task_id}/`.

---

## License

MIT 

