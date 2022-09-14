package query

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"regexp"
	"strconv"
	"strings"
)

var conditions = []string{"ip", "banner", "port", "time", "webinfo.tag", "webinfo.title", "server", "hostname"}

func QueryLogic(query string) (bson.M, error) {
	data := bson.M{}
	queryList := strings.Split(query, ";")
	if len(queryList) > 0 {
		if len(queryList) > 1 || strings.Contains(queryList[0], ":") {
			for _, val := range queryList {
				if strings.Contains(val, ":") {
					tmp := strings.SplitN(val, ":", 2)
					if len(tmp) == 2 {
						switch tmp[0] {
						case "port":
							num, err := strconv.Atoi(tmp[1])
							if err != nil {
								return data, err
							}
							data["port"] = num
						case "banner":
							zhPattern, err := regexp.Compile("[\u4e00-\u9fa5]+")
							if err != nil {
								return data, err
							}
							if zhPattern.MatchString(tmp[1]) {
								data["banner"] = map[string]interface{}{
									"$regex":   tmp[1],
									"$options": 'i',
								}
							} else {
								textQuery, err := mgoTextSplit(tmp[1])
								if err != nil {
									return data, err
								}
								data["$text"] = map[string]interface{}{
									"$search":        textQuery,
									"$caseSensitive": true,
								}
							}
						case "ip":
							data["ip"] = map[string]interface{}{
								"$regex": tmp[1],
							}
						case "server":
							data["server"] = strings.ToLower(tmp[1])
						case "title":
							data["webinfo.title"] = map[string]interface{}{
								"$regex":   tmp[1],
								"$options": 'i',
							}
						case "tag":
							data["webinfo.tag"] = strings.ToLower(tmp[1])
						case "hostname":
							data["hostname"] = map[string]interface{}{
								"$regex":   tmp[1],
								"$options": 'i',
							}
						case "all":
							filterList := make([]map[string]interface{}, 0)
							for _, i := range conditions {
								filterList = append(filterList, map[string]interface{}{
									i: map[string]interface{}{
										"$regex":   tmp[1],
										"$options": "i",
									},
								})
							}
							data["$or"] = filterList
						default:
							data[tmp[0]] = tmp[1]
						}
					}
				}
			}
		} else {
			filterList := make([]map[string]interface{}, 0)
			for _, v := range conditions {
				filterList = append(filterList, map[string]interface{}{
					v: map[string]interface{}{
						"$regex":   queryList[0],
						"$options": "i",
					},
				})
			}
			data["$or"] = filterList
		}
	}

	return data, nil
}

// split text to support mongodb $text match on a phrase
func mgoTextSplit(queryText string) (string, error) {
	sepReg, err := regexp.Compile("[`\\-=~!@#$%^&*()_+\\[\\]{};\\'\\\\:\"|<,./<>?]")
	if err != nil {
		return "", err
	}
	wordList := sepReg.Split(queryText, -1)
	tmpWordList := make([]string, 0)
	for _, tmp := range wordList {
		tmpWordList = append(tmpWordList,
			fmt.Sprintf("\"%s\"", tmp))
	}
	textQuery := strings.Join(tmpWordList, " ")
	return textQuery, nil
}
