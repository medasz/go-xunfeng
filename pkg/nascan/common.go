package nascan

import (
	"strings"

	"go-xunfeng/db"
	"go-xunfeng/models"
)

func FormatCommonItemConfig(value string) []models.CommonItem {
	data := make([]models.CommonItem, 0)
	lines := strings.Split(value, "\n")
	for _, row := range lines {
		res := strings.SplitN(row, "|", 4)
		if len(res) == 4 {
			data = append(data, models.CommonItem{
				Name:     strings.ToLower(res[0]),
				Location: res[1],
				Key:      res[2],
				Value:    res[3],
			})
		}
	}
	return data
}

func FormatDiscernServerConfig(value string) []models.DiscernServer {
	data := make([]models.DiscernServer, 0)
	lines := strings.Split(value, "\n")
	for _, row := range lines {
		res := strings.SplitN(row, "|", 4)
		if len(res) == 4 {
			data = append(data, models.DiscernServer{
				Name: strings.ToLower(res[0]),
				Port: res[1],
				Mode: res[2],
				Reg:  res[3],
			})
		}
	}
	return data
}

func GetNascanConfig() (*models.ConfigNascanInfo, error) {
	data := new(models.ConfigNascanInfo)
	nascanInfo, err := db.GetNascanConfig()
	if err != nil {
		return data, err
	}
	data.ScanList = nascanInfo.ConfigNascan.ScanList.Value
	data.DiscernCms = FormatCommonItemConfig(nascanInfo.ConfigNascan.DiscernCms.Value)
	data.DiscernLang = FormatCommonItemConfig(nascanInfo.ConfigNascan.DiscernLang.Value)
	data.DiscernCon = FormatCommonItemConfig(nascanInfo.ConfigNascan.DiscernCon.Value)
	data.DiscernServer = FormatDiscernServerConfig(nascanInfo.ConfigNascan.DiscernServer.Value)
	data.Masscan = nascanInfo.ConfigNascan.Masscan.Value
	data.Cycle = nascanInfo.ConfigNascan.Cycle.Value
	data.PortList = nascanInfo.ConfigNascan.PortList.Value
	data.Thread = nascanInfo.ConfigNascan.Thread.Value
	data.Timeout = nascanInfo.ConfigNascan.Timeout.Value
	data.WhiteList = nascanInfo.ConfigNascan.WhiteList.Value
	return data, nil
}
