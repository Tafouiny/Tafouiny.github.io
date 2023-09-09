package Tools

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	Models "forum/models"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

const Port = ":8080"

const (
	sessionCookieName = "session_forum_token"
	sessionDuration   = 3600 * time.Second // 60 minutes
)

var Routes = []Models.Route{
	{Path: "/", Handler: IndexHandler, Method: "GET"},
	{Path: "/signin", Handler: SigninHandler, Method: "POST"},
	{Path: "/signup", Handler: SignupHandler, Method: "POST"},
	{Path: "/about", Handler: AboutHandler, Method: "GET"},
	{Path: "/logout", Handler: LogoutHandler},
	{Path: "/postdetails", Handler: PostDetailHandler},
	{Path: "/createpost", Handler: CreatePostHandler, Method: "POST"},
	{Path: "/createcomment", Handler: CreateCommentHandler, Method: "POST"},
	{Path: "/like", Handler: LikePostHandler, Method: "POST"},
	{Path: "/dislike", Handler: DislikePostHandler, Method: "POST"},
	{Path: "/likecomment", Handler: LikeCommentHandler, Method: "POST"},
	{Path: "/dislikecomment", Handler: DislikeCommentHandler, Method: "POST"},
	{Path: "/category", Handler: FilterCategoryHandler, Method: "POST"},
	{Path: "/myposts", Handler: FilterPostsByYourPostsHandler, Method: "POST"},
	{Path: "/likedposts", Handler: FilterPostsByLikedPostsHandler, Method: "POST"},
	{Path: "/postedby", Handler: FilterPostsByHimPostsHandler, Method: "POST"},
}

func InitDb() {
	db := OpenDB()
	defer db.Close()
	users := GetUsers()
	if len(users) == 0 {
		query, err := os.ReadFile("./database/backupsforum.db.sql")
		if err != nil {
			panic(err)
		}
		if _, err := db.Exec(string(query)); err != nil {
			panic(err)
		}
	} else {
		return
	}
}

func MiddlewareError(path string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				StatusInternalServerError(w)
				return
			}
		}()

		if r.URL.Path != path {
			ErrorPageResponse := Models.ErrorPageResponse{
				Code: 404,
				Text: "Oops 🫢! There is no post to comment here!",
			}
			w.WriteHeader(http.StatusNotFound)
			RenderTemplate(w, "error", ErrorPageResponse)
			return
		}

		handler(w, r)
	}
}

func StatusBadRequest(w http.ResponseWriter, text string) {
	ErrorPageResponse := Models.ErrorPageResponse{
		Code: 400,
		Text: text,
	}
	w.WriteHeader(http.StatusBadRequest)
	RenderTemplate(w, "error", ErrorPageResponse)
}

func StatusUnauthorized(w http.ResponseWriter) {
	ErrorPageResponse := Models.ErrorPageResponse{
		Code: 401,
		Text: "Oops 🫢! You have to sign in first!",
	}
	w.WriteHeader(http.StatusUnauthorized)
	RenderTemplate(w, "error", ErrorPageResponse)
}

func StatusInternalServerError(w http.ResponseWriter) {
	ErrorPageResponse := Models.ErrorPageResponse{
		Code: 500,
		Text: "Oops 🫢! Some posts took vacation!",
	}
	w.WriteHeader(http.StatusInternalServerError)
	RenderTemplate(w, "error", ErrorPageResponse)
}

// Load HTML templates
var templates = template.Must(template.ParseGlob("./views/templates/*.html"))

