package main

func main() {
	xlsDir string = "D:/code/gowork"
	fName string = "123"
	xls := filepath.Join(xlsDir, fName)
	fmt.Println(xls)
}