package main

// go run write_file.go > bigfile.log
import (
    "crypto/md5"
    "fmt"
    "io"
    "log"
    "math/rand"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

const (
    WorkName        = "bigfile"
    Kb              = 1024
    Mb              = 1024 * Kb
    Gb              = 1024 * Mb
    Tb              = 1024 * Gb
    ChunkSize       = 2 * Mb
    FixedChar       = "a"
    FileNameSize    = 20
    TruncateMaxSize = 8 * Tb
    Operations      = 100
    randomStrs      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
    NumberOfGoroutines = 1000
    FileSize           = 6 * Tb
    MountDir           = "/mnt/"
    DataDir            = ""
    FileName           = ""
    FilePath           = ""
    Content            = generateContent(ChunkSize)
    Local              = true
)

func ExecCommand(op string, args ...string) error {
    cmd := exec.Command(op, args...)
    output, err := cmd.Output()
    log.Printf("Run command:%s output:%s", op, output)
    return err
}

func failOnErr(err error, msg string) {
    if err != nil {
        panic(fmt.Sprintf("Msg:%s, Error:%s", msg, err.Error()))
    }
}

func IsExist(path string) bool {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
    }
    return true
}

func GetSize(path string) int64 {
    file, err := os.Open(path)
    failOnErr(err, "Open file")
    defer file.Close()
    fileInfo, err := file.Stat()
    failOnErr(err, "Get file info")
    return fileInfo.Size()
}

func GetCurrentDir() string {
    dir, err := os.Getwd()
    failOnErr(err, "Get current dir")
    return dir
}

func generateContent(size int) string {
    var builder strings.Builder
    builder.Grow(size)
    for i := 0; i < size; i++ {
        builder.WriteString(FixedChar)
    }
    return builder.String()
}

func randomContent(size int) string {
    bytes := make([]byte, size)
    for i := 0; i < size; i++ {
        bytes[i] = randomStrs[rand.Intn(len(randomStrs))]
    }
    return string(bytes)
}

func MD5(path string) string {
    file, err := os.Open(path)
    if err != nil {
        log.Fatalf("Open file error:%s", err.Error())
    }
    defer file.Close()
    hasher := md5.New()
    if _, err := io.Copy(hasher, file); err != nil {
        log.Fatalf("Copy file error:%s", err.Error())
    }
    md5Sum := hasher.Sum(nil)
    return fmt.Sprintf("%x", md5Sum)
}

func MD5Stream(fileSize int64) string {
    hasher := md5.New()
    for byteProcessed := int64(0); byteProcessed < fileSize; byteProcessed += ChunkSize {
        if byteProcessed+int64(ChunkSize) > int64(FileSize) {
            hasher.Write([]byte(Content[:(int64(FileSize) - byteProcessed)]))
        } else {
            hasher.Write([]byte(Content))
        }
    }
    md5Sum := hasher.Sum(nil)
    return fmt.Sprintf("%x", md5Sum)
}

func WriteBigFile() {
    file, err := os.OpenFile(FilePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(fmt.Sprintf("Failed to create file:%s, Error:%s\n", FilePath, err.Error()))
    }
    defer file.Close()
    semaphore := make(chan struct{}, NumberOfGoroutines)
    var wg sync.WaitGroup
    writeChunk := func(file *os.File, offset int64, wg *sync.WaitGroup, semaphore chan struct{}) {
        defer wg.Done()
        defer func() { <-semaphore }()
        _, err := file.WriteAt([]byte(Content), offset)
        if err != nil {
            log.Fatal(fmt.Sprintf("Failed to write file:%s, Error:%s", FilePath, err.Error()))
        }
        log.Printf("Succeed toWrite offset:%d", offset)
    }
    if FileSize%ChunkSize != 0 {
        log.Fatalf("FileSize:%d is not a multiple of ChunkSize:%d", FileSize, ChunkSize)
    }
    for i := 0; i < FileSize/ChunkSize; i++ {
        wg.Add(1)
        semaphore <- struct{}{}
        go writeChunk(file, int64(i*ChunkSize), &wg, semaphore)
    }
    wg.Wait()
}

func CompareMd5() {
    actual_md5 := MD5(FilePath)
    fileSize := GetSize(FilePath)
    expected_md5 := MD5Stream(fileSize)
    log.Printf("File: %s, Actual md5:%s, Expected md5:%s", FilePath, actual_md5, expected_md5)
    if actual_md5 != expected_md5 {
        log.Fatalf("Failed to compare md5, File: %s, Actual md5:%s, Expected md5:%s", FilePath, actual_md5, expected_md5)
    }
}

func truncateSize(size int64) {
    err := os.Truncate(FilePath, size)
    failOnErr(err, "Truncate file")
}

func randomOffsetWrite(offset int64, nums int) {
    file, err := os.Open(FilePath)
    failOnErr(err, "Open file")
    defer file.Close()
    for i := 0; i < nums; i++ {
        _, err = file.WriteAt([]byte(Content), offset)
        offset += int64(ChunkSize)
        failOnErr(err, "random Write file")
        log.Printf("Succeed to Write offset:%d, chunkSize:%d", offset, ChunkSize)
    }
}

func AppendWrite(times int) {
    file, err := os.Open(FilePath)
    failOnErr(err, "Open file")
    defer file.Close()
    fileInfo, _ := file.Stat()
    for i := 0; i < times; i++ {
        file.Write([]byte(Content))
        log.Printf("Succeed to Append Write offset:%d, chunkSize:%d", fileInfo.Size(), ChunkSize)
    }
}

func truncateSmall(times int) {
    file, err := os.Open(FilePath)
    failOnErr(err, "Open file")
    defer file.Close()
    fileInfo, _ := file.Stat()
    for i := 0; i < times; i++ {
        targetSize := fileInfo.Size() - int64(ChunkSize)
        err := os.Truncate(FilePath, targetSize)
        failOnErr(err, fmt.Sprintf("Failed to Truncate file to size:%d", targetSize))
        log.Printf("Succeed to truncate offset:%d, chunkSize:%d", fileInfo.Size(), ChunkSize)
    }
}

func run() {
    WriteBigFile()
    CompareMd5()
    for i := 0; i < Operations; i++ {
        AppendWrite(1)
        CompareMd5()
        truncateSmall(1)
        CompareMd5()
    }
}

func localTest() {
    NumberOfGoroutines = 10
    FileSize = 1 * Gb
    run()
}

func MountTest() {
    run()
}

func Init() {
    rand.Seed(time.Now().UnixNano())
    if Local {
        MountDir = filepath.Join(GetCurrentDir(), "mount")
    }
    DataDir = filepath.Join(MountDir, WorkName)
    log.Println("DataDir:", DataDir)
    FileName = randomContent(FileNameSize)
    log.Println("FileName:", FileName)
    FilePath = filepath.Join(DataDir, FileName)
    log.Println("FilePath:", FilePath)
    if err := os.MkdirAll(DataDir, 0755); err != nil {
        log.Fatal(fmt.Sprintf("Failed to create file:%s, Error:%s", FilePath, err.Error()))
    }

}

func main() {
    Local = false
    Init()
    if Local {
        localTest()
    } else {
        MountTest()
    }
}
