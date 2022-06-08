package setting

import (
	"creditCalc/utils"
	"encoding/json"
)

type Setting struct {
	Address   string
	Port      string
	DbAddress string
	DbPort    string
	DbUser    string
	DbPass    string
	DbName    string
}

func LoadSetting(filename string) *Setting {
	var setting Setting
	_bytes, e := utils.LoadFile(filename)
	if e != nil {
		utils.Logger.Println(e)
		return nil
	}

	e = json.Unmarshal(_bytes, &setting)
	if e != nil {
		utils.Logger.Println(e)
		return nil
	}

	return &setting
}
