package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	// 读取 Excel 文件
	xlsx, err := excelize.OpenFile("C:/Users/JY7188/Desktop/data.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取第一个工作表的数据
	rows := xlsx.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 生成 map
	data := make(map[string][]string)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		key := row[0]
		var values []string
		for _, cell := range row[1:] {
			values = append(values, cell)
		}
		data[key] = values
	}

	fmt.Printf("%+v",data)
}





// func FilepathJoin() {
// 	xlsDir string = "D:/code/gowork"
// 	fName string = "123"
// 	xls := filepath.Join(xlsDir, fName)
// 	fmt.Println(xls)
// }


