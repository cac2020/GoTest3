package exam3

import (
	"GoTest3/gorm/exam1"
	"fmt"

	"gorm.io/gorm"
)

/**
* 题目3：钩子函数
	继续使用博客系统的模型。
	要求 ：
	为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
*/

func HookTest(db *gorm.DB) {
	testPostNum(db)
	testCommentStatus(db)
}

func testPostNum(db *gorm.DB) {
	var user exam1.User
	db.Take(&user)
	fmt.Println("初始用户信息文章数量：", user.PostNum)

	//添加文章
	post := exam1.Post{
		Title:    "文章1",
		Content:  "内容1",
		UserID:   user.ID,
		Catagory: "技术",
	}
	db.Create(&post)
	db.Find(&user, user.ID)
	fmt.Println("添加文章后用户信息文章数量：", user.PostNum)

	//删除文章
	db.Delete(&post)
	db.Find(&user, user.ID)
	fmt.Println("删除文章后用户信息文章数量：", user.PostNum)
}

func testCommentStatus(db *gorm.DB) {
	var post exam1.Post
	db.Preload("Comments").Order("comment_num desc").Limit(1).Find(&post)
	fmt.Printf("初始文章信息ID：%d，评论状态：%d，评论数量：%d \n", post.ID, post.CommentStatus, post.CommentNum)

	//删除当前文章下所有评论
	//批量删除 优先考虑性能  传入钩子函数的Comment 默认都是零值  所以此场景 先查询所有评论 然后遍历删除
	for _, comment := range post.Comments {
		db.Delete(&comment)
	}
	db.Find(&post, post.ID)
	fmt.Printf("删除所有评论后文章信息ID：%d，评论状态：%d，评论数量：%d \n", post.ID, post.CommentStatus, post.CommentNum)

	//添加评论
	comment := exam1.Comment{
		Content: "评论1",
		PostID:  post.ID,
	}
	db.Create(&comment)
	db.Find(&post, post.ID)
	fmt.Printf("添加评论后文章信息ID：%d，评论状态：%d，评论数量：%d \n", post.ID, post.CommentStatus, post.CommentNum)
}
