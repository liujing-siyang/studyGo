// 模拟推送裸金属物理机数据
package main
 
import (
   "flag"
   "math/rand"
   "time"
 
   "github.com/prometheus/client_golang/prometheus"
   "github.com/prometheus/client_golang/prometheus/push"
)
 
func main() {
   ExamplePusher_Push()
}
 
var url string
var uuids = []string{"ss-xldtvkrd6itigznov4howgomidmc"}
var itemKeys = []string{"cpu_util", "mem_util", "disk_util_inband", "disk_read_bytes_rate", "disk_read_requests_rate",
   "disk_write_bytes_rate", "disk_write_requests_rate", "network_incoming_bytes_rate_inband", "network_outing_bytes_rate_inband"}
 
func ExamplePusher_Push() {
 
   flag.StringVar(&url, "url", "http://101.37.25.231:9091", "pgw url")
   flag.Parse()
 
   for _, uuid := range uuids {
      pusher := push.New(url, "bare_metal")
      for _, itemKey := range itemKeys {
 
         completionTime := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "mytest",
            Name: itemKey,
            Help: itemKey + "_help",
         })
 
         // 设定推送值
         completionTime.SetToCurrentTime()
         rand.Seed(time.Now().UnixNano())
         value := rand.Intn(100)
         completionTime.Set(float64(value))
 
         pusher.Collector(completionTime)
      }
 
      // 在循环外推送，避免重复的label发生覆盖
      err := pusher.Grouping("uuid", uuid).Push()
 
      if err != nil {
         panic(err)
      }
   }
}