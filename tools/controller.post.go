package Tools

import (
	"math"
	"sort"
	"time"

	Models "forum/models"
)

func CreatePost(Title, Content, uuid string, id_user int, time time.Time) {
	db := OpenDB()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO posts(post_title, post_content, uuid_post_categories, id_user, time) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		HandleError(err, "preparing insertion of post")
		return
	}
	res, err := stmt.Exec(Title, Content, uuid, id_user, time)
	if err != nil {
		HandleError(err, "executing insertion of post")
		return
	}
	_, _ = res.RowsAffected()
}

func GetPostAndStuff(pageNum int) (postandstuff []Models.PostAndStuff, pages []int, err error) {
	db := OpenDB()
	defer db.Close()

	var totalPosts int
	err = db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&totalPosts)
	if err != nil {
		return nil, pages, err
	}

	if totalPosts <= 0 {
		return nil, nil, err
	}

	perPage := 5
	totalPages := int(math.Ceil(float64(totalPosts) / float64(perPage)))

	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	// Calculer l'offset pour la requête SQL-
	offset := (pageNum - 1) * perPage

	rows, err := db.Query(`
    SELECT p.id_post, p.post_title, p.post_content, u.username, u.email, p.uuid_post_categories, p.time
    FROM posts p
    INNER JOIN users u ON p.id_user = u.id_user
    ORDER BY p.time DESC
    LIMIT ? OFFSET ?
    `, perPage, offset)
	if err != nil {
		return nil, pages, err
	}
	defer rows.Close()

	rowsComment, err := db.Query(`SELECT id_post, id_comment, comment_content, username FROM comments NATURAL JOIN users`)
	if err != nil {
		return nil, pages, err
	}
	defer rowsComment.Close()

	var comments []Models.PostComments
	for rowsComment.Next() {
		var comment Models.PostComments
		err := rowsComment.Scan(&comment.ID_Post, &comment.ID_Comment, &comment.Content, &comment.Username)
		if err != nil {
			return nil, pages, err
		}
		comment.Likes = GetLikesFromComment(comment.ID_Comment)
		comment.Dislikes = GetDislikesFromComment(comment.ID_Comment)
		comments = append(comments, comment)
	}

	var posts []Models.PostAndStuff
	for rows.Next() {
		var post Models.PostAndStuff
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Username, &post.User_Email, &post.Uuid_Categories, &post.Time)
		if err != nil {
			return nil, pages, err
		}

		post.Comments = []Models.PostComments{} // Initialize comments for this post
		categories := GetCategoriesPost(post.Uuid_Categories)
		post.Categories = categories
		post.Likes = GetLikeFromPost(post.ID)
		post.Dislikes = GetDislikesFromPost(post.ID)
		for _, comment := range comments {
			if comment.ID_Post == post.ID {
				post.Comments = append(post.Comments, comment)
			}
		}
		posts = append(posts, post)
	}

	return posts, pages, nil
}

func GetCategoriesPost(uuid string) string {
	db := OpenDB()
	defer db.Close()

	categories := ""

	categoryRows, err := db.Query("SELECT id_category FROM post_categories WHERE uuid_post_categories = ?", uuid)
	if err != nil {
		HandleError(err, "Fetching categories from database.")
		return categories
	}
	defer categoryRows.Close()

	for categoryRows.Next() {
		var category int
		if err := categoryRows.Scan(&category); err != nil {
			HandleError(err, "Scanning categories from database.")
			return categories
		}

		var formattedCategory string
		err := db.QueryRow("SELECT content FROM categories WHERE id_category = ?", category).Scan(&formattedCategory)
		if err != nil {
			HandleError(err, "Fetching formatted category content from database.")
			continue
		}

		if categories == "" {
			categories += "#" + formattedCategory
		} else {
			categories += " #" + formattedCategory
		}
	}
	return categories
}

