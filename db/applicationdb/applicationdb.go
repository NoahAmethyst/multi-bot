package applicationdb

import (
	"gorm.io/gorm/clause"
	"multi-bot/entity/entity_pb/app_pb"
	"multi-bot/utils/db_util"
	"multi-bot/utils/log"
)

func GetAllApps() ([]*app_pb.Application, error) {
	var allApps []*app_pb.Application

	mysql := db_util.GetDb()
	err := mysql.Model(&app_pb.Application{}).Scan(&allApps).Error

	if err != nil {
		log.Error().Msgf("get all bots error: %s", err)
	}
	return allApps, err
}

func GetTgApplicationGroups(appId string, botType int32) ([]*TgApplicationGroup, error) {
	cd := map[string]interface{}{
		"app_id":    appId,
		"is_delete": 0,
	}
	if botType != 0 {
		cd["bot_type"] = botType
	}
	var appGroups []*TgApplicationGroup
	mysql := db_util.GetDb()
	err := mysql.Model(&appGroups).Where(cd).Scan(&appGroups).Error
	if err != nil {
		log.Error().Msgf("get all tg groups error: %s", err)
	}
	return appGroups, err
}

func InsertTgApplicationGroup(group *TgApplicationGroup) error {

	mysql := db_util.GetDb()
	set := clause.Assignments(map[string]interface{}{
		"is_delete":  false,
		"group_name": group.GroupName,
	})

	err := mysql.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "app_id"}, {Name: "group_id"}},
		DoUpdates: set,
	}).Model(group).Create(group).Error
	if err != nil {
		log.Error().Msgf("create new tg groups error: %s", err)
	}
	return err
}

func DelTgApplicationGroup(group *TgApplicationGroup) error {
	mysql := db_util.GetDb()
	group.IsDelete = true
	err := mysql.Model(group).Updates(group).Error
	if err != nil {
		log.Error().Msgf("del tg groups error: %s", err)
	}
	return err
}

func ApplicationGroupFirstOrInit(group *TgApplicationGroup) error {
	return db_util.GetDb().FirstOrInit(group).Error
}
