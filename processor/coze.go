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

	chatResp, err := cozeCli.Chat.Create(ctx, req)
	if err != nil {
		fmt.Println("Error creating chats:", err)
		return "", err
	}
	fmt.Println(chatResp)
	fmt.Println(chatResp.LogID())
	chat := chatResp.Chat
	chatID := chat.ID
	conversationID := chat.ConversationID

	// The sdk provide an automatic polling method.
	chat2, err := cozeCli.Chat.CreateAndPoll(ctx, req, nil)
	if err != nil {
		fmt.Println("Error in CreateAndPoll:", err)
		return "", err
	}

	respContent := ""
	for _, msg := range chat2.Messages {
		respContent += msg.Content
	}

	fmt.Printf("Chat ID: %s, Conversation ID: %s\n", chatID, conversationID)
	return respContent, nil
}
