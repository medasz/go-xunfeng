package nascan

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go-xunfeng/db"
	"go-xunfeng/models"
)

var (
	Masscan          int64
	Nachange         int64
	StatisticsAtomic atomic.Value
	ConfigInfoAtomic atomic.Value
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

func Monitor() {
	for {
		err := db.UpdateHeartbeat(bson.M{"$set": bson.M{"up_time": time.Now()}})
		if err != nil {
			log.Println(err)
		}
		statistics := StatisticsAtomic.Load().(map[string]*models.StatisticsInfo)
		if _, ok := statistics[time.Now().Format("2006-01-02")]; !ok {
			statistics[time.Now().Format("2006-01-02")] = &models.StatisticsInfo{
				Add:    0,
				Update: 0,
				Delete: 0,
			}
		}
		err = db.UpdateOrUpsertStatistic(time.Now().Format("2006-01-02"), *statistics[time.Now().Format("2006-01-02")])
		if err != nil {
			log.Println(err)
		}
		newConfig, err := GetNascanConfig()
		if err != nil {
			log.Println(err)
		}
		oldConfig := ConfigInfoAtomic.Load().(*models.ConfigNascanInfo)
		if base64.StdEncoding.EncodeToString([]byte(oldConfig.ScanList)) !=
			base64.StdEncoding.EncodeToString([]byte(newConfig.ScanList)) {
			atomic.AddInt64(&Nachange, 1)
		}
		ConfigInfoAtomic.Store(newConfig)
		time.Sleep(time.Second * 30)
	}
}

// Cruise 存活检测
func Cruise() {
	for {
		now := time.Now()
		week := now.Weekday()
		hour := now.Hour()
		if week >= time.Monday && week <= time.Friday && hour >= 9 && hour <= 18 {
			infos, err := db.GetInfoAll()
			if err != nil {
				log.Println(err)
			}
			configInfo := ConfigInfoAtomic.Load().(*models.ConfigNascanInfo)
			timeout, err := time.ParseDuration(configInfo.Timeout + "s")
			if err != nil {
				log.Println(err)
			}
			for _, info := range infos {
				for {
					if atomic.LoadInt64(&Masscan) != 1 {
						break
					}
					time.Sleep(10)
				}
				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", info.Ip, info.Port), timeout/2)
				if err != nil {
					log.Println(err)
					_, err := db.DeleteInfo(info.Ip, info.Port)
					if err != nil {
						log.Println(err)
						continue
					}
					now = time.Now()
					date := now.Format("2006-01-02")
					log.Println(fmt.Sprintf("%s:%d delete", info.Ip, info.Port))
					statistics := StatisticsAtomic.Load().(map[string]*models.StatisticsInfo)
					statistics[date].Delete += 1
					history := models.History{
						Info:    info,
						DelTime: now.UTC(),
						Type:    "delete",
					}
					_, err = db.CreateHistory(history)
					if err != nil {
						log.Println(err)
					}
				} else {
					conn.Close()
				}
			}
		}
		time.Sleep(time.Hour)
	}
}
