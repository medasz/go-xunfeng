package nascan

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go-xunfeng/db"
	"go-xunfeng/models"
	"go-xunfeng/pkg/tools"
)

type Scan struct {
	ip         string
	port       string
	timeout    time.Duration
	portList   []string
	configInfo *models.ConfigNascanInfo
	banner     string
	server     string
}

func NewScan(ip string, portList []string, configInfo *models.ConfigNascanInfo) (*Scan, error) {
	tmp := &Scan{
		ip:         ip,
		portList:   portList,
		configInfo: configInfo,
	}
	timeout, err := strconv.Atoi(configInfo.Timeout)
	if err != nil {
		return nil, err
	}
	tmp.timeout = time.Duration(timeout) * time.Second
	return tmp, nil
}

func (s *Scan) run() (err error) {
	for _, port := range s.portList {
		s.server = ""
		s.banner = ""
		s.port = port
		// 端口扫描
		log.Println("开始端口扫描：", s.port)
		s.ScanPort()
		log.Println("结束端口扫描：", s.banner)
		if s.banner == "" {
			continue
		}
		//服务识别
		log.Println("开始服务识别：", s.banner)
		s.serverDiscern()
		log.Println("结束服务识别：", s.server)
		if s.server != "" {
			continue
		}
		//尝试访问web
		log.Println("开始访问web：", s.port)
		webInfo := s.tryWeb()
		if webInfo != nil {
			log.Println(fmt.Sprintf("%s:%s is web", s.ip, s.port))
			log.Println(fmt.Sprintf("%s:%s web info %s", s.ip, s.port, webInfo))
			err = db.UpdateInfoAll(s.ip, s.port, bson.M{"$set": bson.M{"banner": s.banner, "server": "web", "webinfo": webInfo,
				"time": time.Now()}})
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func (s *Scan) ip2hostname(ip string) string {
	hostnames, err := net.LookupAddr(ip)
	if err != nil {
		log.Println(err)
		return ""
	}
	return hostnames[0]
}

func (s *Scan) ScanPort() {
	port, err := strconv.Atoi(s.port)
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := net.DialTimeout("tcp", s.ip+":"+s.port, s.timeout/2)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	time.Sleep(200 * time.Millisecond)
	conn.SetReadDeadline(time.Now().Add(s.timeout / 2))
	var buf [1024]byte
	n, err := conn.Read(buf[:])
	if err != nil {
		s.banner = "NULL"
		log.Println(err)
	} else {
		s.banner = string(buf[:n])
	}
	log.Println(s.ip, s.port, "open")
	hostname := s.ip2hostname(s.ip)
	now := time.Now()
	nowStr := now.Format("2006-01-02")
	statistics := StatisticsAtomic.Load().(map[string]*models.StatisticsInfo)
	_, err = db.CreateInfo(models.Info{
		Ip:       s.ip,
		Hostname: hostname,
		Time:     now.UTC(),
		Banner:   s.banner,
		Port:     port,
	})
	if err == nil {
		statistics[nowStr].Add += 1
		StatisticsAtomic.Store(statistics)
		return
	}
	log.Println(err)
	count, err := db.CountAllInfo(bson.M{"ip": s.ip, "port": s.port})
	if err != nil {
		log.Println(err)
		return
	}
	if count > 0 && len(s.banner) > 0 {
		res := db.FindOneInfoAndDelete(bson.M{"ip": s.ip, "port": s.port, "banner": bson.M{"$ne": s.banner}})
		if res.Err() != nil {
			log.Println(res.Err())
			return
		}
		historyInfo := models.Info{}
		err = res.Decode(&historyInfo)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = db.CreateInfo(models.Info{
			Ip:       s.ip,
			Port:     port,
			Hostname: hostname,
			Banner:   s.banner,
			Time:     now.UTC(),
		})
		if err != nil {
			log.Println(err)
			return
		}
		statistics[nowStr].Update += 1
		StatisticsAtomic.Store(statistics)
		_, err = db.CreateHistory(models.History{
			Info:    historyInfo,
			DelTime: now,
			Type:    "update",
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Scan) serverDiscern() {
	// 快速识别
	for _, markInfo := range s.configInfo.DiscernServer {
		if markInfo.Mode == "default" {
			if markInfo.Port == s.port {
				s.server = markInfo.Name
			}
		} else if markInfo.Mode == "banner" {
			reg, err := regexp.Compile("(?im)" + markInfo.Reg)
			if err != nil {
				log.Println(err)
				continue
			}
			if reg.MatchString(s.banner) {
				s.server = markInfo.Name
			}
		}
		if s.server != "" {
			break
		}
	}
	if s.server == "" && !tools.In(s.port, []string{"80", "443", "8080"}) {
		for _, markInfo := range s.configInfo.DiscernServer {
			if !tools.In(markInfo.Mode, []string{"default", "banner"}) {
				conn, err := net.DialTimeout("tcp", s.ip+":"+s.port, s.timeout/2)
				if err != nil {
					log.Println(err)
					continue
				}
				mode, err := strconv.Unquote(`"` + markInfo.Mode + `"`)
				if err != nil {
					log.Println(err)
					continue
				}
				reg, err := strconv.Unquote(`"` + markInfo.Reg + `"`)
				if err != nil {
					log.Println(err)
					continue
				}
				conn.SetDeadline(time.Now().Add(s.timeout / 2))
				conn.Write([]byte(mode))
				buf := make([]byte, 1024)
				_, err = conn.Read(buf)
				if err != nil {
					log.Println(err)
					continue
				}
				fmt.Println("-------mode-------", mode)
				fmt.Println("-------reg-------", reg)
				conn.Close()
				regx, err := regexp.Compile("(?im)" + reg)
				if err != nil {
					log.Println(err)
					continue
				}
				if regx.MatchString(string(buf)) {
					s.server = markInfo.Name
					break
				}
			}
		}
	}
	if s.server != "" {
		log.Println(fmt.Sprintf("%s:%s is %s", s.ip, s.port, s.server))
		if err := db.UpdateInfo(s.ip, s.port, s.server); err != nil {
			log.Println(err)
		}
	}
}

func (s *Scan) tryWeb() *models.WebInfo {
	port, err := strconv.Atoi(s.port)
	if err != nil {
		log.Println(err)
		return nil
	}
	var data *models.WebInfo
	url := "http://" + s.ip + ":" + s.port
	if s.port == "443" {
		url = "https://" + s.ip + ":" + s.port
	}
	client := &http.Client{
		Timeout:   s.timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	header := resp.Header
	resp.Body.Close()
	if len(header) == 0 {
		return nil
	}
	if header.Get("Content-Encoding") != "" && tools.In("gzip", header.Values("Content-Encoding")) {
		gr, err := gzip.NewReader(resp.Body) //初始化gzip reader
		defer gr.Close()                     //函数退出时关闭reader防止内存泄漏
		if err != nil {
			log.Println(err)
			return nil
		}
		html, err = io.ReadAll(gr)
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	fmt.Println(header)
	fmt.Println(string(html))
	htmlCode := s.getCode(string(html), header)
	if len(htmlCode) > 0 && len(htmlCode) < 12 {
		reader, err := charset.NewReaderLabel(htmlCode, strings.NewReader(string(html)))
		if err != nil {
			log.Println(err)
		} else {
			html, err = io.ReadAll(reader)
			if err != nil {
				log.Println(err)
			}
		}
	}
	titleReg, err := regexp.Compile(`(?im)<title>(.*?)</title>`)
	if err != nil {
		log.Println(err)
		return nil
	}
	titleGroup := titleReg.FindStringSubmatch(string(html))
	if len(titleGroup) >= 2 {
		data.Title = titleGroup[1]
	}
	var webBanner string
	for k, v := range header {
		webBanner += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	webBanner += "\r\n\r\n" + string(html)
	s.banner = webBanner
	info, err := db.GetInfoOne(bson.M{
		"ip":   s.ip,
		"port": port,
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	if info.Server == "" {
		tags := s.getTag()
		data.Tag = tags
	} else {
		if int(math.Abs(float64(len([]byte(info.Banner))))) > len([]byte(webBanner))/60 {
			now := time.Now()
			nowStr := now.Format("2006-01-02")
			historyInfo := models.History{
				Info:    *info,
				DelTime: now,
				Type:    "update",
			}
			_, err = db.CreateHistory(historyInfo)
			if err != nil {
				log.Println(err)
			}
			data.Tag = s.getTag()
			statistics := StatisticsAtomic.Load().(map[string]*models.StatisticsInfo)
			statistics[nowStr].Update += 1
			log.Println(fmt.Sprintf("%s:%s update web info", s.ip, s.port))
		}
	}
	return data
}

func (s *Scan) getCode(html string, header http.Header) string {
	reg, err := regexp.Compile(`(?i)<meta.*?charset=(.*?)"(>| |/)`)
	if err != nil {
		log.Println(err)
		return ""
	}
	m := reg.FindStringSubmatch(html)
	if len(m) >= 2 {
		return strings.TrimLeft(m[1], "\"")
	}
	if header.Get("Content-Type") != "" {
		contentType := header.Get("Content-Type")
		contentTypeReg, err := regexp.Compile(`(?i).*?charset=(.*?)(;|$)`)
		if err != nil {
			log.Println(err)
			return ""
		}
		contentTypeGroup := contentTypeReg.FindStringSubmatch(contentType)
		if len(contentTypeGroup) >= 2 {
			return contentTypeGroup[1]
		}
	}
	return ""
}

func (s *Scan) getTag() (tags []string) {
	urlPath := "http://" + net.JoinHostPort(s.ip, s.port)
	if s.port == "443" {
		urlPath = "https://" + net.JoinHostPort(s.ip, s.port)
	}
	client := http.Client{
		Timeout:   s.timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Get(urlPath)
	if err != nil {
		log.Println(err)
		return
	}
	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	html := string(htmlBytes)
	header := resp.Header
	resp.Body.Close()
	for _, disType := range []string{"Discern_cms", "Discern_con", "Discern_lang"} {
		tag := s.discern(header, html, disType, urlPath)
		if tag != "" && !tools.In(tag, tags) {
			tags = append(tags, tag)
		}
	}
	return
}

func (s *Scan) getDisType(disType string) []models.CommonItem {
	data := make([]models.CommonItem, 0)
	switch disType {
	case "Discern_cms":
		data = s.configInfo.DiscernCms
	case "Discern_con":
		data = s.configInfo.DiscernCon
	case "Discern_lang":
		data = s.configInfo.DiscernLang
	}
	return data
}

func (s *Scan) discern(header http.Header, html, disType, urlPath string) string {
	fileTmp := make(map[string]string)
	for _, markInfo := range s.getDisType(disType) {
		switch markInfo.Location {
		case "header":
			if len(header) == 0 {
				return ""
			}
			headerValue := header.Get(markInfo.Key)
			reg, err := regexp.Compile(`(?i)` + markInfo.Value)
			if err != nil {
				log.Println(err)
				continue
			}
			if reg.MatchString(headerValue) {
				return markInfo.Name
			}
		case "file":
			if markInfo.Key == "index" {
				if html == "" {
					return ""
				}
				reg, err := regexp.Compile(`(?i)` + markInfo.Value)
				if err != nil {
					log.Println(err)
					continue
				}
				if reg.MatchString(html) {
					return markInfo.Name
				}
			} else {
				var reHtml string
				if _, ok := fileTmp[markInfo.Key]; ok {
					reHtml = fileTmp[markInfo.Key]
				} else {
					client := http.Client{
						Timeout:   s.timeout,
						Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
					}
					resp, err := client.Get(urlPath)
					if err != nil {
						log.Println(err)
						return ""
					}
					htmlBytes, err := io.ReadAll(resp.Body)
					if err != nil {
						log.Println(err)
						return ""
					}
					reHtml = string(htmlBytes)
					fileTmp[markInfo.Key] = reHtml
				}
				reg, err := regexp.Compile(`(?i)` + markInfo.Value)
				if err != nil {
					log.Println(err)
				}
				if reg.MatchString(reHtml) {
					return markInfo.Name
				}
			}
		}
	}
	return ""
}
