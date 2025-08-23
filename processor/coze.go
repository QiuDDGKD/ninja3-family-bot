package processor

import (
	"context"
	"fmt"

	"github.com/coze-dev/coze-go"
)

type CozeConf struct {
	ApiToken string `yaml:"ApiToken"`
	BotID    string `yaml:"BotID"`
	UserID   string `yaml:"UserID"`
	BaseUrl  string `yaml:"BaseUrl"`
}

type CozeContext struct {
	GroupId string
}

type Coze struct {
	ctx  *CozeContext
	conf *CozeConf
}

func NewCoze(ctx *CozeContext, conf *CozeConf) *Coze {
	return &Coze{
		ctx:  ctx,
		conf: conf,
	}
}

func (c *Coze) GetResponse(input string) (string, error) {
	authCli := coze.NewTokenAuth(c.conf.ApiToken)
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(c.conf.BaseUrl))

	ctx := context.Background()
	req := &coze.CreateChatsReq{
		BotID:  c.conf.BotID,
		UserID: c.ctx.GroupId,
		Messages: []*coze.Message{
			coze.BuildUserQuestionText(input, nil),
		},
	}

	// The sdk provide an automatic polling method.
	timeout := 5
	chat2, err := cozeCli.Chat.CreateAndPoll(ctx, req, &timeout)
	if err != nil {
		fmt.Println("Error in CreateAndPoll:", err)
		return "", err
	}

	respContent := ""
	for _, msg := range chat2.Messages {
		fmt.Println("Message:", msg.Content, "MessageType:", msg.Type)
	}
	for _, msg := range chat2.Messages {
		if msg.Type != coze.MessageTypeAnswer {
			continue
		}

		respContent = msg.Content
		break
	}

	return respContent, nil
}
