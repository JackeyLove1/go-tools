package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "sync"
)

var (
    FileSize = 6 * Tb
)

const (
    FilePath           = "/mnt/bytenas/bigfile/wcv9ETLlaYfFknfpSuFf"
    NumberOfGoroutines = 2000
    FixedChar          = "a"
    PadChar            = "\x00"
    Kb                 = 1024
    Mb                 = 1024 * Kb
    Gb                 = 1024 * Mb
    Tb                 = 1024 * Gb
    ChunkSize          = 2 * Mb
)

func failOnErr(err error, msg string) {
    if err != nil {
        panic(fmt.Sprintf("Msg:%s, Error:%s", msg, err.Error()))
    }
}

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
    file, err := os.OpenFile(FilePath, os.O_RDONLY, 0666)
    failOnErr(err, "Open file")
    defer file.Close()
    fileInfo, err := file.Stat()
    failOnErr(err, "Get file info")
    log.Println("Start to compare content ... ")
    fileSize := fileInfo.Size()
    if fileSize != int64(FileSize) {
        log.Fatalf("FileSize:%d is not equal to FileSize:%d", fileSize, FileSize)
    }
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
    for i := 0; i < int(FileSize/ChunkSize); i++ {
        wg.Add(1)
        semaphore <- struct{}{}
        go readChunk(file, int64(i*ChunkSize), byte(FixedChar[0]), &wg, semaphore)
    }
    /*
       for i := int(ValidSize / ChunkSize); i < int(fileSize/ChunkSize); i++ {
           wg.Add(1)
           semaphore <- struct{}{}
           go readChunk(file, int64(i*ChunkSize), byte(PadChar[0]), &wg, semaphore)
       }
    */
    wg.Wait()
    log.Printf("Succeed to compare content, file size:%d, validDataSize:%d\n", fileSize, FileSize)
}

func main() {
    CompareContent()
}
