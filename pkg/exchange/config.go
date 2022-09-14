package exchange

import (
	"go-xunfeng/models"
	"strings"
)

func NascanConfigMongoToTemplate(nasConfig *models.Nascan) []models.ConfigExternal {
	data := make([]models.ConfigExternal, 0)
	//mascan
	tmp := strings.Split(nasConfig.ConfigNascan.Masscan.Value, "|")
	data = append(data, models.ConfigExternal{
		Value:      tmp[2],
		Info:       nasConfig.ConfigNascan.Masscan.Info,
		Help:       nasConfig.ConfigNascan.Masscan.Help,
		IsList:     false,
		ConfigType: "Masscan",
		IsSwitch:   tmp[0] == "0",
		Speed:      tmp[1],
	})
	//thread
	data = append(data, models.ConfigExternal{
		Value:      nasConfig.ConfigNascan.Thread.Value,
		Info:       nasConfig.ConfigNascan.Thread.Info,
		Help:       nasConfig.ConfigNascan.Thread.Help,
		IsList:     false,
		ConfigType: "Thread",
		IsSwitch:   false,
		Speed:      "",
	})
	//timeout
	data = append(data, models.ConfigExternal{
		Value:      nasConfig.ConfigNascan.Timeout.Value,
		Info:       nasConfig.ConfigNascan.Timeout.Info,
		Help:       nasConfig.ConfigNascan.Timeout.Help,
		IsList:     false,
		ConfigType: "Timeout",
		IsSwitch:   false,
		Speed:      "",
	})
	//cycle
	data = append(data, models.ConfigExternal{
		Value:      nasConfig.ConfigNascan.Cycle.Value,
		Info:       nasConfig.ConfigNascan.Cycle.Info,
		Help:       nasConfig.ConfigNascan.Cycle.Help,
		IsList:     false,
		ConfigType: "Cycle",
		IsSwitch:   false,
		Speed:      "",
	})
	//cms
	data = append(data, models.ConfigExternal{
		Value:      nasConfig.ConfigNascan.DiscernCms.Value,
		Info:       nasConfig.ConfigNascan.DiscernCms.Info,
		Help:       nasConfig.ConfigNascan.DiscernCms.Help,
		IsList:     true,
		ConfigType: "Discern_cms",
		IsSwitch:   false,
		Speed:      "",
	})
	//lang
	data = append(data, models.ConfigExternal{
		Value:      nasConfig.ConfigNascan.DiscernLang.Value,
		Info:       nasConfig.ConfigNascan.DiscernLang.Info,
		Help:       nasConfig.ConfigNascan.DiscernLang.Help,
		IsList:     true,
		ConfigType: "Discern_lang",
		IsSwitch:   false,
		Speed:      "",
	})
	//con
	data = append(data, models.ConfigExternal{
		Value:      nasConfig.ConfigNascan.DiscernCon.Value,
		Info:       nasConfig.ConfigNascan.DiscernCon.Info,
		Help:       nasConfig.ConfigNascan.DiscernCon.Help,
		IsList:     true,
		ConfigType: "Discern_con",
		IsSwitch:   false,
		Speed:      "",
	})
	//scan_list
	data = append(data, models.ConfigExternal{
		Value:      nasConfig.ConfigNascan.ScanList.Value,
		Info:       nasConfig.ConfigNascan.ScanList.Info,
		Help:       nasConfig.ConfigNascan.ScanList.Help,
		IsList:     true,
		ConfigType: "Scan_list",
		IsSwitch:   false,
		Speed:      "",
	})
	//port_list
	res := strings.Split(nasConfig.ConfigNascan.PortList.Value, "|")
	data = append(data, models.ConfigExternal{
		Value:      res[1],
		Info:       nasConfig.ConfigNascan.PortList.Info,
		Help:       nasConfig.ConfigNascan.PortList.Help,
		IsList:     true,
		ConfigType: "Port_list",
		IsSwitch:   res[0] == "0",
		Speed:      "",
	})
	//server
	data = append(data, models.ConfigExternal{
		Value:      nasConfig.ConfigNascan.DiscernServer.Value,
		Info:       nasConfig.ConfigNascan.DiscernServer.Info,
		Help:       nasConfig.ConfigNascan.DiscernServer.Help,
		IsList:     true,
		ConfigType: "Discern_server",
		IsSwitch:   false,
		Speed:      "",
	})
	//white_list
	data = append(data, models.ConfigExternal{
		Value:      nasConfig.ConfigNascan.WhiteList.Value,
		Info:       nasConfig.ConfigNascan.WhiteList.Info,
		Help:       nasConfig.ConfigNascan.WhiteList.Help,
		IsList:     true,
		ConfigType: "White_list",
		IsSwitch:   false,
		Speed:      "",
	})
	return data
}
func VulscanConfigMongoToTemplate(vulConfig *models.Vulscan) []models.ConfigExternal {
	data := make([]models.ConfigExternal, 0)
	//thread
	data = append(data, models.ConfigExternal{
		Value:      vulConfig.ConfigVulscan.Thread.Value,
		Info:       vulConfig.ConfigVulscan.Thread.Info,
		Help:       vulConfig.ConfigVulscan.Thread.Help,
		IsList:     false,
		ConfigType: "Thread",
		IsSwitch:   false,
		Speed:      "",
	})
	//timeout
	data = append(data, models.ConfigExternal{
		Value:      vulConfig.ConfigVulscan.Timeout.Value,
		Info:       vulConfig.ConfigVulscan.Timeout.Info,
		Help:       vulConfig.ConfigVulscan.Timeout.Help,
		IsList:     false,
		ConfigType: "Timeout",
		IsSwitch:   false,
		Speed:      "",
	})
	//white_list
	data = append(data, models.ConfigExternal{
		Value:      vulConfig.ConfigVulscan.WhiteList.Value,
		Info:       vulConfig.ConfigVulscan.WhiteList.Info,
		Help:       vulConfig.ConfigVulscan.WhiteList.Help,
		IsList:     true,
		ConfigType: "White_list",
		IsSwitch:   false,
		Speed:      "",
	})
	//password_dic
	data = append(data, models.ConfigExternal{
		Value:      vulConfig.ConfigVulscan.PasswordDic.Value,
		Info:       vulConfig.ConfigVulscan.PasswordDic.Info,
		Help:       vulConfig.ConfigVulscan.PasswordDic.Help,
		IsList:     true,
		ConfigType: "Scan_list",
		IsSwitch:   false,
		Speed:      "",
	})
	return data
}
