package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"ninja3-family-bot/processor"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/interaction/webhook"
	"gopkg.in/yaml.v3"
)

const (
	host_ = "0.0.0.0"
	port_ = 9000
	path_ = "/qqbot"
)

var p *processor.Processor

func main() {
	conf := loadConfig()
	p = processor.NewProcessor(conf)

	// 注册事件处理函数
	_ = event.RegisterHandlers(
		// 注册c2c消息处理函数
		GroupATMessageEventHandler(),
	)
	//注册回调处理函数
	http.HandleFunc(path_, func(writer http.ResponseWriter, request *http.Request) {
		webhook.HTTPHandler(writer, request, conf.QQBotCredentials)
	})
	// 启动http服务监听端口
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host_, port_), nil); err != nil {
		log.Fatal("setup server fatal:", err)
	}
}

// GroupATMessageEventHandler 实现处理 at 消息的回调
func GroupATMessageEventHandler() event.GroupATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content))
		return p.ProcessGroupMessage(input, data)
	}
}

func loadConfig() *processor.ProcessorConfig {
	content, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	config := &processor.ProcessorConfig{}
	if err := yaml.Unmarshal(content, config); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return config
}
