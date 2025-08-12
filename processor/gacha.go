package processor

import (
	"errors"
	"ninja3-family-bot/model"
	"ninja3-family-bot/tools"
	"strconv"
)

type AbyssRecordFilter struct {
	DamageRankMin *int
	DamageMin     *int
	TimesMin      *int
}

func ParseFilter(params []string) AbyssRecordFilter {
	filter := AbyssRecordFilter{}
	for i := 0; i+1 < len(params); i += 2 {
		tp, value := params[i], params[i+1]
		v, err := strconv.Atoi(value)
		if err != nil {
			continue
		}

		switch tp {
		case "排名":
			filter.DamageRankMin = &v
		case "伤害":
			filter.DamageMin = &v
		case "次数":
			filter.TimesMin = &v
		}
	}

	return filter
}

// 抽取队长
func (p *Processor) GachaCaptain(num int) ([]string, error) {
	var captains []model.AbyssCaptain
	if err := p.DB.Where(&model.AbyssCaptain{Enabled: true}).Find(&captains).Error; err != nil {
		return nil, errors.New("查询队长失败了喵~")
	}
	if len(captains) == 0 {
		return nil, errors.New("没有队长信息喵~")
	}

	if num > len(captains) {
		num = len(captains)
	}

	// 随机 num 个小于 len(captains) 且不重复的整数
	indices := tools.RandInts(num, len(captains))
	result := make([]string, 0, num)
	for _, index := range indices {
		captain := captains[index]
		result = append(result, captain.Nickname)
	}
	return result, nil
}

// 抽取成员
func (p *Processor) GachaMember(num int, filterParams []string) ([]string, error) {
	filter := ParseFilter(filterParams)
	stat := p.DB.Order("damage DESC")
	if filter.DamageRankMin != nil {
		stat = stat.Limit(*filter.DamageRankMin)
	}
	if filter.DamageMin != nil {
		stat = stat.Where("damage >= ?", *filter.DamageMin)
	}
	if filter.TimesMin != nil {
		stat = stat.Where("times >= ?", *filter.TimesMin)
	}

	stat = stat.Where("date = ?", tools.GetLastFriday())
	records := make([]model.AbyssRecord, 0)
	if err := stat.Find(&records).Error; err != nil {
		return nil, errors.New("查询记录失败了喵~")
	}

	if len(records) == 0 {
		return nil, errors.New("没有符合条件的记录喵~")
	}

	if num > len(records) {
		num = len(records)
	}
	// 随机 num 个小于 len(records) 且不重复的整数
	indices := tools.RandInts(num, len(records))
	result := make([]string, 0, num)
	for _, index := range indices {
		record := records[index]
		result = append(result, record.Nickname)
	}
	return result, nil
}
