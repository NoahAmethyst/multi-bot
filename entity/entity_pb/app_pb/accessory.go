package app_pb

import (
	"database/sql/driver"
	"encoding/json"
	"multi-bot/utils/log"
)

func (c *CmdConfigs) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

func (c *CmdConfigs) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	if err != nil {
		log.Error().Msgf("json marshal error:%s", err.Error())
	}
	return string(b), err
}
