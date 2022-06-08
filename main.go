package main

import (
	"creditCalc/database"
	"creditCalc/router"
	"creditCalc/setting"
)

func main() {
	s := setting.LoadSetting("setting.json")
	database.Connection(s)
	_ = router.Initialized().Run(s.Address + ":" + s.Port)
}
