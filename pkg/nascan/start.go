package nascan

import (
	"go-xunfeng/models"
	"go-xunfeng/pkg/slices"
	"go-xunfeng/pkg/tools"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Start struct {
	ConfigInfo *models.ConfigNascanInfo
}

func NewStart(configInfo *models.ConfigNascanInfo) *Start {
	return &Start{
		configInfo,
	}
}

func (s Start) Run() {
	scanList := strings.Split(s.ConfigInfo.ScanList, "\n")
	whiteList := strings.Split(s.ConfigInfo.WhiteList, "\n")
	mode := strings.Split(s.ConfigInfo.Masscan, "|")[0]
	icmpOn := strings.Split(s.ConfigInfo.PortList, "|")[0]
	var allIpList []string
	for _, scan := range scanList {
		// 解析IP段
		ips, err := Handler(scan)
		if err != nil {
			log.Println(err)
			continue
		}
		// 清理白名单IP
		for _, white := range whiteList {
			ips = slices.DeleteElement(ips, white)
		}
		// 开始扫描IP
		if mode == "1" {
			// 开启masscan扫描
			//masscanPath := strings.Split(s.ConfigInfo.Masscan, "|")[2]
			//masscanRate := strings.Split(s.ConfigInfo.Masscan, "|")[1]
			if icmpOn == "1" {
				// icmp探测
				ips = s.getAcIps(ips)
			}
			// masscan扫描
			atomic.AddInt64(&Masscan, 1)
			ipPortMap := s.masscan(ips)
			atomic.AddInt64(&Masscan, -1)
			if len(ipPortMap) == 0 {
				continue
			}
			s.scanStart(ips, ipPortMap)
		} else {
			allIpList = append(allIpList, ips...)
		}
	}
	if mode == "0" {
		if icmpOn == "1" {
			allIpList = s.getAcIps(allIpList)
		}
		s.scanStart(allIpList, nil)
	}
}

func (s Start) scanStart(ips []string, ipPortsMap map[string][]string) {
	threat, err := strconv.Atoi(s.ConfigInfo.Thread)
	if err != nil {
		log.Println(err)
		return
	}
	mode := strings.Split(s.ConfigInfo.Masscan, "|")[0]
	queue := make(chan string, 200)
	go func() {
		for _, ip := range ips {
			queue <- ip
		}
		close(queue)
	}()
	wg := &sync.WaitGroup{}
	wg.Add(threat)
	for i := 0; i < threat; i++ {
		go func() {
			defer func() {
				wg.Done()
				if err := recover(); err != nil {
					log.Println(err)
					return
				}
			}()
			for ip := range queue {
				var portList []string
				if mode == "1" {
					portList = ipPortsMap[ip]
				} else {
					portList = strings.Split(strings.Split(s.ConfigInfo.PortList, "|")[1], "\n")
				}
				scan, err := NewScan(ip, portList, s.ConfigInfo)
				if err != nil {
					log.Println(err)
					continue
				}
				err = scan.run()
				if err != nil {
					log.Println(err)
				}
			}
		}()
	}
	wg.Wait()
	time.Sleep(time.Hour)
}

func (s Start) masscan(ips []string) map[string][]string {
	data := make(map[string][]string)
	if len(ips) == 0 {
		return data
	}
	masscanInfo := strings.Split(s.ConfigInfo.Masscan, "|")
	if len(masscanInfo) < 3 {
		log.Println("masscanInfo is err:", masscanInfo)
		return data
	}
	m := NewMasscaner()
	m.SetSystemPath(masscanInfo[2])
	m.SetPorts("1-65535")
	m.SetRate(masscanInfo[1])
	m.SetArgs("-iL", "target.log", "--randomize-hosts")
	// 创建文件，写入扫描目标
	ipsText := strings.Join(ips, "\n")
	err := tools.OpenOrCreateTargetFile(ipsText)
	if err != nil {
		log.Println(err)
		return data
	}
	// masscan扫描
	log.Println("masscan扫描开始")
	err = m.Run()
	if err != nil {
		log.Println(err)
		return data
	}
	// 删除扫描目标文件
	os.Remove("target.log")
	log.Println("masscan扫描结束")
	// 解析扫描结果
	results, err := m.Parse()
	if err != nil {
		log.Println(err)
		return data
	}
	for _, result := range results {
		if v, ok := data[result.Address.Addr]; ok {
			for _, port := range result.Ports {
				if !tools.In(port.Portid, v) {
					data[result.Address.Addr] = append(data[result.Address.Addr], port.Portid)
				}
			}
		} else {
			data[result.Address.Addr] = result.Ports.ToPortList()
		}
	}
	return data
}

/*
go实现icmp探测参考链接：
https://colobu.com/2023/04/26/write-the-ping-tool-in-Go/
https://colobu.com/2024/05/20/implemenmt-pping-in-go/
https://colobu.com/2023/09/10/mping-a-multi-targets-high-frequency-pressure-measuring-and-detection-tool/
https://github.com/mdlayher/icmpx
https://github.com/SayHe110/easy-ping/blob/master/main.go
https://github.com/go-ping/ping/blob/master/ping.go#L704
https://www.ctyun.cn/developer/article/408613428564037
https://colobu.com/2019/06/01/packet-capture-injection-and-analysis-gopacket/
https://github.com/EvanLi/programming-book-2/tree/master/Go
*/
func (s Start) getAcIps(ips []string) []string {
	nas, err := NewNascan()
	if err != nil {
		return []string{}
	}
	return nas.mPing(ips)
}
