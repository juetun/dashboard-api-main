package wrappers

import (
	"html/template"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/models"
)

type PostStore struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Tags     []int  `json:"tags"`
	Summary  string `json:"summary"`
	Content  string `json:"content"`
}

type CateStore struct {
	Name        string `json:"name" valid:"Required;MaxSize(100)"`
	DisplayName string `json:"display_name" valid:"Required;MaxSize(100)"`
	ParentId    int    `json:"parent_id" valid:"Min(0)"`
	SeoDesc     string `json:"seo_desc" valid:"Required;MaxSize(250)"`
}

// 传参校验用
func (c *CateStore) Message() map[string]common.ValidationMessage {
	return map[string]common.ValidationMessage{
		"Name.Required.":        {Code: 402000002, Message: "请输入分类名称"},
		"Name.MaxSize.":         {Code: 402000006, Message: "分类名称不超过100个字符"},
		"DisplayName.Required.": {Code: 402000003, Message: "请输入分类别名"},
		"DisplayName.MaxSize.":  {Code: 402000007, Message: "分类别名不超过100个字符"},
		"ParentId.Min.":         {Code: 402000004, Message: "上级类型不能小于0"},
		"SeoDesc.Required.":     {Code: 402000005, Message: "请填写SEO描述"},
		"SeoDesc.MaxSize.":      {Code: 402000008, Message: "SEO描述不超过250个字符"},
	}
}

type TagStore struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	SeoDesc     string `json:"seo_desc"`
}

type LinkStore struct {
	Name  string `json:"name"`
	Link  string `json:"link"`
	Order int    `json:"order"`
}



type ConsolePostList struct {
	Post     ConsolePost  `json:"post"`
	Tags     []ConsoleTag `json:"tags"`
	Category ConsoleCate  `json:"category"`
	View     ConsoleView  `json:"view"`
	Author   ConsoleUser  `json:"author"`
}

type ConsolePost struct {
	Id        int             `json:"id"`
	Uid       string          `json:"uid"`
	UserId    int             `json:"user_id"`
	Title     string          `json:"title"`
	Summary   string          `json:"summary"`
	Original  string          `json:"original"`
	Content   string          `json:"content"`
	Password  string          `json:"password"`
	DeletedAt base.TimeNormal `json:"deleted_at"`
	CreatedAt base.TimeNormal `json:"created_at"`
	UpdatedAt base.TimeNormal `json:"updated_at"`
}

type ConsoleTag struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	SeoDesc     string `json:"seo_desc"`
	Num         int    `json:"num"`
}

type ConsoleCate struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	SeoDesc     string `json:"seo_desc"`
}

type ConsoleUser struct {
	UserHid string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Status  int    `json:"status"`
}

type ConsoleSystem struct {
	Title        string `json:"title"`
	Keywords     string `json:"keywords"`
	Theme        int    `json:"theme"`
	Description  string `json:"description"`
	RecordNumber string `json:"record_number"`
}

type ConsoleView struct {
	Num int `json:"num"`
}

type IndexPostList struct {
	PostListArr *[]ConsolePostList
	Paginate    utils.Paginate
}

type IndexPost struct {
	Id        int             `json:"id"`
	Uid       string          `json:"uid"`
	UserId    int             `json:"user_id"`
	Title     string          `json:"title"`
	Summary   string          `json:"summary"`
	Original  string          `json:"original"`
	Content   template.HTML   `json:"content"`
	Password  string          `json:"password"`
	DeletedAt base.TimeNormal `json:"deleted_at"`
	CreatedAt base.TimeNormal `json:"created_at"`
	UpdatedAt base.TimeNormal `json:"updated_at"`
}

type IndexPostDetail struct {
	Post     IndexPost      `json:"post"`
	Tags     []ConsoleTag   `json:"tags"`
	Category ConsoleCate    `json:"category"`
	View     ConsoleView    `json:"view"`
	Author   ConsoleUser    `json:"author"`
	LastPost *models.ZPosts `json:"last_post"`
	NextPost *models.ZPosts `json:"next_post"`
}

type IndexGithubParam struct {
	GithubName         string
	GithubRepo         string
	GithubClientId     string
	GithubClientSecret string
	GithubLabels       string
}

type Category struct {
	Cates models.ZCategories `json:"cates"`
	Html  string             `json:"html"`
}

type IndexCategory struct {
	Cates models.ZCategories `json:"cates"`
	Html  template.HTML      `json:"html"`
}

//
type TagCount struct {
	PostId int `gorm:"column:post_id;" json:"post_id"`
	TagId  int `gorm:"column:tag_id;" json:"tag_id"`
	Count  int `gorm:"column:count;" json:"count"`
}
