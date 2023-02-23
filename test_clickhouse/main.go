package main

import (
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	// "gorm.io/driver/clickhouse"
	// "gorm.io/gorm"
	// "gorm.io/gorm/schema"
)

type VisitRankingOrg struct {
	SignBranchOrgId   int64  `gorm:"column:sign_branch_org_id"`   // 分公司id
	SignBranchOrgName string `gorm:"column:sign_branch_org_name"` // 分公司名称
	SignOrgId         int64  `gorm:"column:sign_org_id"`          // 小组id
	SignOrgTitle      string `gorm:"column:sign_org_title"`       // 小组名称
}

func main() {
	clickhouse_chproxy()
}

func clickhouseCon() {
	//dsn := "tcp://212.64.27.146:9000?database=smarthr-crm&username=data_code&password=RL12AnHNdE6XS606&read_timeout=10&write_timeout=20"
	// dsn := "tcp://172.17.2.5:9000?database=YGBI&username=default&password=RL12AnHNdE6XS606&read_timeout=10&write_timeout=20"
	// db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{
	// 	NamingStrategy: schema.NamingStrategy{
	// 		SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
	// 	},
	// })
	// if err != nil {
	// 	fmt.Println("connect fail")
	// } else {
	// 	fmt.Println("connect success")
	// 	fmt.Println(db)
	// }
	// sql := `SELECT o4.id sign_branch_org_id, o4.org_title sign_branch_org_name, o5.id sign_org_id, o5.org_title sign_org_title
	// 		FROM org o3  -- 区
	// 		LEFT JOIN org o4 ON o3.id = o4.pid -- 分公司
	// 		LEFT JOIN org o5 ON o4.id = o5.pid -- 小组
	// 		WHERE o3.pid = ?
	// 		AND o5.org_title LIKE '销售%'
	// 		AND o5.status = 1`
	// var list []VisitRankingOrg
	// err = db.Table("org").Raw(sql, 727451181164638208).Scan(&list).Error
	// if err != nil {
	// 	return
	// }
	// for index, value := range list {
	// 	fmt.Println(index, value)
	// }

	// type Org struct {
	// 	BranchId   int64  `gorm:"column:new_sign_branch_org_id"`
	// 	BranchName string `gorm:"column:new_sign_branch_org_name"`
	// }
	// var orgList []Org
	// var branchIds []int64
	// var branchNames []string
	// err = db.Table("org_three_level").Select("new_sign_branch_org_id, new_sign_branch_org_name").
	// 	Where("new_department_id = ?", 727451181164638208).
	// 	Where("new_department_id IS NOT NULL").
	// 	Group("new_sign_branch_org_id, new_sign_branch_org_name").Scan(&orgList).
	// 	Error
	// if err != nil {
	// 	return
	// }
	// for _, value := range orgList {
	// 	branchIds = append(branchIds, value.BranchId)
	// 	branchNames = append(branchNames, value.BranchName)
	// }
	// for index, value := range branchIds {
	// 	fmt.Println(index, value)
	// }
	// for index, value := range branchNames {
	// 	fmt.Println(index, value)
	// }
	// connect, err := sql.Open("clickhouse", "tcp://172.17.2.5:9000?debug=true")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := connect.Ping(); err != nil {
	// 	fmt.Println("connect fail")
	// 	return
	// }
	// connect.Close()
}

func clickhouse_chproxy() {
	err := Connect()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("connect success")
	}
}

func Connect() error {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", "212.64.27.146", 9080)},
		Auth: clickhouse.Auth{
			Database: "test1_smarthr_finance",
			Username: "dev",
			Password: "dev123",
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		// DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Protocol: clickhouse.HTTP, //使用chproxy指定HTTP协议
	})
	err := conn.Ping()
	if err != nil {
		return err
	}
	sql := `
		select id from bill_invoice bi limit 10
	`
	rows, err := conn.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		fmt.Println(id)
	}
	return err
}

func ConnectDSN() error {
	conn, err := sql.Open("clickhouse", fmt.Sprintf("clickhouse://%s:%d?username=%s&password=%s", "212.64.27.146", 9080, "dev", "dev123"))
	if err != nil {
		return err
	}
	return conn.Ping()
}
