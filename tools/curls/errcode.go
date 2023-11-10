package curls

import "errors"

var (
    CharsetDecoderError    = errors.New("failed to decode context")
    ProxyError             = errors.New("failed to set proxy or proxy is unreachable ")
    DownloadFileEmptyError = errors.New("download dir is empty")
)
