package main

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/spaolacci/murmur3"
)

func main() {
	client := Conn()
	var data =`{"age":18,"month":12}`
	_, err := client.Do("hset", "myhash", "test",data)
	if err != nil {
		fmt.Println("set string failed", err)
		return
	}
	resset, err := redis.String(client.Do("hget", "myhash","test"))
	if err != nil {
		fmt.Println("set string failed", err)
		return
	}
	fmt.Println(resset)
	defer client.Close()

}

func Conn() (client redis.Conn) {
	setdb := redis.DialDatabase(2)
	setPasswd := redis.DialPassword("hi7812dpa")
	client, err := redis.Dial("tcp", "172.17.2.27:6379",setdb,setPasswd)
	if err != nil {
		fmt.Println("redis connect failed,", err)
		return
	}

	fmt.Println("redis connect success")
	return
}

func test() uint32 {
	var data = `{"timeLiteral_988506922681307182":"2022-03","salesOrgId_988506922681307183":814595866286469142,"productName_988506922681307183":"金柚康康（ODS）","salesStaffTitle_988506922681307183":"朱清清(JY06214)","entryTitle_988506922685501482":"已核对","digital_988506922681307190":null,"contractId_988506922681307183":867046106645741568,"fiId":988750684670164999,"contractCode_988506922681307183":"DG-2021000055（生效）","creatTime":"2022-06-21 10:22:13","fdId":988503722490269759,"salesStaffId_988506922681307183":704272469783015424,"contractSubjectName_988506922681307183":"杭州今元标矩科技有限公司","evaluation_988506922685501479":"1,183.58","fdVersion":12,"customerName_988506922681307183":"上海生腾数据科技有限公司","productId_988506922681307183":131,"digital_988506922685501477":null,"digital_988506922685501476":null,"isDelete":false,"digital_988506922685501478":null,"updateTime":"2022-06-21 10:22:13","digital_988506922681307188":null,"salesOrgTitle_988506922681307183":"杭州今元标矩科技有限公司/产品技术部/技术部/测试部","digital_988506922681307189":null,"evaluation_988506922685501480":"1,116.59","digital_988506922681307186":4.5,"digital_988506922681307187":12,"digital_988506922681307184":34.23,"digital_988506922681307185":33,"contractSubjectId_988506922681307183":4,"customerId_988506922681307183":781510974451576832,"dataSource":null,"digital_988506922685501481":33445}`
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("结果序列化失败")
	}
	hasher := murmur3.New32()
	hasher.Write(jsonData)
	return hasher.Sum32()
}
