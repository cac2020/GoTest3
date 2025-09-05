package exam1

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/gorm"
)

/**
* 题目1：模型定义
* 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
* 要求 ：
* 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
* 编写Go代码，使用Gorm创建这些模型对应的数据库表。
 */
type User struct {
	gorm.Model
	Name     string `gorm:"size:50;not null;"` //`gorm:"size:50;not null;unique"`
	Age      sql.NullInt16
	Password sql.NullString
	Email    sql.NullString
	Active   sql.NullBool `gorm:"default:true"`
	PostNum  int          `gorm:"default:0"` //文章数量
	Posts    []Post       `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
}

type Post struct {
	gorm.Model
	Title         string `gorm:"size:50;not null;"`
	Content       string `gorm:"not null;"`
	Catagory      string
	Active        sql.NullBool `gorm:"default:true"`
	CommentNum    int          `gorm:"default:0"` //评论数量
	CommentStatus int          `gorm:"default:0"` //状态 0-无评论 1：已评论
	UserID        uint
	Comments      []Comment `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE;"`
}

func (u *Post) AfterSave(tx *gorm.DB) (err error) {
	//新增文章	后更新用户文章数量
	var user User
	tx.Find(&user, "id = ?", u.UserID)
	user.PostNum += 1
	//仅更新指定字段
	tx.Model(&User{}).Where("id = ?", u.UserID).Update("post_num", user.PostNum)
	return
}

func (u *Post) AfterDelete(tx *gorm.DB) (err error) {
	//删除文章后 更新用户文章数量
	var user User
	tx.Find(&user, "id = ?", u.UserID)
	user.PostNum -= 1
	//仅更新指定字段
	tx.Model(&User{}).Where("id = ?", u.UserID).Update("post_num", user.PostNum)
	return
}

type Comment struct {
	gorm.Model
	Content string       `gorm:"not null;"`
	Active  sql.NullBool `gorm:"default:true"`
	PostID  uint
}

func (c *Comment) AfterSave(tx *gorm.DB) (err error) {
	//新增评论后	更新新文章评论数量、评论状态
	var post Post
	tx.Find(&post, "id = ?", c.PostID)
	post.CommentNum += 1
	if post.CommentNum > 0 {
		post.CommentStatus = 1
	}
	//仅更新指定字段
	tx.Model(&Post{}).Where("id = ?", c.PostID).Updates(Post{CommentNum: post.CommentNum, CommentStatus: post.CommentStatus})
	return
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	//删除评论后 更新新文章评论数量、评论状态
	var commentCount int64
	commentStatus := 1
	if err := tx.Model(&Comment{}).
		Where("post_id = ? AND deleted_at IS NULL", c.PostID).
		Count(&commentCount).Error; err != nil {
		return err
	}
	if commentCount <= 0 {
		commentStatus = 0
	}
	fmt.Printf("更新前参数：postID：%d，评论数量：%d，评论状态：%d \n", c.PostID, commentCount, commentStatus)
	//仅更新指定字段
	//tx.Model(&Post{}).Where("id = ?", c.PostID).Updates(Post{CommentNum: int(commentCount), CommentStatus: commentStatus})
	tx.Model(&Post{}).Where("id = ?", c.PostID).Updates(map[string]interface{}{"comment_num": int(commentCount), "comment_status": commentStatus})

	//更新后强制查询一次查看结果
	//var post Post
	//tx.Model(&Post{}).Where("id = ?", c.PostID).Select("comment_num,comment_status").First(&post)
	//fmt.Printf("删除评论后文章信息ID：%d，评论状态：%d，评论数量：%d \n", c.PostID, post.CommentStatus, post.CommentNum)
	return
}

func CreateTable(db *gorm.DB) {
	//创建表
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatal("创建表失败:", err)
		panic(err)
	} else {
		fmt.Println("创建表成功")
	}

	//插入数据
	//1.创建用户
	users := []User{
		{Name: "张三", Age: sql.NullInt16{Int16: 18, Valid: true}, Email: sql.NullString{String: "123@qq.com", Valid: true},
			Posts: []Post{
				{
					Title:    "华山论剑200",
					Content:  "葵花宝典就是垃圾",
					Catagory: "武林大会",
					Comments: []Comment{
						{
							Content: "岳不群就是个变态",
						},
					},
				},
			}},
		{Name: "李四", Age: sql.NullInt16{Int16: 18, Valid: true}, Email: sql.NullString{String: "456@qq.com", Valid: true},
			Posts: []Post{
				{
					Title:    "华山论剑200",
					Content:  "葵花宝典就是垃圾",
					Catagory: "武林大会",
					Comments: []Comment{
						{
							Content: "岳不群就是个变态",
						},
					},
				},
			}},
		{Name: "王五", Age: sql.NullInt16{Int16: 18, Valid: true}, Email: sql.NullString{String: "789@qq.com", Valid: true},
			Posts: []Post{
				{
					Title:    "华山论剑1",
					Content:  "葵花宝典就是垃圾",
					Catagory: "武林大会",
					Comments: []Comment{
						{
							Content: "岳不群就是个变态",
						},
					},
				},
				{
					Title:    "华山论剑2",
					Content:  "葵花宝典就是垃圾",
					Catagory: "武林大会",
					Comments: []Comment{
						{
							Content: "岳不群就是个变态",
						},
						{
							Content: "岳不群就是个变态",
						},
					},
				},
				{
					Title:    "华山论剑3",
					Content:  "葵花宝典就是垃圾",
					Catagory: "武林大会",
					Comments: []Comment{
						{
							Content: "岳不群就是个变态",
						},
						{
							Content: "岳不群就是个变态",
						},
						{
							Content: "岳不群就是个变态",
						},
					},
				},
			}},
		{Name: "赵六", Age: sql.NullInt16{Int16: 18, Valid: true}, Email: sql.NullString{String: "987@qq.com", Valid: true},
			Posts: []Post{
				{
					Title:    "华山论剑",
					Content:  "葵花宝典就是垃圾",
					Catagory: "武林大会",
					Comments: []Comment{
						{
							Content: "岳不群就是个变态",
						},
						{
							Content: "令狐冲为什么要去华山论剑",
						},
						{
							Content: "令狐冲的华山论剑讲的是啥",
						},
						{
							Content: "各大门派的看门绝技是什么",
						},
					},
				},
				{
					Title:    "攻打光明顶",
					Content:  "打败邪教，夺取屠龙刀",
					Catagory: "武林大会",
					Comments: []Comment{
						{
							Content: "张无忌使用的那招叫什么",
						},
						{
							Content: "寒冰真气厉害吗",
						},
						{
							Content: "周芷若喜欢张无忌吗",
						},
						{
							Content: "赵敏有没有带领蒙古人前去参加",
						},
						{
							Content: "灭绝师太更年期了",
						},
						{
							Content: "灭绝师太那把倚天剑厉害还是屠龙刀厉害",
						},
					},
				},
			},
		},
	}
	db.Create(&users)
}
