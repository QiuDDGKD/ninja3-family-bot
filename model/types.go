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
