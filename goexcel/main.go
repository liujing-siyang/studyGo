package main

import (
	"github.com/johmue/goexcelwin"
)

func main() {
	xb := goexcelwin.Book{}
	xb.CreateXLSX("./libxl/bin64/libxl.dll")

	xb.SetKey("<License Name>", "<License Key>")

	xb.SetLocale("UTF-8")
	xb.SetRgbMode(1)

	xs := xb.AddSheet("Table1", nil)

	xs.WriteStr(1, 1, "Hello!", nil)
	xs.WriteNum(1, 2, 100, nil)

	xb.Save("test.xls")
}