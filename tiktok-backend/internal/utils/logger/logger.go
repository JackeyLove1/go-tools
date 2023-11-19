package logger

import (
    "os"

    "github.com/rs/zerolog"
    "ticktok/init"
    "ticktok/internal/utils/files"
)

var GlobalLogger zerolog.Logger

func InitLogger(config init.LogConfig) {
    var err error
    var file *os.File
    if config.LogFileWritten {
        if exists, _ := files.PathExists(config.LogFilePath); exists {
            file, err = os.OpenFile(config.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        } else {
            file, err = os.Create(config.LogFilePath)
        }
        if err != nil {
            GlobalLogger = init.GetStdOutLogger()
            GlobalLogger.Error().Msg("Get Logger failed")
        }
        GlobalLogger = zerolog.New(file)
    } else {
        GlobalLogger = init.GetStdOutLogger()
    }
    GlobalLogger.Level(zerolog.InfoLevel)
}
