package model

type User struct {
	ID       string `gorm:"primaryKey;column:id"`
	Nickname string `gorm:"column:nickname"`
	ATK      int    `gorm:"column:atk"`
}

type AbyssSignUp struct {
	Date     string `gorm:"primaryKey;column:date"`
	UserID   string `gorm:"primaryKey;column:user_id"`
	Nickname string `gorm:"column:nickname"`
	ATK      int    `gorm:"column:atk"`
}

type AbyssLeave struct {
	Date     string `gorm:"primaryKey;column:date"`
	UserID   string `gorm:"primaryKey;column:user_id"`
	Nickname string `gorm:"column:nickname"`
	Reason   string `gorm:"column:reason"`
}

type BattleSignUp struct {
	Date     string `gorm:"primaryKey;column:date"`
	UserID   string `gorm:"primaryKey;column:user_id"`
	Nickname string `gorm:"column:nickname"`
	ATK      int    `gorm:"column:atk"`
	Tp       string `gorm:"column:tp"`
}

var BattleTypeMap = map[string]struct{}{
	"先锋": {},
	"副将": {},
	"主将": {},
	"王牌": {},
	"头目": {},
}

var BattleTypes = []string{
	"先锋",
	"副将",
	"主将",
	"王牌",
	"头目",
}
