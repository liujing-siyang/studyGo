package main

import (
	"context"
	"fmt"
	"strings"

	"log"

	"github.com/360EntSecGroup-Skylar/excelize"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type myrow struct {
	Where                       string
	Consult_depart              string
	Consult_depart_id           string
	Consult_depart_manager_name string
	Consult_depart_manager_id   string
}

func test2() {
	xlsx, err := excelize.OpenFile("C:/Users/JY7188/Desktop/协商单_对应关系0905.xlsx")
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

	myrows := make([]myrow, 0)
	for i, row := range rows {
		if i < 1 {
			continue
		}
		r := myrow{
			Where:                       row[0],
			Consult_depart:              row[1],
			Consult_depart_id:           row[2],
			Consult_depart_manager_name: row[3],
			Consult_depart_manager_id:   row[4],
		}
		myrows = append(myrows, r)
	}
	for _, v := range myrows {
		fmt.Printf(`select ownership,ownership_id,ownership_manager_name,ownership_manager_id 
		from yy_divided_report where ownership = '%s';`, v.Where)
		fmt.Println()
		fmt.Println()
		fmt.Printf(`UPDATE yy_divided_report
		SET 
			consult_depart = '%s',
			consult_depart_id = '%s',
			consult_depart_manager_name = '%s',
			consult_depart_manager_id = '%s'
		where consult_depart = '%s';`, v.Consult_depart, v.Consult_depart_id, v.Consult_depart_manager_name, v.Consult_depart_manager_id, v.Where)
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

// 运营体系协商明细表
type YyDividedReport struct {
	Id                       int64   `gorm:"column:id" json:"id"`
	StatisticsMonth          string  `gorm:"column:statistics_month" json:"statistics_month"`                       //业绩月
	CustCn                   string  `gorm:"column:cust_cn" json:"cust_cn"`                                         //客户名称
	ComboTitle               string  `gorm:"column:combo_title" json:"combo_title"`                                 //签单产品
	ContCode                 string  `gorm:"column:cont_code" json:"cont_code"`                                     //合同编码
	Function                 string  `gorm:"column:function" json:"function"`                                       //协商职能
	ConsultDepart            string  `gorm:"column:consult_depart" json:"consult_depart"`                           //协商部门
	ConsultDepartId          string  `gorm:"column:consult_depart_id" json:"consult_depart_id"`                     //协商部门id
	SellerUserNameCode       string  `gorm:"column:seller_user_name_code" json:"seller_user_name_code"`             //所属销售
	SellerUserStaffId        string  `gorm:"column:seller_user_staff_id" json:"seller_user_staff_id"`               //所属销售staff_id
	SellerUserOrgId          string  `gorm:"column:seller_user_org_id" json:"seller_user_org_id"`                   //所属销售org_id
	ProjectLeader            string  `gorm:"column:project_leader" json:"project_leader"`                           //项目负责人
	ProjectLeaderStaffId     string  `gorm:"column:project_leader_staff_id" json:"project_leader_staff_id"`         //项目负责人staff_id
	ProjectLeaderOrgId       string  `gorm:"column:project_leader_org_id" json:"project_leader_org_id"`             //项目负责人org_id
	PerformanceDivideRate    float64 `gorm:"column:performance_divide_rate" json:"performance_divide_rate"`         //业绩分成比例
	DepartCommissionRate     float64 `gorm:"column:depart_commission_rate" json:"depart_commission_rate"`           //部门提成比列
	ConsultDepartManagerName string  `gorm:"column:consult_depart_manager_name" json:"consult_depart_manager_name"` //协商部门负责人
	ConsultDepartManagerId   string  `gorm:"column:consult_depart_manager_id" json:"consult_depart_manager_id"`     //协商部门负责人staff_id
	Ownership                string  `gorm:"column:ownership" json:"ownership"`                                     //业绩归属部门
	OwnershipId              string  `gorm:"column:ownership_id" json:"ownership_id"`                               //业绩归属部门id
	OwnershipManagerName     string  `gorm:"column:ownership_manager_name" json:"ownership_manager_name"`           //业绩归属部门负责人
	OwnershipManagerId       string  `gorm:"column:ownership_manager_id" json:"ownership_manager_id"`               //业绩归属部门负责人staff_id
	Income                   float64 `gorm:"column:income" json:"income"`                                           //营业收入
	Cost                     float64 `gorm:"column:cost" json:"cost"`                                               //成本
	ProjectGrossProfit       float64 `gorm:"column:project_gross_profit" json:"project_gross_profit"`               //项目毛利
	CommissionBase           float64 `gorm:"column:commission_base" json:"commission_base"`                         //提成基数
	DividedRate              float64 `gorm:"column:divided_rate" json:"divided_rate"`                               //分成比例
	OverAdvancePayment       float64 `gorm:"column:over_advance_payment" json:"over_advance_payment"`               //超45天垫付费用
	IsCommissionCom          string  `gorm:"column:is_commission_com" json:"is_commission_com"`                     //提成是否计算
	CommissionFrozen         string  `gorm:"column:commission_frozen" json:"commission_frozen"`                     //提成是否冻结
	Helper                   string  `gorm:"column:helper" json:"helper"`                                           //协助人
	HelperOrgName            string  `gorm:"column:helper_org_name" json:"helper_org_name"`                         //协助人部门
	HelperBranchOrgId        string  `gorm:"column:helper_branch_org_id" json:"helper_branch_org_id"`               //协助人部门orgid
	HelperStaffId            string  `gorm:"column:helper_staff_id" json:"helper_staff_id"`                         //协助人staff_id
	HelperOrgId              string  `gorm:"column:helper_org_id" json:"helper_org_id"`                             //协助人org_id
	ConsultIncome            float64 `gorm:"column:consult_income" json:"consult_income"`                           //协商营收
	ConsultCost              float64 `gorm:"column:consult_cost" json:"consult_cost"`                               //协商成本
	ConsultGrossProfit       float64 `gorm:"column:consult_gross_profit" json:"consult_gross_profit"`               //协商毛利
	ConsultCommissionBase    float64 `gorm:"column:consult_commission_base" json:"consult_commission_base"`         //协商提成基数
	ConsultCommissionAccrual float64 `gorm:"column:consult_commission_accrual" json:"consult_commission_accrual"`   //协商提成计提
	ConsultCompanyCost       float64 `gorm:"column:consult_company_cost" json:"consult_company_cost"`               //协商经费
	CommissionDivided        float64 `gorm:"column:commission_divided" json:"commission_divided"`                   //提成分成
	CodeRemark               string  `gorm:"column:code_remark" json:"code_remark"`                                 //合同备注
	IsUpdate                 int8    `gorm:"column:is_update" json:"is_update"`                                     //是否更新，1：更新
}

func GetListbyStaff() (list []*YyDividedReport, err error) {
	sql := `
		select * from yy_divided_report 
		where statistics_month >= '202301' and statistics_month <= '202307'
	`
	err = YjbbMysql.Raw(sql).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return
}

func GetListbyPro() (list []*YyDividedReport, err error) {
	sql := `
	select * from yy_divided_report 
	where statistics_month >= '202301' and statistics_month <= '202307'
	and locate("项目",consult_depart)
	`
	err = YjbbMysql.Raw(sql).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return
}

var (
	YjbbMysql *gorm.DB
	mysqlurl  = "yjbb:lmj89tyewdj$hui45@tcp(sh-cdb-5gqdlfey.sql.tencentcdb.com:61787)/yjbb_produce?charset=utf8&parseTime=true"
)

func init() {
	yjbbMysql, err := gorm.Open(mysql.Open(mysqlurl), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	YjbbMysql = yjbbMysql
}

func UpadteInSelectDB(db *gorm.DB, t interface{}, update map[string]interface{}, where string, args ...interface{}) (err error) {
	tx := db.Begin()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			tx.Rollback()
			err = fmt.Errorf("捕获异常: %v", panicErr)
		}
	}()
	err = Update(context.TODO(), tx, &t, where, update, args...)
	if err != nil {
		return
	}
	tx.Commit()
	return nil
}

func Update(ctx context.Context, sourceDB *gorm.DB, m interface{}, where string, update map[string]interface{}, args ...interface{}) error {
	err := sourceDB.Model(m).Where(where, args...).
		Updates(update).Error
	if err != nil {
		sourceDB.Rollback()
		return fmt.Errorf("ConsultClient:Update where=%s: %w", where, err)
	}
	return nil
}

func InitExcel() {
	xlsx, err := excelize.OpenFile("C:/Users/JY7188/Desktop/2023-08科技运营&营销体系组织架构调整对应关系.xlsx")
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
	staffOrgMap := make(map[string]string, 0)
	for i, v := range rows {
		if i == 0 {
			continue
		}
		staffOrgMap[v[1]] = v[14]
	}
	list, err := GetListbyStaff()
	if err != nil {
		return
	}
	for _, v := range list {
		org, ok := staffOrgMap[v.SellerUserStaffId]
		if ok {
			err = UpadteInSelectDB(YjbbMysql, YyDividedReport{}, map[string]interface{}{
				"seller_user_org_id": org,
			}, "id  = ?", v.Id)
			if err != nil {
				log.Print(err.Error())
				return
			}
		}
	}
}

func InitExcel2() {
	xlsx, err := excelize.OpenFile("C:/Users/JY7188/Desktop/2023-08科技运营&营销体系组织架构调整对应关系.xlsx")
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
	staffOrgMap := make(map[string][]string, 0)
	for i, v := range rows {
		if i == 0 {
			continue
		}
		staffOrgMap[v[1]] = []string{v[14], strings.Join([]string{v[3], v[5], v[7], v[9]}, "/"), v[10]}
	}
	secondxlsx, err := excelize.OpenFile("C:/Users/JY7188/Desktop/2023协商单列表.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取第一个工作表的数据
	secondrows := secondxlsx.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	monthCodeMap := make(map[string]struct{})
	for i,v := range secondrows{
		if i == 0 {
			continue
		}
		monthCodeMap[v[0]+v[1]] = struct{}{}
	}
	list, err := GetListbyPro()
	if err != nil {
		return
	}
	for _, v := range list {
		// 过滤不需要更新的数据
		if _,ok := monthCodeMap[v.StatisticsMonth+v.ContCode];ok{
			continue
		}
		org, ok := staffOrgMap[v.ProjectLeaderStaffId]
		if ok {
			err = UpadteInSelectDB(YjbbMysql, YyDividedReport{}, map[string]interface{}{
				"project_leader_org_id": org[0],
				"consult_depart":        org[1],
				"consult_depart_id":     org[2],
			}, "id  = ?", v.Id)
			if err != nil {
				log.Print(err.Error())
				return
			}
		}
	}
}

type temp struct{
	Where                       string
	Consult_depart_manager_name string
	Consult_depart_manager_id   string
}

func test3() {
	xlsx, err := excelize.OpenFile("C:/Users/JY7188/Desktop/123.xlsx")
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

	myrows := make([]myrow, 0)
	for i, row := range rows {
		if i < 1 {
			continue
		}
		r := myrow{
			Where:                       row[0],
			Consult_depart_manager_name: row[1],
			Consult_depart_manager_id:   row[2],
		}
		myrows = append(myrows, r)
	}
	for _, v := range myrows {
		fmt.Printf(`select consult_depart_manager_name,consult_depart_manager_id 
		from yy_divided_report where consult_depart = '%s';`, v.Where)
		fmt.Println()
		fmt.Println()
		fmt.Printf(`UPDATE yy_divided_report
		SET 
			consult_depart_manager_name = '%s',
			consult_depart_manager_id = '%s'
		where consult_depart = '%s';`, v.Consult_depart_manager_name, v.Consult_depart_manager_id, v.Where)
		fmt.Println()
		fmt.Println()
	}
}

func InitExcel4() {
	xlsx, err := excelize.OpenFile("C:/Users/JY7188/Desktop/2023-08科技运营&营销体系组织架构调整对应关系.xlsx")
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
	staffOrgMap := make(map[string][]string, 0)
	for i, v := range rows {
		if i == 0 {
			continue
		}
		staffOrgMap[v[1]] = []string{v[14],v[13],v[10]}
	}
	list, err := GetListbyStaff()
	if err != nil {
		return
	}
	for _, v := range list {
		org, ok := staffOrgMap[v.HelperStaffId]
		if ok {
			err = UpadteInSelectDB(YjbbMysql, YyDividedReport{}, map[string]interface{}{
				"helper_org_id": org[0],
				"helper_org_name":org[1],
				"helper_branch_org_id":org[2],
			}, "id  = ?", v.Id)
			if err != nil {
				log.Print(err.Error())
				return
			}
		}
	}
}



func main() {
	// test2()
	// InitExcel()
	// InitExcel2()
	// test3()
	InitExcel4()
}
