package processor

import (
	"context"
	"errors"
	"log"
	"ninja3-family-bot/model"
	"ninja3-family-bot/tools"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
)

type Processor struct {
	Ctx    context.Context
	Api    openapi.OpenAPI
	DB     *gorm.DB
	Family *model.Family
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

	db.AutoMigrate(
		&model.Family{},
		&model.GroupFamilyRelation{},
		&model.User{},
		&model.AbyssSignUp{},
		&model.AbyssLeave{},
		&model.BattleSignUp{},
		&model.BattleLeave{},
		&model.AbyssCaptain{},
	)

	return &Processor{
		Ctx: ctx,
		Api: api,
		DB:  db,
	}
}

func (p *Processor) ProcessGroupMessage(input string, data *dto.WSGroupATMessageData) error {
	relation := &model.GroupFamilyRelation{}
	if err := p.DB.Where("group_id = ?", data.GroupID).First(relation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.DB.Create(&model.GroupFamilyRelation{
				GroupID:  data.GroupID,
				FamilyID: "test",
			})
			relation = &model.GroupFamilyRelation{
				GroupID:  data.GroupID,
				FamilyID: "test",
			}
		} else {
			p.Api.PostGroupMessage(p.Ctx, data.GroupID, &dto.MessageToCreate{
				MsgID:   data.ID,
				Content: "家族绑定信息查询失败了喵~",
			})
			return nil // 家族绑定信息查询失败，直接返回
		}
	}

	family := &model.Family{}
	if err := p.DB.Where("id = ?", relation.FamilyID).First(family).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Api.PostGroupMessage(p.Ctx, data.GroupID, &dto.MessageToCreate{
				MsgID:   data.ID,
				Content: "没有找到家族信息喵~",
			})
			return nil
		}
		p.Api.PostGroupMessage(p.Ctx, data.GroupID, &dto.MessageToCreate{
			MsgID:   data.ID,
			Content: "家族信息查询失败了喵~",
		})
		return nil // 家族信息查询失败，直接返回
	}
	p.Family = family

	splits := tools.GetSplits(input)
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
