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
type (
	DiscernServer struct {
		Name string
		Port string
		Mode string
		Reg  string
	}

	CommonItem struct {
		Name     string
		Location string
		Key      string
		Value    string
	}

	ConfigNascanInfo struct {
		ScanList      string          `json:"Scan_list"`
		DiscernCms    []CommonItem    `bson:"Discern_cms"`
		DiscernCon    []CommonItem    `bson:"Discern_con"`
		DiscernLang   []CommonItem    `bson:"Discern_lang"`
		DiscernServer []DiscernServer `bson:"Discern_server"`
		PortList      string          `bson:"Port_list"`
		Masscan       string          `bson:"Masscan"`
		Timeout       string          `bson:"Timeout"`
		Cycle         string          `bson:"Cycle"`
		Thread        string          `bson:"Thread"`
		WhiteList     string          `bson:"White_list"`
	}
)