// RenderTemplate renders the specified template with the given data
func RenderTemplate(w http.ResponseWriter, tmpl string, i interface{}) {
	htmlFile := tmpl + ".html"
	err := templates.ExecuteTemplate(w, htmlFile, i)
	if err != nil {
		StatusInternalServerError(w)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Catch the cookie sent by login handler
	userSession, auth := IsAuthenticated(r)
	categories := GetCategories()

	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}

	idUsername := GetIDUserFromSession(userSession)

	pageNum, err := strconv.Atoi(pageParam)
	if err != nil {
		StatusBadRequest(w, "Error with the pages")
		return
	}

	data, totalPages, err := GetPostAndStuff(pageNum)

	if pageNum > totalPages[len(totalPages)-1] || pageNum <= 0 {
		ErrorPageResponse := Models.ErrorPageResponse{
			Code: 404,
			Text: "Oops 🫢! There is no post to comment here!",
		}
		w.WriteHeader(http.StatusNotFound)
		RenderTemplate(w, "error", ErrorPageResponse)
		return
	}

	if err != nil {
		StatusBadRequest(w, "Oops 🫢! some posts have been retained by the letter carrier!")
		return
	}
	var sessionModel Models.Session
	if auth {
		sessionModel = userSession
	} else {
		sessionModel = Models.Session{}
	}
	topUsers := GetUsersWithMostLikes()
	htmlResponse := Models.HtmlResponse{
		Session:    sessionModel,
		Data:       data,
		Categories: categories,
		TotalPages: totalPages,
		IdUser:     idUsername,
		TopUser:    topUsers,
		PageNum:    pageNum,
	}

	RenderTemplate(w, "index", htmlResponse)
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	_, auth := IsAuthenticated(r)
	if auth {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, ok := CheckCredentials(email, password)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "<script>alert('Woops! Email Or password is incorrect.')</script>")
			RenderTemplate(w, "signin", nil)
			return
		}

		sessionToken := uuid.Must(uuid.NewV4()).String()
		expiresAt := time.Now().Add(sessionDuration)

		// Check if an user already have a session. If it is the case it deletes the ancient session and creates a new.
		// We cannot have two differents session on our forum -> audit
		session, checkSession := UserHasAlreadyASession(user.Username)
		if checkSession {
			DeleteSession(session.Token)
			CreateSession(sessionToken, user.Username, expiresAt)
		} else {
			CreateSession(sessionToken, user.Username, expiresAt)
		}

		http.SetCookie(w, &http.Cookie{
			Name:    sessionCookieName,
			Value:   sessionToken,
			Expires: expiresAt,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		RenderTemplate(w, "signin", nil)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(sessionCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		StatusBadRequest(w, "Oops 🫢! some posts have been retained by the letter carrier!")
		return
	}

	sessionToken := c.Value

	DeleteSession(sessionToken)

	// We need to let the client know that the cookie is expired In the response, we set the session token to an empty
	http.SetCookie(w, &http.Cookie{
		Name:    sessionCookieName,
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	_, auth := IsAuthenticated(r)
	if auth {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")

		email := r.FormValue("email")
		password := r.FormValue("password")
		confPassword := r.FormValue("cpassword")

		if password != confPassword {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "<script>alert('Woops! Password is different.')</script>")
			RenderTemplate(w, "signup", nil)
			return
		}

		if isBlank(username) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "<script>alert('Woops! Invalid username.')</script>")
			RenderTemplate(w, "signup", nil)
			return
		}
		if isBlank(email) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "<script>alert('Woops! Invalid email.')</script>")
			RenderTemplate(w, "signup", nil)
			return
		}

		_, ok := IfUserExist(username, email)
		if ok {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "<script>alert('Woops! Email Or username already exists.')</script>")
			RenderTemplate(w, "signup", nil)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			HandleError(err, "hasign password using bcrypt")
			StatusBadRequest(w, "Oops 🫢! some posts have been retained by the letter carrier!")
			return
		}

		CreateUser(email, string(hashedPassword), username)
		fmt.Fprintf(w, "<script>alert('Good! Registration succesffully made.')</script>")
	}
	RenderTemplate(w, "signup", nil)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	// Catch the cookie sent by login handler
	userSession, auth := IsAuthenticated(r)
	if auth {
		HtmlResponse := Models.HtmlResponse{
			Session:    userSession,
			Data:       nil,
			Categories: GetCategories(),
		}
		RenderTemplate(w, "about", HtmlResponse)
		return
	}

	// If there's no valid session or no cookie, display nav bar disconected
	HtmlResponse := Models.HtmlResponse{
		Session:    Models.Session{},
		Data:       nil,
		Categories: GetCategories(),
	}
	RenderTemplate(w, "about", HtmlResponse)
}

