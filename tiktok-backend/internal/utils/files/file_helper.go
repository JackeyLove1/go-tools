package files

import (
    "bytes"
    "errors"
    "io"
    "mime/multipart"
    "os"
    "path"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    "tiktok-backend/init"
    "tiktok-backend/internal/utils/constants"
)

// PathExists Judge a file is existed or not
func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if errors.Is(os.ErrNotExist, err) {
        return false, nil
    }
    return false, err
}

// CheckFileExt check the extension of a file name
func CheckFileExt(fileName string) bool {
    ext := path.Ext(fileName)
    ext = string(bytes.ToLower([]byte(ext)))
    for _, legalExt := range init.VideoConf.AllowedExts {
        if legalExt == ext {
            return true
        }
    }
    return false
}

func CheckFileSize(fileSize int64) bool {
    return fileSize <= init.VideoConf.UploadMaxSize*constants.MB
}

func GetFileNameWithoutExt(fileName string) string {
    ext := path.Ext(fileName)
    return strings.TrimSuffix(fileName, ext)
}

func SaveFileToLocal(savePath string, data *multipart.FileHeader) (string, error) {
    if exists, _ := PathExists(savePath); !exists {
        err := os.Mkdir(savePath, os.ModePerm)
        if err != nil {
            return "", err
        }
    }

    src, err := data.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()
    timeLog := time.Now().Unix()
    fileName := GetFileNameWithoutExt(data.Filename)
    fileName += strconv.FormatInt(timeLog, 10) + path.Ext(data.Filename)
    // Caution: this operation will inject "../" or "/" leak
    out, err := os.Create(filepath.Join(savePath, fileName))
    if err != nil {
        return "", err
    }
    defer out.Close()

    _, err = io.Copy(out, src)
    return fileName, err
}

func SaveDataToLocal(savePath string, data *[]byte, filename string) (string, error) {
    if exists, _ := PathExists(savePath); !exists {
        err := os.Mkdir(savePath, os.ModePerm)
        if err != nil {
            return "", err
        }
    }
    src := bytes.NewReader(*data)
    timeLog := time.Now().Unix()
    fileName := GetFileNameWithoutExt(filename)
    filename += strconv.FormatInt(timeLog, 10) + path.Ext(filename)
    out, err := os.Create(filepath.Join(savePath, filename))
    if err != nil {
        return "", err
    }
    defer out.Close()
    _, err = io.Copy(out, src)
    return fileName, err
}
