package Models

import (
	"database/sql"
	"net/http"
	"time"
)

type Server struct {
	DB *sql.DB
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

type User struct {
	ID       int
	Email    string
	Password string
	Username string
}

type Post struct {
	ID          int
	Content     string
	ID_Category int
	ID_User     int
}

type PostAndStuff struct {
	ID              int
	Title           string
	Content         string
	Username        string
	User_Email      string
	Uuid_Categories string
	Time            time.Time
	Categories      string
	Comments        []PostComments
	Likes           int
	Dislikes        int
}
type LikeUser struct {
	UserId     int
	LikesCount int
	Username   string
}
type Category struct {
	ID      int
	Content string
}

type Comment struct {
	ID      int
	ID_Post int
	Content string
}
type PostComments struct {
	ID_Post    int
	ID_Comment int
	Content    string
	Username   string
	Likes      int
	Dislikes   int
}

// each session contains the username of the user and the time at which it expires
type Session struct {
	Token    string
	Username string
	Expiry   time.Time
}

type HtmlResponse struct {
	Session    Session
	Data       interface{}
	Categories []Category
	Users      []User
	TotalPages []int
	IdUser     int
	TopUser    []LikeUser
	PageNum    int
}

type ErrorPageResponse struct {
	Code int
	Text string
}
