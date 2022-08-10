package obj_util

import (
	"encoding/json"
	"multi-bot/utils/log"
)

func MapToStruct(m interface{}, out interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		log.Error().Msgf("MapToStruct Error %s", err)
		return err
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		log.Error().Msgf("MapToStruct Error %s", err)
		return err
	}

	return nil
}
