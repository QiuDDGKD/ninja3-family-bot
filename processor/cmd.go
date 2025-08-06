package processor

import (
	"errors"
	"ninja3-family-bot/model"
	"ninja3-family-bot/tools"
	"strconv"

	"github.com/tencent-connect/botgo/dto"
	"gorm.io/gorm"
)

type CMDProcessor func(*dto.WSGroupATMessageData, ...string) error

func (p *Processor) GetCMDProcessor(cmd string) (CMDProcessor, error) {
	switch cmd {
	case "/深渊报名":
		return p.AbyssSignUp, nil
	}

	return nil, errors.New("不知道你要干嘛喵~")
}

func (p *Processor) AbyssSignUp(data *dto.WSGroupATMessageData, params ...string) error {
	var User model.User
	if len(params) < 2 {
		if err := p.DB.Where("id = ?", data.Author.ID).First(&User).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("还没有登记信息，需要传入登记信息喵~")
			}
		}
	} else {
		nickname := params[0]
		atkStr := params[1]
		atk, err := strconv.Atoi(atkStr)
		if err != nil {
			return errors.New("面板必须是数字喵~")
		}

		User = model.User{
			ID:       data.Author.ID,
			Nickname: nickname,
			ATK:      atk,
		}

		if err := p.DB.Save(&User).Error; err != nil {
			return errors.New("保存用户信息失败了喵~")
		}
	}

	abyssSignUp := model.AbyssSignUp{
		Date:     tools.GetNextFriday(),
		UserID:   User.ID,
		Nickname: User.Nickname,
		ATK:      User.ATK,
	}
	if err := p.DB.Save(&abyssSignUp).Error; err != nil {
		return errors.New("报名失败了喵~")
	}

	_, err := p.Api.PostGroupMessage(p.Ctx, data.GroupID, dto.MessageToCreate{
		MsgID:   data.ID,
		Content: "报名成功了喵~",
	})
	if err != nil {
		return errors.New("回复消息失败了喵~")
	}

	return nil
}
