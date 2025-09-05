package exam2

import (
	"GoTest3/gorm/exam1"
	"fmt"

	"gorm.io/gorm"
)

/**
* 题目2：关联查询
	基于上述博客系统的模型定义。
	要求 ：
	编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	编写Go代码，使用Gorm查询评论数量最多的文章信息。
*/

func QryTest(db *gorm.DB) {
	//查询某个用户发布的所有文章及其对应的评论信息。
	var user exam1.User
	db.Take(&user)
	qryPostsByUserID(db, user.ID)

	//查询评论数量最多的文章信息。
	qryPostsWithMostComments(db)
}
func qryPostsByUserID(db *gorm.DB, userID uint) {
	var posts []exam1.Post
	db.Where("user_id = ?", userID).Preload("Comments").Find(&posts)
	fmt.Printf("用户%d发布的文章以及评论：%v \n", userID, posts)
}

func qryPostsWithMostComments(db *gorm.DB) {
	//查询评论数最多的文章ID
	var result struct {
		PostID uint
		Count  int
	}
	db.Model(&exam1.Comment{}).Select("post_id, count(1)").Group("post_id").Order("count(1) desc").Limit(1).Take(&result)

	//查询文章
	var post exam1.Post
	db.Where("id = ?", result.PostID).Preload("Comments").Find(&post)
	fmt.Printf("评论数最多的文章信息：%v \n", post)
}
