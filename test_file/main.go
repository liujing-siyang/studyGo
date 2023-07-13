package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type myrow struct {
	Where                       string
	Consult_depart              string
	Consult_depart_id           string
	Consult_depart_manager_name string
	Consult_depart_manager_id   string
}

func main() {
	xlsx, err := excelize.OpenFile("C:/Users/JY7188/Desktop/协商单_对应关系0705.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取第一个工作表的数据
	rows := xlsx.GetRows("Sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}

	myrows := make([]myrow, 0)
	for i, row := range rows {
		if i < 2 {
			continue
		}
		r := myrow{
			Where :row[0],
			Consult_depart:              row[2],
			Consult_depart_id:           row[4],
			Consult_depart_manager_name: row[3],
			Consult_depart_manager_id:   row[5],
		}
		myrows = append(myrows, r)
	}
	for _ ,v := range myrows{
		fmt.Printf(`select consult_depart,consult_depart_id,consult_depart_manager_name,consult_depart_manager_id 
		from yy_divided_report where consult_depart = '%s';`,v.Where)
		fmt.Println()
		fmt.Println()
		fmt.Printf(`UPDATE yy_divided_report
		SET 
			consult_depart = '%s',
			consult_depart_id = '%s',
			consult_depart_manager_name = '%s',
			consult_depart_manager_id = '%s'
		where consult_depart = '%s';`, v.Consult_depart,v.Consult_depart_id,v.Consult_depart_manager_name,v.Consult_depart_manager_id,v.Where)
		fmt.Println()
		fmt.Println()
	}
}

func test() {
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

	fmt.Printf("%+v", data)
}

// func FilepathJoin() {
// 	xlsDir string = "D:/code/gowork"
// 	fName string = "123"
// 	xls := filepath.Join(xlsDir, fName)
// 	fmt.Println(xls)
// }