func GetPostsByUUIDs(uuids []string) ([]Models.PostAndStuff, []int, error) {
	db := OpenDB()
	defer db.Close()

	// Fetch all comments for all posts
	commentQuery := `SELECT id_post, id_comment, comment_content, username FROM comments NATURAL JOIN users`
	rowsComment, err := db.Query(commentQuery)
	if err != nil {
		return nil, nil, err
	}
	defer rowsComment.Close()

	commentsByPost := make(map[int][]Models.PostComments)
	for rowsComment.Next() {
		var comment Models.PostComments
		err := rowsComment.Scan(&comment.ID_Post, &comment.ID_Comment, &comment.Content, &comment.Username)
		if err != nil {
			return nil, nil, err
		}
		comment.Likes = GetLikesFromComment(comment.ID_Comment)
		comment.Dislikes = GetDislikesFromComment(comment.ID_Comment)
		commentsByPost[comment.ID_Post] = append(commentsByPost[comment.ID_Post], comment)
	}

	var posts []Models.PostAndStuff
	for _, uuid := range uuids {
		postQuery := `
            SELECT p.id_post, p.post_title, p.post_content, u.username, u.email, p.uuid_post_categories, p.time
            FROM posts p
            INNER JOIN users u ON p.id_user = u.id_user
            WHERE p.uuid_post_categories = ?
            ORDER BY p.time DESC
        `
		rows, err := db.Query(postQuery, uuid)
		if err != nil {
			return nil, nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var post Models.PostAndStuff
			err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Username, &post.User_Email, &post.Uuid_Categories, &post.Time)
			if err != nil {
				return nil, nil, err
			}

			post.Comments = commentsByPost[post.ID]
			categories := GetCategoriesPost(post.Uuid_Categories)
			post.Categories = categories

			posts = append(posts, post)
		}
	}

	totalPages := 1
	pages := make([]int, totalPages)
	for i := range pages {
		pages[i] = i + 1
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})
	return posts, pages, nil
}

func GetUUIDsByCategoryID(categoryID int) ([]string, error) {
	db := OpenDB()
	defer db.Close()

	var uuids []string

	rows, err := db.Query("SELECT uuid_post_categories FROM post_categories WHERE id_category = ?", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			return nil, err
		}
		uuids = append(uuids, uuid)
	}

	return uuids, nil
}

func GetPostsByUserID(userID int) (postandstuff []Models.PostAndStuff, pages []int, err error) {
	db := OpenDB()
	defer db.Close()

	var totalPosts int
	err = db.QueryRow("SELECT COUNT(*) FROM posts WHERE id_user = ?", userID).Scan(&totalPosts)
	if err != nil {
		return nil, pages, err
	}

	if totalPosts <= 0 {
		return nil, nil, err
	}

	totalPages := 1

	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	rows, err := db.Query(`
        SELECT p.id_post, p.post_title, p.post_content, u.username, u.email, p.uuid_post_categories, p.time
        FROM posts p
        INNER JOIN users u ON p.id_user = u.id_user
        WHERE p.id_user = ?
        ORDER BY p.time DESC
    `, userID)
	if err != nil {
		return nil, pages, err
	}
	defer rows.Close()

	rowsComment, err := db.Query(`SELECT id_post, id_comment, comment_content, username FROM comments NATURAL JOIN users`)
	if err != nil {
		return nil, pages, err
	}
	defer rowsComment.Close()

	var comments []Models.PostComments
	for rowsComment.Next() {
		var comment Models.PostComments
		err := rowsComment.Scan(&comment.ID_Post, &comment.ID_Comment, &comment.Content, &comment.Username)
		if err != nil {
			return nil, pages, err
		}
		comment.Likes = GetLikesFromComment(comment.ID_Comment)
		comment.Dislikes = GetDislikesFromComment(comment.ID_Comment)
		comments = append(comments, comment)
	}

	var posts []Models.PostAndStuff
	for rows.Next() {
		var post Models.PostAndStuff
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Username, &post.User_Email, &post.Uuid_Categories, &post.Time)
		if err != nil {
			return nil, pages, err
		}

		post.Comments = []Models.PostComments{} // Initialize comments for this post
		categories := GetCategoriesPost(post.Uuid_Categories)
		post.Categories = categories
		post.Likes = GetLikeFromPost(post.ID)
		post.Dislikes = GetDislikesFromPost(post.ID)
		for _, comment := range comments {
			if comment.ID_Post == post.ID {
				post.Comments = append(post.Comments, comment)
			}
		}
		posts = append(posts, post)
	}

	return posts, pages, nil
}

