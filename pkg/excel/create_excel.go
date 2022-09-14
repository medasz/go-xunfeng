package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"go-xunfeng/models"
)

var headers = []string{
	"IP", "端口", "主机名", "风险等级", "漏洞描述", "插件类型", "任务名称", "时间", "扫描批次",
}

func CreateTable(list []*models.AllData, name string) (*excelize.File, error) {
	file := excelize.NewFile()
	index := file.NewSheet(name)
	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		file.SetCellValue(name, cell, header)
	}
	for i, tmp := range list {
		cell, err := excelize.CoordinatesToCellName(1, i+2)
		if err != nil {
			continue
		}
		file.SetCellValue(name, cell, tmp.Ip)

		cell, err = excelize.CoordinatesToCellName(2, i+2)
		if err != nil {
			continue
		}
		file.SetCellValue(name, cell, tmp.Port)

		cell, err = excelize.CoordinatesToCellName(3, i+2)
		if err != nil {
			continue
		}
		file.SetCellValue(name, cell, tmp.Hostname)

		cell, err = excelize.CoordinatesToCellName(4, i+2)
		if err != nil {
			continue
		}
		file.SetCellValue(name, cell, tmp.VulLevel)

		cell, err = excelize.CoordinatesToCellName(5, i+2)
		if err != nil {
			continue
		}
		file.SetCellValue(name, cell, tmp.Info)

		cell, err = excelize.CoordinatesToCellName(6, i+2)
		if err != nil {
			continue
		}
		file.SetCellValue(name, cell, tmp.VulName)

		cell, err = excelize.CoordinatesToCellName(7, i+2)
		if err != nil {
			continue
		}
		file.SetCellValue(name, cell, tmp.Title)

		cell, err = excelize.CoordinatesToCellName(8, i+2)
		if err != nil {
			continue
		}
		file.SetCellValue(name, cell, tmp.Time)

		cell, err = excelize.CoordinatesToCellName(9, i+2)
		if err != nil {
			continue
		}
		file.SetCellValue(name, cell, tmp.LastScan)

	}
	file.SetActiveSheet(index)
	return file, nil
}
