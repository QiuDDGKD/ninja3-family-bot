package processor

import (
	"errors"
	"fmt"
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
	case "/查询深渊报名":
		return p.QueryAbyssSignUp, nil
	case "/深渊请假":
		return p.AbyssLeave, nil
	case "/查询深渊请假":
		return p.QueryAbyssLeave, nil
	case "/登记":
		return p.Register, nil
	case "/家族战报名":
		return p.BattleSignUp, nil
	case "/查询家族战报名":
		return p.QueryBattleSignUp, nil
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

	// 删除请假记录
	date := tools.GetNextFriday()
	if err := p.DB.Where("date = ? AND user_id = ?", date, User.ID).Delete(&model.AbyssLeave{}).Error; err != nil {
		return errors.New("取消请假失败了喵~")
	}

	abyssSignUp := model.AbyssSignUp{
		Date:     date,
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

// 查询深渊报名
func (p *Processor) QueryAbyssSignUp(data *dto.WSGroupATMessageData, params ...string) error {
	var abyssSignUps []model.AbyssSignUp
	if err := p.DB.Where("date = ?", tools.GetNextFriday()).Order("atk desc").Find(&abyssSignUps).Error; err != nil {
		return errors.New("查询报名信息失败了喵~")
	}

	if len(abyssSignUps) == 0 {
		return errors.New("没有人报名喵~")
	}

	var response string = fmt.Sprintf("\n目前报名总人数： %d 人\n", len(abyssSignUps))
	response += "报名名单：\n"
	for _, signUp := range abyssSignUps {
		response += signUp.Nickname + " - 面板: " + strconv.Itoa(signUp.ATK) + "\n"
	}

	_, err := p.Api.PostGroupMessage(p.Ctx, data.GroupID, dto.MessageToCreate{
		MsgID:   data.ID,
		Content: response,
	})
	if err != nil {
		return errors.New("回复消息失败了喵~")
	}

	return nil
}

// 深渊请假
func (p *Processor) AbyssLeave(data *dto.WSGroupATMessageData, params ...string) error {
	if len(params) < 1 {
		return errors.New("需要传入请假理由喵~")
	}

	reason := params[0]

	var user model.User
	if err := p.DB.Where("id = ?", data.Author.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("还没有登记信息，需要先登记喵~")
		}
		return errors.New("查询用户信息失败了喵~")
	}

	// 删除报名记录
	date := tools.GetNextFriday()
	if err := p.DB.Where("date = ? AND user_id = ?", date, user.ID).Delete(&model.AbyssSignUp{}).Error; err != nil {
		return errors.New("取消报名失败了喵~")
	}

	abyssLeave := model.AbyssLeave{
		Date:     date,
		UserID:   user.ID,
		Nickname: user.Nickname,
		Reason:   reason,
	}

	if err := p.DB.Save(&abyssLeave).Error; err != nil {
		return errors.New("请假失败了喵~")
	}

	_, err := p.Api.PostGroupMessage(p.Ctx, data.GroupID, dto.MessageToCreate{
		MsgID:   data.ID,
		Content: "请假成功了喵~",
	})
	if err != nil {
		return errors.New("回复消息失败了喵~")
	}

	return nil
}

// 查询深渊请假
func (p *Processor) QueryAbyssLeave(data *dto.WSGroupATMessageData, params ...string) error {
	var abyssLeaves []model.AbyssLeave
	if err := p.DB.Where("date = ?", tools.GetNextFriday()).Order("user_id").Find(&abyssLeaves).Error; err != nil {
		return errors.New("查询请假信息失败了喵~")
	}

	if len(abyssLeaves) == 0 {
		return errors.New("没有人请假喵~")
	}

	var response string = fmt.Sprintf("\n目前请假总人数： %d 人\n", len(abyssLeaves))
	response += "请假名单：\n"
	for _, leave := range abyssLeaves {
		response += leave.Nickname + " - 理由: " + leave.Reason + "\n"
	}

	_, err := p.Api.PostGroupMessage(p.Ctx, data.GroupID, dto.MessageToCreate{
		MsgID:   data.ID,
		Content: response,
	})
	if err != nil {
		return errors.New("回复消息失败了喵~")
	}

	return nil
}

func (p *Processor) Register(data *dto.WSGroupATMessageData, params ...string) error {
	if len(params) < 2 {
		return errors.New("需要传入昵称和面板喵~")
	}

	nickname := params[0]
	atkStr := params[1]
	atk, err := strconv.Atoi(atkStr)
	if err != nil {
		return errors.New("面板必须是数字喵~")
	}

	user := model.User{
		ID:       data.Author.ID,
		Nickname: nickname,
		ATK:      atk,
	}

	if err := p.DB.Save(&user).Error; err != nil {
		return errors.New("保存用户信息失败了喵~")
	}

	_, err = p.Api.PostGroupMessage(p.Ctx, data.GroupID, dto.MessageToCreate{
		MsgID:   data.ID,
		Content: "登记成功了喵~",
	})
	if err != nil {
		return errors.New("回复消息失败了喵~")
	}

	return nil
}

func (p *Processor) BattleSignUp(data *dto.WSGroupATMessageData, params ...string) error {
	if len(params) < 1 {
		return errors.New("需要传入报名的类型喵~")
	}

	tp := params[0]
	if _, ok := model.BattleTypeMap[tp]; !ok {
		return errors.New("类型必须是 先锋, 副将, 主将, 王牌, 头目 之一喵~")
	}

	var user model.User
	if err := p.DB.Where("id = ?", data.Author.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("还没有登记信息，需要先登记喵~")
		}
		return errors.New("查询用户信息失败了喵~")
	}

	battleSignUp := model.BattleSignUp{
		Date:     tools.GetNextBattleDate(),
		UserID:   user.ID,
		Nickname: user.Nickname,
		ATK:      user.ATK,
		Tp:       tp,
	}

	if err := p.DB.Save(&battleSignUp).Error; err != nil {
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

// 查询家族战报名
func (p *Processor) QueryBattleSignUp(data *dto.WSGroupATMessageData, params ...string) error {
	var battleSignUps []model.BattleSignUp
	if err := p.DB.Where("date = ?", tools.GetNextBattleDate()).Order("tp").Find(&battleSignUps).Error; err != nil {
		return errors.New("查询报名信息失败了喵~")
	}

	if len(battleSignUps) == 0 {
		return errors.New("没有人报名喵~")
	}

	tpSignUpMap := make(map[string][]model.BattleSignUp)
	for _, signUp := range battleSignUps {
		tpSignUpMap[signUp.Tp] = append(tpSignUpMap[signUp.Tp], signUp)
	}

	var response string = fmt.Sprintf("\n目前报名总人数： %d 人\n", len(battleSignUps))
	response += "报名名单：\n"
	for _, tp := range model.BattleTypes {
		if signUps, ok := tpSignUpMap[tp]; ok {
			response += fmt.Sprintf("\n%s：\n", tp)
			for _, signUp := range signUps {
				response += signUp.Nickname + " - 面板: " + strconv.Itoa(signUp.ATK) + "\n"
			}
		}
	}

	_, err := p.Api.PostGroupMessage(p.Ctx, data.GroupID, dto.MessageToCreate{
		MsgID:   data.ID,
		Content: response,
	})
	if err != nil {
		return errors.New("回复消息失败了喵~")
	}

	return nil
}

// 抽奖
func (p *Processor) Gacha(data *dto.WSGroupATMessageData, params ...string) error {
	tp, numStr := params[0], params[1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return errors.New("抽奖数量必须是数字喵~")
	}
	if num <= 0 {
		return errors.New("抽奖数量必须大于 0 喵~")
	}

	result := make([]string, 0, num)
	switch tp {
	case "队长":
		result, err = p.GachaCaptain(num)
	case "成员":
		result, err = p.GachaMember(num, params[2:])
	default:
		err = errors.New("抽奖类型必须是 队长 或 成员 喵~")
	}
	if err != nil {
		return err
	}

	response := fmt.Sprintf("抽取 %d 个 %s 的结果：\n", num, tp)
	for i, name := range result {
		response += fmt.Sprintf("%d. %s\n", i+1, name)
	}
	if len(result) == 0 {
		response = "没有抽到人喵~"
	}

	if len(response) > 2000 {
		response = "抽奖结果太长了喵~ 请减少抽奖数量或类型。"
		return errors.New(response)
	}

	_, err = p.Api.PostGroupMessage(p.Ctx, data.GroupID, &dto.MessageToCreate{
		MsgID:   data.ID,
		Content: response,
	})
	if err != nil {
		return errors.New("回复消息失败了喵~")
	}
	return nil
}
