package archive

import (
    "archive/zip"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

func DownloadAndSaveFiles(urls []string, allowedTypes map[string]struct{}, dir string) (map[string]string, error) {
    result := make(map[string]string)
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        return result, err
    }

    for _, link := range urls {
        ext := strings.ToLower(filepath.Ext(link))
        if _, ok := allowedTypes[ext]; !ok {
            result[link] = "file type not allowed"
            continue
        }
        resp, err := http.Get(link)
        if err != nil || resp.StatusCode != http.StatusOK {
            result[link] = "unable to fetch"
            if resp != nil {
                resp.Body.Close()
            }
            continue
        }
        _, fname := filepath.Split(link)
        outPath := filepath.Join(dir, fname)
        out, err := os.Create(outPath)
        if err != nil {
            result[link] = "failed to save file"
            resp.Body.Close()
            continue
        }
        if _, err := io.Copy(out, resp.Body); err != nil {
            result[link] = "failed to write file"
            out.Close()
            resp.Body.Close()
            continue
        }
        out.Close()
        resp.Body.Close()
        result[link] = ""
    }
    return result, nil
}

func ZipFolder(srcFolder, zipPath string) error {
    zipfile, err := os.Create(zipPath)
    if err != nil {
        return err
    }
    defer zipfile.Close()
    zw := zip.NewWriter(zipfile)
    defer zw.Close()

    entries, err := os.ReadDir(srcFolder)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        if entry.IsDir() || entry.Name() == "archive.zip" {
            continue
        }
        fpath := filepath.Join(srcFolder, entry.Name())
        f, err := os.Open(fpath)
        if err != nil {
            continue
        }
        w, err := zw.Create(entry.Name())
        if err != nil {
            f.Close()
            continue
        }
        io.Copy(w, f)
        f.Close()
    }
    return nil
}
