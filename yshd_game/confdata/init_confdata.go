package confdata

//tabtoy --mode=exportorv2 --json_out=config.json Globals.xlsx task.xlsx
//tabtoy --mode=exportorv2 --go_out=task.go --combinename=Config Globals.xlsx task.xlsx
//
var ConfigData *ConfigTable

func InitConfData() {
	ConfigData = NewConfigTable()
	/*
		if err := ConfigData.Load("../config.json"); err != nil {
			panic(err)
		}
	*/
}
