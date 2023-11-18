package init

import (
    "fmt"
    "math/rand"
    "os"
    "strings"
    "time"

    "github.com/rs/zerolog"
    "gopkg.in/ini.v1"
)

const (
    configFilePath = "./configs/config.ini"
)

var stdOutLogger = zerolog.New(os.Stdout)

type KafkaProducerConfig struct {
    Host            string
    Port            string
    RequireACKs     string
    Partitioner     string
    ReturnSuccesses bool
}

type KafkaConsumerConfig struct {
    Host string
    Port string
}

type ossConfig struct {
    Url             string
    Bucket          string
    BucketDirectory string
    AccessKeyID     string
    AccessKeySecret string
}

type videoConfig struct {
    SavePath      string
    AllowedExts   []string
    UploadMaxSize int64
}

type UserConfig struct {
    PasswordEncrypted bool
}

type LogConfig struct {
    LogFileWritten bool
    LogFilePath    string
}

type RpcConfig struct {
    UserServiceHost     string
    UserServicePort     string
    VideoServiceHost    string
    VideoServicePort    string
    FavoriteServiceHost string
    FavoriteServicePort string
    CommentServiceHost  string
    CommentServicePort  string
    FollowServiceHost   string
    FollowServicePort   string
    MessageServiceHost  string
    MessageServicePort  string
}

// parse file config
var (
    Port       string
    dbHost     string
    dbPort     string
    dbUser     string
    dbPassWord string
    dbName     string
    dbLogLevel string

    rdbHost string
    rdbPort string

    FeedListLength  int
    kafkaServerConf KafkaProducerConfig
    kafkaClientConf KafkaConsumerConfig

    OssConf ossConfig

    VideoConf videoConfig

    UserConf UserConfig

    LogConf LogConfig

    RpcCSConf RpcConfig
    RpcSDConf RpcConfig
)

func InitConfig() {
    stdOutLogger.Info().Msg("init config")
    f, err := ini.Load(configFilePath)
    if err != nil {
        panic(fmt.Errorf("failed to init Config, err:%w", err))
    }
    rand.Seed(time.Now().Unix())

    loadServer(f)
    loadDb(f)
    loadRdb(f)
    loadKafkaServer(f)
    loadKafkaClient(f)
    loadFeed(f)
    loadOss(f)
    loadVideo(f)
    loadUser(f)
    loadLog(f)
    loadRpcCSConf(f)
    loadRpcSDConf(f)
}

func loadServer(file *ini.File) {
    s := file.Section("server")
    Port = s.Key("Port").MustString("8888")
}

func loadDb(file *ini.File) {
    s := file.Section("database")
    dbName = s.Key("DbName").MustString("tic-tok")
    dbPort = s.Key("DbPort").MustString("3306")
    dbHost = s.Key("DbHost").MustString("127.0.0.1")
    dbUser = s.Key("DbUser").MustString("")
    dbPassWord = s.Key("DbPassWord").MustString("")
    dbLogLevel = s.Key("LogLevel").MustString("error")
}

func loadRdb(file *ini.File) {
    s := file.Section("redis")
    rdbHost = s.Key("UserServiceHost").MustString("127.0.0.1")
    rdbPort = s.Key("Port").MustString("6379")
}

func loadKafkaServer(file *ini.File) {
    s := file.Section("kafkaProducer")
    kafkaServerConf.Host = s.Key("UserServiceHost").MustString("127.0.0.1")
    kafkaServerConf.Port = s.Key("Port").MustString("9092")
    kafkaServerConf.RequireACKs = s.Key("RequireACKs").MustString("WaitForAll")
    kafkaServerConf.Partitioner = s.Key("ProducerPartitioner").MustString("NewRandomPartitioner")
    kafkaServerConf.ReturnSuccesses = s.Key("ProducerReturnSuccesses").MustBool(true)
}

func loadKafkaClient(file *ini.File) {
    s := file.Section("kafkaConsumer")
    kafkaClientConf.Host = s.Key("UserServiceHost").MustString("127.0.0.1")
    kafkaClientConf.Port = s.Key("Port").MustString("9092")
}

