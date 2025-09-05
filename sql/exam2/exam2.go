package exam2

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

/*
*
  假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
	要求 ：
	- 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/

func TransferTest() {
	//1.创建数据库连接
	db := getDB()
	defer db.Close()

	//2.建表 并插入若干条数据
	createTable(db)

	//3.插入数据
	initData(db)

	//4.转账
	transfer(db, "A", "B", 100)

	//5.转账后余额查询
	queryAccounts(db)
}

func getDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("创建数据库连接失败:", err)
		panic(err)
	} else {
		fmt.Println("创建数据库连接成功")
	}

	return db
}

func createTable(db *sql.DB) {
	schemaTransactions := `
    CREATE TABLE IF NOT EXISTS transactions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_account_id TEXT NOT NULL,
		to_account_id TEXT NOT NULL,
        amount INTEGER NOT NULL
    );`
	_, err1 := db.Exec(schemaTransactions)
	if err1 != nil {
		log.Fatal("创建账户表失败:", err1)
		panic(err1)
	} else {
		fmt.Println("创建创建账户表成功")
	}

	schemaAccounts := `
    CREATE TABLE IF NOT EXISTS accounts (
        id TEXT PRIMARY KEY,
        balance INTEGER NOT NULL
    );`
	_, err2 := db.Exec(schemaAccounts)
	if err2 != nil {
		log.Fatal("创建交易表失败:", err2)
		panic(err2)
	} else {
		fmt.Println("创建创建交易表成功")
	}
}

func initData(db *sql.DB) {
	inSql := `INSERT INTO accounts (id,balance) VALUES (?,?)`
	db.Exec(inSql, "A", 1000)
	db.Exec(inSql, "B", 800)
	fmt.Println("初始化数据成功")
}

func transfer(db *sql.DB, from, to string, amount int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("创建事务失败:", err)
		return
	}

	//检查余额
	var balance int
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", from).Scan(&balance)
	if err != nil {
		tx.Rollback()
		log.Fatal("查询账户"+from+"余额失败，转账失败，交易回滚", err)
		return
	}
	fmt.Println("账户"+from+"转账前余额为：", balance)
	if balance < amount {
		tx.Rollback()
		log.Fatal("账户" + from + "余额不足，转账失败，交易回滚")
		return
	}

	//执行转账
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, from)
	if err != nil {
		tx.Rollback()
		log.Fatal("转账失败，账户"+from+"扣款失败，交易回滚", err)
		return
	}
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, to)
	if err != nil {
		tx.Rollback()
		log.Fatal("转账失败，账户"+to+"加款失败，交易回滚", err)
		return
	}
	_, err = tx.Exec("INSERT INTO transactions (from_account_id, to_account_id, amount) VALUES (?, ?, ?)", from, to, amount)
	if err != nil {
		tx.Rollback()
		log.Fatal("转账失败，记录交易失败，交易回滚", err)
		return
	}
	//提交事务
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatal("转账失败，提交事务失败", err)
		return
	}
	fmt.Println("转账成功")
}

func queryAccounts(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM accounts")
	if err != nil {
		log.Fatal("查询数据失败:", err)
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var balance int
		rows.Scan(&id, &balance)
		fmt.Printf("账户：%s,余额：%d \n", id, balance)
	}
}
