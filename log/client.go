// log/server.go中有服务器端的逻辑，但是客户端的服务想使用这个service还是很麻烦
// 为了让客户端的服务方便的使用log/server.go
package log

import (
	"bytes"
	"fmt"
	"io"
	stlog "log"
	"net/http"
	"sds/registry"
)

// 客户端服务本地写日志的设置
func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)                               // 服务端log设置了时间戳，但客户端不需要
	stlog.SetOutput(&clientLogger{url: serviceURL}) // 客户端记录日志的输出应该指向服务端的logger
}

type clientLogger struct {
	url string
}

var _ io.Writer = (*clientLogger)(nil)

// clientLogger需要实现io.Writer这个接口，所以实现以下Write方法
func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer([]byte(data))
	res, err := http.Post(cl.url+"/log", "text/plain", b) // 写到服务端
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to send log message. Service responed with code %v", res.StatusCode)
	}
	return len(data), nil
}
