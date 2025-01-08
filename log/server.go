package log

import (
	"io"
	stlog "log" // 要自己写一个日志服务，但是要用到标准log库
	"net/http"
	"os"
)

var log *stlog.Logger

type fileLog string // 日志服务是一个web服务，接受post请求，把post请求内容写入日志中

// fileLog是写入文件的路径
func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

func Run(destination string) {
	// 创建一个新的日志记录器，指定输出目标，并加上前缀
	log = stlog.New(fileLog(destination), "go: ", stlog.LstdFlags)
}

// 注册函数，将日志服务注册在localhost:4000/log里面
func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := io.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func write(message string) {
	log.Printf("%v\n", message)
}
