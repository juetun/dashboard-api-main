package pojos

import (
	"html/template"

	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/web/models"
)

type PostStore struct {
	Title    string `json:"title"`
	Category int    `json:"category"`
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
func (c *CateStore) Message() map[string]int {
	return map[string]int{
		"Name.Required":        402000002,
		"Name.MaxSize":         402000006,
		"DisplayName.Required": 402000003,
		"DisplayName.MaxSize":  402000007,
		"ParentId.Min":         402000004,
		"SeoDesc.Required":     402000005,
		"SeoDesc.MaxSize":      402000008,
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

type AuthLogin struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Captcha    string `json:"captcha"`
	CaptchaKey string `json:"captchaKey"`
}

type AuthRegister struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status int    `json:"status"`
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
	Paginate    Paginate
}

type Paginate struct {
	Limit   int `json:"limit"`
	Count   int `json:"count"`
	Total   int `json:"total"`
	Last    int `json:"last"`
	Current int `json:"current"`
	Next    int `json:"next"`
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
