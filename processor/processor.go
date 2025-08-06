package processor

import (
	"context"
	"log"
	"ninja3-family-bot/model"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
)

type Processor struct {
	Ctx context.Context
	Api openapi.OpenAPI
	DB  *gorm.DB
}

type ProcessorConfig struct {
	QQBotCredentials *token.QQBotCredentials `yaml:"QQBotCredentials"`
	MysqlDSN         string                  `yaml:"MysqlDSN"`
}

func NewProcessor(conf *ProcessorConfig) *Processor {
	ctx := context.Background()
	// 创建 oauth2 标准 token source
	tokenSource := token.NewQQBotTokenSource(conf.QQBotCredentials)
	// 启动自动刷新 access token 协程
	if err := token.StartRefreshAccessToken(ctx, tokenSource); err != nil {
		log.Fatalln(err)
	}
	// 初始化 openapi，正式环境
	api := botgo.NewOpenAPI(conf.QQBotCredentials.AppID, tokenSource).WithTimeout(5 * time.Second).SetDebug(true)

	// 使用 GORM 连接 MySQL
	db, err := gorm.Open(mysql.Open(conf.MysqlDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(&model.User{}, &model.AbyssSignUp{})

	return &Processor{
		Ctx: ctx,
		Api: api,
		DB:  db,
	}
}

func (p *Processor) ProcessGroupMessage(input string, data *dto.WSGroupATMessageData) error {
	splits := strings.Split(input, " ")
	if len(splits) < 1 {
		p.Api.PostGroupMessage(p.Ctx, data.GroupID, &dto.MessageToCreate{
			MsgID:   data.ID,
			Content: "怎么不说话喵~",
		})
		return nil // 没有命令，直接返回
	}

	cmd := splits[0]
	params := splits[1:]
	processor, err := p.GetCMDProcessor(cmd)
	if err != nil {
		p.Api.PostGroupMessage(p.Ctx, data.GroupID, &dto.MessageToCreate{
			MsgID:   data.ID,
			Content: err.Error(),
		})
		return nil // 错误信息已发送，直接返回
	}

	if err := processor(data, params...); err != nil {
		p.Api.PostGroupMessage(p.Ctx, data.GroupID, &dto.MessageToCreate{
			MsgID:   data.ID,
			Content: err.Error(),
		})
		return nil // 错误信息已发送，直接返回
	}

	return nil
}
