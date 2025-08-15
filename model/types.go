package model

import "time"

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

type BattleLeave struct {
	Date     string `gorm:"primaryKey;column:date"`
	UserID   string `gorm:"primaryKey;column:user_id"`
	Nickname string `gorm:"column:nickname"`
	Reason   string `gorm:"column:reason"`
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

type AbyssCaptain struct {
	Nickname string `gorm:"primaryKey;column:nickname"`
	Enabled  bool   `gorm:"column:enabled"`
}

type AbyssRecord struct {
	Id       int       `gorm:"primaryKey;column:id"`
	Uid      string    `gorm:"column:uid;uniqueIndex:uniq_uid_date"`
	Date     time.Time `gorm:"type:DATE;column:date;uniqueIndex:uniq_uid_date"`
	Damage   int       `gorm:"column:damage"`
	Times    int       `gorm:"column:times"`
	Nickname string    `gorm:"column:nickname"`
}
