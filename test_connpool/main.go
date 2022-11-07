package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var dbOA *gorm.DB

func InitDB() {
	var err error
	DbUrl := "tcp://212.64.27.146:9000?database=mbi_oa_data&username=data_code&password=RL12AnHNdE6XS606&read_timeout=10&write_timeout=20" //product
	dbOA, err = gorm.Open(clickhouse.Open(DbUrl), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("clickhouse dbOa connect success")
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := dbOA.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(5)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(10)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Second * 20)
}

func main() {
	InitDB()
	sqlDB, _ := dbOA.DB()
	DBstatus := sqlDB.Stats()
	fmt.Printf("%d --- %d --- %d \n",DBstatus.OpenConnections,DBstatus.InUse,DBstatus.Idle)
	err := sqlDB.Ping()
	fmt.Println(err)
	sql := `SELECT count() from bpo_employed_info final`
	var wg sync.WaitGroup
	count := 0
	sleepcout := 0
	rand.Seed(time.Now().Unix())
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(index int) {
			r := rand.Intn(5)
			if r == 1 {
				time.Sleep(time.Second * 25)
				fmt.Printf("index is [%d]\n", index)
				sleepcout++
			}
			res, err := dbselect(sql)
			if err != nil {
				fmt.Printf("index [%d],err :%s\n", index, err)
			}
			DBstatus = sqlDB.Stats()
			fmt.Printf("%d --- %d --- %d \n",DBstatus.OpenConnections,DBstatus.InUse,DBstatus.Idle)
			fmt.Printf("res is [%d]\n", res)
			count++
			wg.Done()
		}(i)
	}
	wg.Wait()
	DBstatus = sqlDB.Stats()
	fmt.Printf("%d --- %d --- %d \n",DBstatus.OpenConnections,DBstatus.InUse,DBstatus.Idle)
	fmt.Println(count)
	fmt.Println(sleepcout)
}

func dbselect(sql string) (result int, err error) {
	err = dbOA.Raw(sql).Scan(&result).Error
	return
}