func loadFeed(file *ini.File) {
    s := file.Section("feed")
    FeedListLength = s.Key("ListLength").MustInt(30)
}

func loadOss(file *ini.File) {
    s := file.Section("oss")
    OssConf.Url = s.Key("Url").MustString("")
    OssConf.Bucket = s.Key("Bucket").MustString("")
    OssConf.BucketDirectory = s.Key("BucketDirectory").MustString("")
    OssConf.AccessKeyID = s.Key("AccessKeyID").MustString("")
    OssConf.AccessKeySecret = s.Key("AccessKeySecret").MustString("")
}

func loadVideo(file *ini.File) {
    s := file.Section("video")
    VideoConf.SavePath = s.Key("SavePath").MustString("../userdata/")
    videoExts := s.Key("AllowedExts").MustString("mp4,wmv,avi")
    VideoConf.AllowedExts = strings.Split(videoExts, ",")
    VideoConf.UploadMaxSize = s.Key("UploadMaxSize").MustInt64(1024)
}

func loadUser(file *ini.File) {
    s := file.Section("user")
    UserConf.PasswordEncrypted = s.Key("PasswordEncrypted").MustBool(false)
}

func loadLog(file *ini.File) {
    s := file.Section("log")
    LogConf.LogFileWritten = s.Key("FileLogWritten").MustBool(false)
    LogConf.LogFilePath = s.Key("LogFilePath").MustString("./logdata/logFile.txt")
}

func loadRpcCSConf(file *ini.File) {
    s := file.Section("rpcCS")
    RpcCSConf.UserServiceHost = s.Key("UserServiceHost").MustString("127.0.0.1")
    RpcCSConf.UserServicePort = s.Key("UserServicePort").MustString(":50051")
    RpcCSConf.VideoServiceHost = s.Key("VideoServiceHost").MustString("127.0.0.1")
    RpcCSConf.VideoServicePort = s.Key("VideoServicePort").MustString(":50052")
    RpcCSConf.FavoriteServiceHost = s.Key("FavoriteServiceHost").MustString("127.0.0.1")
    RpcCSConf.FavoriteServicePort = s.Key("FavoriteServicePort").MustString(":50053")
    RpcCSConf.CommentServiceHost = s.Key("CommentServiceHost").MustString("127.0.0.1")
    RpcCSConf.CommentServicePort = s.Key("CommentServicePort").MustString(":50054")
    RpcCSConf.FollowServiceHost = s.Key("FollowServiceHost").MustString("127.0.0.1")
    RpcCSConf.FollowServicePort = s.Key("FollowServicePort").MustString(":50055")
    RpcCSConf.MessageServiceHost = s.Key("MessageServiceHost").MustString("127.0.0.1")
    RpcCSConf.MessageServicePort = s.Key("MessageServicePort").MustString(":50056")

}

func loadRpcSDConf(file *ini.File) {
    s := file.Section("rpcSD")
    RpcSDConf.UserServiceHost = s.Key("UserServiceHost").MustString("127.0.0.1")
    RpcSDConf.UserServicePort = s.Key("UserServicePort").MustString(":50061")
    RpcSDConf.VideoServiceHost = s.Key("VideoServiceHost").MustString("127.0.0.1")
    RpcSDConf.VideoServicePort = s.Key("VideoServicePort").MustString(":50062")
    RpcSDConf.FavoriteServiceHost = s.Key("FavoriteServiceHost").MustString("127.0.0.1")
    RpcSDConf.FavoriteServicePort = s.Key("FavoriteServicePort").MustString(":50063")
    RpcSDConf.CommentServiceHost = s.Key("CommentServiceHost").MustString("127.0.0.1")
    RpcSDConf.CommentServicePort = s.Key("CommentServicePort").MustString(":50064")
    RpcSDConf.FollowServiceHost = s.Key("FollowServiceHost").MustString("127.0.0.1")
    RpcSDConf.FollowServicePort = s.Key("FollowServicePort").MustString(":50065")
    RpcSDConf.MessageServiceHost = s.Key("MessageServiceHost").MustString("127.0.0.1")
    RpcSDConf.MessageServicePort = s.Key("MessageServicePort").MustString(":50066")
}