func CreatePostCategory(uuid string, id_category int) {
	db := OpenDB()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO post_categories(uuid_post_categories, id_category) values(?,?)")
	if err != nil {
		HandleError(err, "preparing insertion of post_categories")
		return
	}
	_, err = stmt.Exec(uuid, id_category)
	if err != nil {
		HandleError(err, "Excecuting insertion of post_categories")
		return
	}
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userSession, auth := IsAuthenticated(r)

	if !auth {
		StatusUnauthorized(w)
		return
	}
	categories := GetCategories()
	if r.Method == http.MethodPost {
		content := r.FormValue("content")
		//content = strings.ReplaceAll(content, "\n", "<br>")
		title := r.FormValue("title")
		categories_form := r.Form["category"]

		if OnlySpace(title) || OnlySpace(content) || len(categories_form) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "<script>alert('Woops! Post title,content or category may be empty !')</script>")

			HtmlResponse := Models.HtmlResponse{
				Session: userSession,
				Data:    categories,
			}
			RenderTemplate(w, "createpost", HtmlResponse)
			return
		}

		uuid_categories := uuid.Must(uuid.NewV4()).String()

		for _, ID_Category := range categories_form {
			ID, err := strconv.Atoi(ID_Category)
			if err != nil {
				StatusBadRequest(w, "Oops 🫢! some posts have been retained by the letter carrier!")
				return
			}
			CreatePostCategory(uuid_categories, ID)
		}
		user_id := GetIDUserFromSession(userSession)
		createdAt := time.Now()
		CreatePost(title, content, uuid_categories, user_id, createdAt)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	HtmlResponse := Models.HtmlResponse{
		Session: userSession,
		Data:    categories,
	}

	RenderTemplate(w, "createpost", HtmlResponse)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userSession, auth := IsAuthenticated(r)

	if !auth {
		StatusUnauthorized(w)
		return
	}

	if r.Method == http.MethodPost {

		content := r.FormValue("comment")
		post_ID := r.URL.Query().Get("postid")
		postID, err := strconv.Atoi(post_ID)
		if err != nil {
			StatusBadRequest(w, "Oops 🫢! some posts have been retained by the letter carrier!")
			return
		}

		url := "/postdetails?postid=" + post_ID

		if OnlySpace(content) {
			http.Redirect(w, r, url, http.StatusSeeOther)
			return
		}
		user_id := GetIDUserFromSession(userSession)
		CreateComment(postID, user_id, content)

		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	}
}

func FilterCategoryHandler(w http.ResponseWriter, r *http.Request) {
	userSession, _ := IsAuthenticated(r)
	var posts []Models.PostAndStuff

	categoryIDStr := r.URL.Query().Get("categoryID")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		StatusBadRequest(w, "Oops 🫢! some posts have been retained by the letter carrier!")
		return
	}

	uuids, err := GetUUIDsByCategoryID(categoryID)
	if err != nil {
		StatusInternalServerError(w)
		return
	}

	posts, pages, err := GetPostsByUUIDs(uuids)
	if err != nil {
		StatusInternalServerError(w)
		return
	}

	HtmlResponse := Models.HtmlResponse{
		Session:    userSession,
		Data:       posts,
		TotalPages: pages,
		Categories: GetCategories(),
		TopUser:    GetUsersWithMostLikes(),
		IdUser:     GetIDUserFromSession(userSession),
		PageNum:    1,
	}
	RenderTemplate(w, "index", HtmlResponse)
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	userSession, authSession := IsAuthenticated(r)
	if !authSession {
		StatusUnauthorized(w)
		return
	}

	post_id := r.URL.Query().Get("postid")
	postID, err := strconv.Atoi(post_id)
	if err != nil {
		StatusBadRequest(w, "Oops 🫢! some posts have been retained by the letter carrier!")
		return
	}
	userID := GetIDUserFromSession(userSession)

	if UserHasAlreadyLiked(postID, userID) {
		RemoveLike(postID, userID)
	} else if UserHasAlreadyDisliked(postID, userID) {
		RemoveDislike(postID, userID)
		LikePost(postID, userID)
	} else {
		LikePost(postID, userID)
	}

	url := "/postdetails?postid=" + post_id
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	userSession, authSession := IsAuthenticated(r)
	if !authSession {
		StatusUnauthorized(w)
		return
	}

	post_id := r.URL.Query().Get("postid")
	postID, err := strconv.Atoi(post_id)
	if err != nil {
		StatusBadRequest(w, "Oops 🫢! some posts have been retained by the letter carrier!")
		return
	}
	userID := GetIDUserFromSession(userSession)

	if UserHasAlreadyLiked(postID, userID) {
		RemoveLike(postID, userID)
		DislikePost(postID, userID)
	} else if UserHasAlreadyDisliked(postID, userID) {
		RemoveDislike(postID, userID)
	} else {
		DislikePost(postID, userID)
	}

	url := "/postdetails?postid=" + post_id
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func FilterPostsByYourPostsHandler(w http.ResponseWriter, r *http.Request) {
	userSession, auth := IsAuthenticated(r)

	if !auth {
		StatusUnauthorized(w)
		return
	}

	idUsername := GetIDUserFromSession(userSession)

	var posts []Models.PostAndStuff
	posts, pages, err := GetPostsByUserID(idUsername)
	if err != nil {
		StatusInternalServerError(w)
		return
	}

	var sessionModel Models.Session
	if auth {
		sessionModel = userSession
	} else {
		sessionModel = Models.Session{}
	}

	HtmlResponse := Models.HtmlResponse{
		Session:    sessionModel,
		Data:       posts,
		TotalPages: pages,
		Categories: GetCategories(),
		IdUser:     idUsername,
		TopUser:    GetUsersWithMostLikes(),
		PageNum:    1,
	}
	RenderTemplate(w, "index", HtmlResponse)
}

