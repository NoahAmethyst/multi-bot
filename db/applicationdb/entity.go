package applicationdb

type TgApplicationGroup struct {
	Id        uint   `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	AppId     string `gorm:"column:app_id"`
	GroupId   string `gorm:"column:group_id"`
	GroupName string `gorm:"column:group_name"`
	IsDelete  bool   `gorm:"column:is_delete"`
	BotType   int32  `gorm:"column:bot_type" json:"bot_type"`
}

func (TgApplicationGroup) TableName() string {
	return "tg_application_group"
}
