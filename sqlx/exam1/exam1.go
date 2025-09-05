package exam1

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

/**
 * 题目1：使用SQL扩展库进行查询
假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
*/

type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     int
}

func Test() {
	//1. 连接数据库
	db := getDB()
	defer db.Close()
	//2. 建表
	createTable(db)
	//3. 插入数据
	initData(db)
	//4. 查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	qryTechEmps(db)
	//5. 查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
	qryHighestSalary(db)
}

func getDB() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("创建数据库连接失败:", err)
		panic(err)
	} else {
		fmt.Println("创建数据库连接成功")
	}

	return db
}
func createTable(db *sqlx.DB) {
	sql := `
		CREATE TABLE IF NOT EXISTS employees (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT ,
			department TEXT ,
			salary INTEGER
		);
	`
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal("创建表失败:", err)
		panic(err)
	} else {
		fmt.Println("创建表成功")
	}
}

func initData(db *sqlx.DB) {
	insSql := `INSERT INTO employees (name,department,salary) VALUES (?,?,?)`
	db.Exec(insSql, "张三", "技术部", 5000)
	db.Exec(insSql, "李四", "技术部", 1000)
	db.Exec(insSql, "王五", "产品部", 800)
	db.Exec(insSql, "赵六", "产品部", 900)
	db.Exec(insSql, "孙七", "销售部", 700)
	db.Exec(insSql, "周八", "销售部", 800)
	db.Exec(insSql, "吴九", "技术部", 900)
	db.Exec(insSql, "郑十", "技术部", 1000)
	fmt.Println("初始化数据成功")
}

func qryTechEmps(db *sqlx.DB) {
	employees := []Employee{}
	err := db.Select(&employees, "SELECT * FROM employees WHERE department = ?", "技术部")
	if err != nil {
		log.Fatal("查询数据失败:", err)
		panic(err)
	}
	fmt.Println("技术部人员：", employees)
}

func qryHighestSalary(db *sqlx.DB) {
	employee := Employee{}
	err := db.Get(&employee, "SELECT * FROM employees ORDER BY salary DESC LIMIT 1")
	if err != nil {
		log.Fatal("查询数据失败:", err)
		panic(err)
	}
	fmt.Println("最高工资员工：", employee)
}
