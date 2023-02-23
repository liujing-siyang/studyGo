package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


// 获取表单数据的筛选条件
type FormDataReq struct {
	FdId           int              `json:"fdId"` // 表单id
	FilterItemList []FilterItemList `json:"filterItemList"`
}

type FilterItemList struct {
	FieldName      string      `json:"fieldName"`
	FieldQueryType string      `json:"fieldQueryType"`
	Value          interface{} `json:"value"`
}
// 获取的表单数据
type FormDataResponse struct {
	Flag      bool      `json:"flag"`      //	false:失败 true:成功
	ErrorCode string    `json:"errorCode"` //	S000000为成功，其他为错误编码
	ErrorMsg  string    `json:"errorMsg"`
	Data      []RowData `json:"data"`
}

type RowData struct {
	Fild            int    `json:"fiId"` // 表单实例id
	FdId            int    `json:"fdId"` // 表单id
	DataSource      string `json:"dataSource"`
	StatisticsMonth string `json:"timeLiteral_1034486175136915471"` // 业绩月
	IsDelete        bool   `json:"isDelete"`                        // 删除标记
}

func main() {
	formDataReq := FormDataReq{
		FdId:          1034506425127575588,
		FilterItemList: []FilterItemList{},
	}
	filterItemList := FilterItemList{
		FieldName:      "timeLiteral_1034508944474021911",
		FieldQueryType: "equal",
		Value:          "2022-12",
	}
	formDataReq.FilterItemList = append(formDataReq.FilterItemList, filterItemList)
	jsonData, err := json.Marshal(formDataReq)
	if err != nil {
		fmt.Println("111")
	}
	var formDataRes FormDataResponse
	var resp *http.Response
	var respData []byte
	resp, err = http.Post("", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println("2222")
	}
	defer resp.Body.Close()
	// 读取响应
	respData, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("333")
	}
	// logrus.Infoln("push response: ", string(respData))
	err = json.Unmarshal(respData, &formDataRes)
	if err != nil {
		fmt.Println("444")
	}
	fmt.Println(formDataRes)
	myjsonData, err := json.Marshal(formDataRes.Data)
	fmt.Println(string(myjsonData))
}
