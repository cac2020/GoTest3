package exam1

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

/*
*
  - 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
  - 要求 ：
    编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
    编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
    编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
    编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
*/

/**
 * 使用Go内置sql库 数据库使用sqlite3
 */
func StudentTest() {
	//1.创建数据库连接
	db, err1 := sql.Open("sqlite3", ":memory:")
	if err1 != nil {
		log.Fatal("创建数据库连接失败:", err1)
		panic(err1)
	} else {
		fmt.Println("创建数据库连接成功")
	}
	defer db.Close()

	//2.建表
	schema := `
    CREATE TABLE IF NOT EXISTS students (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        age INTEGER NOT NULL,
        grade TEXT NOT NULL
    );`
	_, err2 := db.Exec(schema)
	if err2 != nil {
		log.Fatal("创建表失败:", err2)
		panic(err2)
	} else {
		fmt.Println("创建创建表成功")
	}

	//3.编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
	insertSql := `INSERT INTO students (name, age, grade) VALUES (?, ?, ?)`
	ret1, err3 := db.Exec(insertSql, "张三", 20, "三年级")
	if err3 != nil {
		log.Fatal("插入数据失败:", err3)
		panic(err3)
	} else {
		id, _ := ret1.LastInsertId()
		cnt, _ := ret1.RowsAffected()
		fmt.Printf("插入数据成功：ID= %d ，插入条数： %d \n", id, cnt)
	}

	//4.编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	querySql := `SELECT * FROM students WHERE age > ?`
	rows, err4 := db.Query(querySql, 18)
	if err4 != nil {
		log.Fatal("查询数据失败:", err4)
		panic(err4)
	} else {
		for rows.Next() {
			var id int
			var name string
			var age int
			var grade string
			rows.Scan(&id, &name, &age, &grade)
			fmt.Printf("ID: %d, Name: %s, Age: %d, Grade: %s\n", id, name, age, grade)
		}
	}
	defer rows.Close()

	//5.编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	updateSql := `UPDATE students SET grade = ? WHERE name = ?`
	ret2, err5 := db.Exec(updateSql, "四年级", "张三")
	if err5 != nil {
		log.Fatal("更新数据失败:", err5)
		panic(err5)
	} else {
		cnt, _ := ret2.RowsAffected()
		fmt.Printf("更新数据成功：更新条数： %d \n", cnt)
	}
	//6.编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	delSql := `DELETE FROM students WHERE age < ?`
	ret3, err6 := db.Exec(delSql, 15)
	if err6 != nil {
		log.Fatal("删除数据失败:", err6)
		panic(err6)
	} else {
		cnt, _ := ret3.RowsAffected()
		fmt.Printf("删除数据成功：删除条数： %d \n", cnt)
	}

}
