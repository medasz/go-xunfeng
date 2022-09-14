package nascan

import (
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"go-xunfeng/db"
	"go-xunfeng/models"
	"go-xunfeng/pkg/nascan"
)

func init() {
	configInfo, err := nascan.GetNascanConfig()
	if err != nil {
		panic(err)
	}
	log.Println("获取配置成功")
	nascan.ConfigInfoAtomic.Store(configInfo)
	statistic, err := db.GetStatistics()
	if err != nil {
		panic(err)
	}
	statistics := make(map[string]*models.StatisticsInfo)
	statistics[time.Now().Format("2006-01-02")] = &statistic.Info
	nascan.StatisticsAtomic.Store(statistics)
}

func Run() {
	//fmt.Println(nascan.StatisticsAtomic.Load())
	//fmt.Println(nascan.ConfigInfoAtomic.Load())
	go nascan.Monitor()
	go nascan.Cruise()
	acData := map[string]string{}
	for {
		now := time.Now()
		nowHour := now.Hour()
		nowDay := now.Day()
		nowDate := now.Format("20060102")
		configInfo := nascan.ConfigInfoAtomic.Load().(*models.ConfigNascanInfo)
		cycles := strings.SplitN(configInfo.Cycle, "|", 2)
		if len(cycles) < 2 {
			log.Printf("格式错误：%s\n", configInfo.Cycle)
			continue
		}
		cyDay, err := strconv.Atoi(cycles[0])
		if err != nil {
			log.Println(err)
			continue
		}
		acHour, err := strconv.Atoi(cycles[1])
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("扫描规则:", configInfo.Cycle)
		naChange := atomic.LoadInt64(&nascan.Nachange)
		if _, ok := acData[nowDate]; (!ok && nowDay%cyDay == 0 && nowHour == acHour) || naChange == 1 {
			atomic.AddInt64(&nascan.Nachange, -1)
			acData[nowDate] = ""
			log.Println("开始扫描")
			r := nascan.NewStart(configInfo)
			r.Run()
		}
		time.Sleep(time.Minute)
	}
}
