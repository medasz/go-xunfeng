package models

type ConfigBase struct {
	Value string `bson:"value"`
	Info  string `bson:"info"`
	Help  string `bson:"help"`
}

type Nascan struct {
	ConfigType   string       `bson:"type"`
	ConfigNascan ConfigNascan `bson:"config"`
}

type ConfigNascan struct {
	ScanList      ConfigBase `bson:"Scan_list"`
	DiscernCms    ConfigBase `bson:"Discern_cms"`
	DiscernCon    ConfigBase `bson:"Discern_con"`
	DiscernLang   ConfigBase `bson:"Discern_lang"`
	DiscernServer ConfigBase `bson:"Discern_server"`
	PortList      ConfigBase `bson:"Port_list"`
	Masscan       ConfigBase `bson:"Masscan"`
	Timeout       ConfigBase `bson:"Timeout"`
	Cycle         ConfigBase `bson:"Cycle"`
	Thread        ConfigBase `bson:"Thread"`
	WhiteList     ConfigBase `bson:"White_list"`
}

type Vulscan struct {
	ConfigType    string        `bson:"type"`
	ConfigVulscan ConfigVulscan `bson:"config"`
}

type ConfigVulscan struct {
	PasswordDic ConfigBase `bson:"Password_dic"`
	Timeout     ConfigBase `bson:"Timeout"`
	Thread      ConfigBase `bson:"Thread"`
	WhiteList   ConfigBase `bson:"White_list"`
}

type ConfigExternal struct {
	Value      string
	Info       string
	Help       string
	IsList     bool
	ConfigType string
	IsSwitch   bool
	Speed      string
}
