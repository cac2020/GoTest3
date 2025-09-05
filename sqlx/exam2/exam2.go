package exam2

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

/**
*题目2：实现类型安全映射
假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
要求 ：
定义一个 Book 结构体，包含与 books 表对应的字段。
编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/

type Book struct {
	ID     int
	Title  sql.NullString
	Author sql.NullString
	Price  sql.NullInt64
}

func Test() {
	//1. 连接数据库
	db := getDB()
	defer db.Close()
	//2. 建表
	createTable(db)
	//3. 插入数据
	initData(db)
	//4. 查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中
	qryBooks(db)
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
		CREATE TABLE IF NOT EXISTS book (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT ,
			author TEXT ,
			price INTEGER
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
	db.Exec(`INSERT INTO book (title,author,price) VALUES (?,?,?)`, "Go技术指南", "张三", 30)
	db.Exec(`INSERT INTO book (title,author,price) VALUES (?,?,?)`, "Go语言圣经", "李四", 40)
	db.Exec(`INSERT INTO book (title,author,price) VALUES (?,?,?)`, "Web3.0入门", "王五", 20)
	db.Exec(`INSERT INTO book (title,author,price) VALUES (?,?,?)`, "币圈指南", "赵四", 78)
	db.Exec(`INSERT INTO book (title,price) VALUES (?,?)`, "区块链宝典", 69)
	db.Exec(`INSERT INTO book (title,author,price) VALUES (?,?,?)`, "联盟链", "赵六", 88)
	db.Exec(`INSERT INTO book (title,author,price) VALUES (?,?,?)`, "以太坊技术内幕", "王五", 95)
	db.Exec(`INSERT INTO book (title,author) VALUES (?,?)`, "比特币技术内幕", "赵四")
	fmt.Println("初始化数据成功")
}

func qryBooks(db *sqlx.DB) {
	var books []Book
	err := db.Select(&books, "SELECT * FROM book WHERE price > ?", 50)
	if err != nil {
		log.Fatal("查询失败:", err)
		panic(err)
	} else {
		fmt.Println("查询成功")
	}
	for _, book := range books {
		fmt.Println(book)
	}

}
