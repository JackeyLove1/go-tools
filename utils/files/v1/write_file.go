package main

// go run main.go > run.log 1>&2
import (
    "crypto/md5"
    "fmt"
    "io"
    "io/ioutil"
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
    WorkName              = "bigfile"
    Kb                    = 1024
    Mb                    = 1024 * Kb
    Gb                    = 1024 * Mb
    Tb                    = 1024 * Gb
    ChunkSize             = 2 * Mb
    FixedChar             = "a"
    PadChar               = "\x00"
    FileNameSize          = 20
    TruncateMaxSize int64 = 8 * Tb
    Operations            = 100
    randomStrs            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
    NumberOfGoroutines       = 2000
    FileSize                 = 6 * Tb
    MountDir                 = "/mnt"
    DataDir                  = ""
    FileName                 = ""
    FilePath                 = ""
    Content                  = generateContent(ChunkSize)
    Local                    = true
    GlobalSize         int64 = 0 // for single thread truncate and random write
    ValidSize          int64 = 0
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
    log.Printf("Succeed to Create file:%s\n", FilePath)
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
        log.Printf("Succeed to Write file: %s, offset:%d", FilePath, offset)
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
    ValidSize = int64(FileSize)
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

// Expected all content is GlobalSize * FixedChar
func isChunkConsistOf(chunk []byte, char byte) bool {
    for idx, b := range chunk {
        if b != char {
            log.Printf("Failed to compare content char, File: %s,idx:%d, actual:%s, expected all:%s \n", FilePath, idx, string(b), string(char))
            return false
        }
    }
    return true
}

func CompareContent() {
    file, err := os.Open(FilePath)
    failOnErr(err, "Open file")
    defer file.Close()
    fileInfo, err := file.Stat()
    failOnErr(err, "Get file info")
    log.Println("Start to compare content ... ")
    fileSize := fileInfo.Size()
    semaphore := make(chan struct{}, NumberOfGoroutines)
    var wg sync.WaitGroup
    readChunk := func(file *os.File, offset int64, targetChar byte, wg *sync.WaitGroup, semaphore chan struct{}) {
        defer wg.Done()
        defer func() { <-semaphore }()
        chunk, err := ioutil.ReadAll(io.NewSectionReader(file, offset, int64(ChunkSize)))
        failOnErr(err, "Read chunk")
        if !isChunkConsistOf(chunk, targetChar) {
            log.Fatalf("Failed to compare content, File: %s, offset:%d, actual: %s, expected all:%s \n", FilePath, offset, string(chunk), string(targetChar))
        }
    }
    for i := 0; i < int(ValidSize/ChunkSize); i++ {
        wg.Add(1)
        semaphore <- struct{}{}
        go readChunk(file, int64(i*ChunkSize), byte(FixedChar[0]), &wg, semaphore)
    }
    for i := int(ValidSize / ChunkSize); i < int(fileSize/ChunkSize); i++ {
        wg.Add(1)
        semaphore <- struct{}{}
        go readChunk(file, int64(i*ChunkSize), byte(PadChar[0]), &wg, semaphore)
    }
    wg.Wait()
    log.Printf("Succeed to compare content, file size:%d, validDataSize:%d\n", fileSize, ValidSize)
}

func truncateSize(size int64) {
    err := os.Truncate(FilePath, size)
    failOnErr(err, "Truncate file")
}

func AppendWrite(times int) {
    file, err := os.Open(FilePath)
    failOnErr(err, "Open file")
    defer file.Close()
    fileInfo, _ := file.Stat()
    for i := 0; i < times; i++ {
        file.Write([]byte(Content))
        log.Printf("Succeed to Append Write to file: %s, offset:%d, chunkSize:%d", FileName, fileInfo.Size(), ChunkSize)
    }
}

func Max(a, b int64) int64 {
    if a < b {
        return b
    }
    return a
}

func RandomOffsetWrite(times int) {
    file, err := os.OpenFile(FilePath, os.O_WRONLY, 0644)
    failOnErr(err, "Open file")
    defer file.Close()
    offset := rand.Int63() % int64(ValidSize/2)
    for i := 0; i < times; i++ {
        _, err = file.WriteAt([]byte(Content), offset)
        failOnErr(err, "random Write file")
        ValidSize = Max(ValidSize, offset+int64(ChunkSize))
        log.Printf("Succeed to Random Write offset:%d, chunkSize:%d\n", offset, ChunkSize)
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
        ValidSize = targetSize
        log.Printf("Succeed to truncate small file:%s truncateSize:%d\n", FileName, targetSize)
    }
}

func truncateLarge(times int) {
    file, err := os.Open(FilePath)
    failOnErr(err, "Failed to truncateLarge Open file")
    defer file.Close()
    fileInfo, _ := file.Stat()
    for i := 0; i < times; i++ {
        targetSize := fileInfo.Size() + int64(ChunkSize)
        if targetSize > TruncateMaxSize {
            log.Fatalf("Failed to truncateLarge, execeed MaxTruncateSize: %d, targetSize:%d, TruncateMaxSize:%d", TruncateMaxSize,
                targetSize, TruncateMaxSize)
        }
        err = os.Truncate(FilePath, targetSize)
        failOnErr(err, fmt.Sprintf("Failed to Truncate file to size:%d", targetSize))
        log.Printf("Succeed to truncate large file:%s truncateSize:%d\n", FileName, targetSize)
    }
}

func run() {
    WriteBigFile()
    CompareContent()
    for i := 0; i < Operations; i++ {
        // AppendWrite(1)
        CompareContent()
        truncateSmall(1)
        CompareContent()
        truncateLarge(1)
        CompareContent()
        RandomOffsetWrite(1)
        CompareContent()
    }
}

func localTest() {
    NumberOfGoroutines = 10
    FileSize = 100 * Mb
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
    Local = true
    Init()
    if Local {
        localTest()
    } else {
        MountTest()
    }
}
