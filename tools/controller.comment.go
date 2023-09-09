package Tools

import (
	"database/sql"

	Models "forum/models"
)

func ReadComment(rows *sql.Rows) []Models.Comment {
	var comment Models.Comment
	var comments []Models.Comment
	for rows.Next() {
		err := rows.Scan(&comment.ID, &comment.Content)
		if err != nil {
			HandleError(err, "Scanning tables comments from db")
			return comments
		}
		comments = append(comments, comment)
	}
	return comments
}

func CreateComment(id_post, id_user int, comment_content string) {
	db := OpenDB()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO comments(id_post, id_user, comment_content) values(?,?,?)")
	if err != nil {
		HandleError(err, "preparing insertion of comment")
		return
	}
	_, err = stmt.Exec(id_post, id_user, comment_content)
	if err != nil {
		HandleError(err, "Excecuting insertion of comment")
		return
	}
}

func GetComments() []Models.Comment {
	var comment Models.Comment
	var comments []Models.Comment
	db := OpenDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM comments")
	if err != nil {
		HandleError(err, "Fetching categories database.")
		return nil
	}
	for rows.Next() {
		err := rows.Scan(&comment.ID, &comment.ID_Post, &comment.Content)
		if err != nil {
			HandleError(err, "Fetching categories database.")
			return comments
		}
		comments = append(comments, comment)
	}
	return comments
}
