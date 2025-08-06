package processor

import (
	"context"
	"log"
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
	QQBotCredentials *token.QQBotCredentials
	MysqlDSN         string
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

	return &Processor{
		Ctx: ctx,
		Api: api,
		DB:  db,
	}
}

func (p *Processor) ProcessGroupMessage(input string, data *dto.WSGroupATMessageData) error {
	// 在这里可以使用 p.DB 进行数据库操作
	p.Api.PostGroupMessage(p.Ctx, data.GroupID, &dto.MessageToCreate{})
	return nil
}