func GetPostsLikedByUserID(userID int) (postandstuff []Models.PostAndStuff, pages []int, err error) {
	db := OpenDB()
	defer db.Close()

	var totalPosts int
	err = db.QueryRow("SELECT COUNT(*) FROM post_likes WHERE id_user = ? AND like_type = 1", userID).Scan(&totalPosts)
	if err != nil {
		return nil, pages, err
	}

	if totalPosts <= 0 {
		return nil, nil, err
	}

	rows, err := db.Query(`
		SELECT p.id_post, p.post_title, p.post_content, u.username, u.email, p.uuid_post_categories, p.time
		FROM posts p
		INNER JOIN post_likes pl ON p.id_post = pl.id_post
		INNER JOIN users u ON p.id_user = u.id_user
		WHERE pl.id_user = ? AND pl.like_type = 1 
        ORDER BY p.time DESC
    `, userID)
	if err != nil {
		return nil, pages, err
	}
	defer rows.Close()

	rowsComment, err := db.Query(`SELECT id_post, id_comment, comment_content, username FROM comments NATURAL JOIN users`)
	if err != nil {
		return nil, pages, err
	}
	defer rowsComment.Close()

	var comments []Models.PostComments
	for rowsComment.Next() {
		var comment Models.PostComments
		err := rowsComment.Scan(&comment.ID_Post, &comment.ID_Comment, &comment.Content, &comment.Username)
		if err != nil {
			return nil, pages, err
		}
		comment.Likes = GetLikesFromComment(comment.ID_Comment)
		comment.Dislikes = GetDislikesFromComment(comment.ID_Comment)
		comments = append(comments, comment)
	}

	var posts []Models.PostAndStuff
	for rows.Next() {
		var post Models.PostAndStuff
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Username, &post.User_Email, &post.Uuid_Categories, &post.Time)
		if err != nil {
			return nil, pages, err
		}

		post.Comments = []Models.PostComments{} // Initialize comments for this post
		categories := GetCategoriesPost(post.Uuid_Categories)
		post.Categories = categories
		post.Likes = GetLikeFromPost(post.ID)
		post.Dislikes = GetDislikesFromPost(post.ID)
		for _, comment := range comments {
			if comment.ID_Post == post.ID {
				post.Comments = append(post.Comments, comment)
			}
		}
		posts = append(posts, post)
	}

	return posts, []int{1}, nil
}

func GetPostsByID(pageNum int, postID int) (Models.PostAndStuff, []int, error) {
	db := OpenDB()
	defer db.Close()

	var post Models.PostAndStuff
	perPage := 5
	offset := (pageNum - 1) * perPage

	// Fetch all comments for all posts
	commentQuery := `SELECT id_post, id_comment, comment_content, username FROM comments NATURAL JOIN users`
	rowsComment, err := db.Query(commentQuery)
	if err != nil {
		return post, nil, err
	}
	defer rowsComment.Close()

	commentsByPost := make(map[int][]Models.PostComments)
	for rowsComment.Next() {
		var comment Models.PostComments
		err := rowsComment.Scan(&comment.ID_Post, &comment.ID_Comment, &comment.Content, &comment.Username)
		if err != nil {
			return post, nil, err
		}
		comment.Likes = GetLikesFromComment(comment.ID_Comment)
		comment.Dislikes = GetDislikesFromComment(comment.ID_Comment)
		commentsByPost[comment.ID_Post] = append(commentsByPost[comment.ID_Post], comment)
	}

	postQuery := `
        SELECT p.id_post, p.post_title, p.post_content, u.username, u.email, p.uuid_post_categories, p.time
        FROM posts p
        INNER JOIN users u ON p.id_user = u.id_user
        WHERE p.id_post = ?
        ORDER BY p.time DESC
		LIMIT ? OFFSET ?
    `
	row := db.QueryRow(postQuery, postID, perPage, offset)
	if err != nil {
		return post, nil, err
	}

	err = row.Scan(&post.ID, &post.Title, &post.Content, &post.Username, &post.User_Email, &post.Uuid_Categories, &post.Time)
	if err != nil {
		return post, nil, err
	}

	post.Comments = commentsByPost[post.ID]
	post.Likes = GetLikeFromPost(postID)
	post.Dislikes = GetDislikesFromPost(postID)
	categories := GetCategoriesPost(post.Uuid_Categories)
	post.Categories = categories

	totalPages := (1 + perPage - 1) / perPage
	pages := make([]int, totalPages)
	for i := range pages {
		pages[i] = i + 1
	}
	return post, pages, nil
}