func FilterPostsByHimPostsHandler(w http.ResponseWriter, r *http.Request) {
	userSession, auth := IsAuthenticated(r)

	postedby := r.URL.Query().Get("user")

	idUsername := GetUserByUsername(postedby)

	var posts []Models.PostAndStuff
	posts, pages, err := GetPostsByUserID(idUsername)
	if err != nil {
		StatusInternalServerError(w)
		return
	}

	var sessionModel Models.Session
	if auth {
		sessionModel = userSession
	} else {
		sessionModel = Models.Session{}
	}

	HtmlResponse := Models.HtmlResponse{
		Session:    sessionModel,
		Data:       posts,
		TotalPages: pages,
		Categories: GetCategories(),
		IdUser:     idUsername,
		TopUser:    GetUsersWithMostLikes(),
	}
	RenderTemplate(w, "index", HtmlResponse)
}

func FilterPostsByLikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	userSession, auth := IsAuthenticated(r)

	if !auth {
		StatusUnauthorized(w)
		return
	}

	idUsername := GetIDUserFromSession(userSession)

	var posts []Models.PostAndStuff
	posts, pages, err := GetPostsLikedByUserID(idUsername)
	if err != nil {
		StatusInternalServerError(w)
		return
	}

	var sessionModel Models.Session
	if auth {
		sessionModel = userSession
	} else {
		sessionModel = Models.Session{}
	}
	topUsers := GetUsersWithMostLikes()
	HtmlResponse := Models.HtmlResponse{
		Session:    sessionModel,
		Data:       posts,
		TotalPages: pages,
		Categories: GetCategories(),
		IdUser:     idUsername,
		TopUser:    topUsers,
		PageNum:    1,
	}
	RenderTemplate(w, "index", HtmlResponse)
}

func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
	userSession, _ := IsAuthenticated(r)

	post_id := r.URL.Query().Get("postid")
	postID, err := strconv.Atoi(post_id)
	if err != nil {
		StatusBadRequest(w, "Oops 🫢! Invalid postId.")
		return
	}

	posts, pages, err := GetPostsByID(1, postID)
	if err != nil {
		StatusBadRequest(w, "Oops 🫢! Invalid postId.")
		return
	}

	topUsers := GetUsersWithMostLikes()
	HtmlResponse := Models.HtmlResponse{
		Session:    userSession,
		Data:       posts,
		TotalPages: pages,
		Categories: GetCategories(),
		TopUser:    topUsers,
		PageNum:    1,
	}
	RenderTemplate(w, "post-details", HtmlResponse)
}

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	userSession, auth := IsAuthenticated(r)
	if !auth {
		StatusUnauthorized(w)
		return
	}
	comment_id := r.URL.Query().Get("comment_id")
	post_id := r.URL.Query().Get("postid")
	commentID, err := strconv.Atoi(comment_id)
	if err != nil {
		http.Error(w, "/", http.StatusBadRequest)
	}
	userID := GetIDUserFromSession(userSession)

	if UserHasAlreadyLikedComment(commentID, userID) {
		RemoveLikeFromComment(commentID, userID)
	} else if UserHasAlreadyDislikedComment(commentID, userID) {
		RemoveDislikeFromComment(commentID, userID)
		LikeComment(commentID, userID)
	} else {
		LikeComment(commentID, userID)
	}
	url := "/postdetails?postid=" + post_id
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	userSession, auth := IsAuthenticated(r)
	if !auth {
		StatusUnauthorized(w)
		return
	}
	comment_id := r.URL.Query().Get("comment_id")
	post_id := r.URL.Query().Get("postid")
	commentID, err := strconv.Atoi(comment_id)
	if err != nil {
		http.Error(w, "/", http.StatusBadRequest)
	}
	userID := GetIDUserFromSession(userSession)
	if UserHasAlreadyLikedComment(commentID, userID) {
		RemoveLikeFromComment(commentID, userID)
		DislikeComment(commentID, userID)
	} else if UserHasAlreadyDislikedComment(commentID, userID) {
		RemoveDislikeFromComment(commentID, userID)
	} else {
		DislikeComment(commentID, userID)
	}
	url := "/postdetails?postid=" + post_id
	http.Redirect(w, r, url, http.StatusSeeOther)
}
