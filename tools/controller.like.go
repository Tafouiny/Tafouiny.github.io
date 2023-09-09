package Tools

import (
	"database/sql"

	Models "forum/models"
)

func LikePost(postID, userID int) {
	db := OpenDB()
	defer db.Close()
	// Faire des vérifications avant d'insert on peut faire un update
	if !UserHasAlreadyLiked(postID, userID) {
		_, err := db.Exec("INSERT INTO post_likes (id_post, id_user, like_type) VALUES (?, ?, ?)",
			postID, userID, 1)
		if err != nil {
			HandleError(err, "Insert post like")
			return
		}
	}
}

func DislikePost(postID, userID int) {
	db := OpenDB()
	defer db.Close()
	// Faire des vérifications avant d'insert on peut faire un update
	if !UserHasAlreadyDisliked(postID, userID) {
		_, err := db.Exec("INSERT INTO post_likes (id_post, id_user, like_type) VALUES (?, ?, ?)",
			postID, userID, -1)
		if err != nil {
			HandleError(err, "Insert post like")
			return
		}
	}
}

// Retirer son like laisser en neutre
func RemoveLike(postID, userID int) {
	db := OpenDB()
	defer db.Close()

	_, err := db.Exec("DELETE FROM post_likes WHERE id_post = ? AND id_user = ? AND like_type = 1", postID, userID)
	if err != nil {
		HandleError(err, "Deleting from post_likes")
		return
	}
}

func RemoveDislike(postID, userID int) {
	db := OpenDB()
	defer db.Close()

	_, err := db.Exec("DELETE FROM post_likes WHERE id_post = ? AND id_user = ? AND like_type = -1", postID, userID)
	if err != nil {
		HandleError(err, "Deleting from post_likes")
		return
	}
}

// Vérifie si l'utilisateur a déjà aimé un post
func UserHasAlreadyLiked(postID, userID int) bool {
	db := OpenDB()
	defer db.Close()
	var exists bool
	err := db.QueryRow("SELECT 1 FROM post_likes WHERE id_post = ? AND id_user = ? AND like_type = 1",
		postID, userID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		HandleError(err, "UserHasAlreadyLiked")
	}
	return exists
}

// Vérifie si l'utilisateur a déjà disliké un post
func UserHasAlreadyDisliked(postID, userID int) bool {
	db := OpenDB()
	defer db.Close()
	var exists bool
	err := db.QueryRow("SELECT 1 FROM post_likes WHERE id_post = ? AND id_user = ? AND like_type = -1",
		postID, userID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		HandleError(err, "UserHasAlreadyLiked")
	}
	return exists
}

func GetLikeFromPost(postID int) int {
	var Likes int
	db := OpenDB()
	defer db.Close()
	err := db.QueryRow("SELECT COUNT(*) FROM post_likes WHERE id_post = ? AND like_type = 1 ", postID).Scan(&Likes)
	if err != nil {
		HandleError(err, "GetLikeFromPost")
	}
	return Likes
}

func GetDislikesFromPost(postID int) int {
	var Likes int
	db := OpenDB()
	defer db.Close()
	err := db.QueryRow("SELECT COUNT(*) FROM post_likes WHERE id_post = ? AND like_type = -1 ", postID).Scan(&Likes)
	if err != nil {
		HandleError(err, "GetLikeFromPost")
	}
	return Likes
}

func GetUsersWithMostLikes() []Models.LikeUser {
	var users []Models.LikeUser
	db := OpenDB()
	defer db.Close()
	rows, err := db.Query(`
		SELECT p.id_user, COUNT(*) AS like_count, u.username  FROM posts p
		JOIN post_likes pl ON p.id_post = pl.id_post
		JOIN users u ON p.id_user = u.id_user
		WHERE like_type = 1 
		GROUP BY p.id_user
		ORDER BY like_count DESC LIMIT 5;
	`)
	if err != nil {
		HandleError(err, "Handle Post_likes errors ")
		return users
	}
	for rows.Next() {
		var user Models.LikeUser
		err := rows.Scan(&user.UserId, &user.LikesCount, &user.Username)
		if err != nil {
			return users
		}
		users = append(users, user)
	}
	defer rows.Close()
	return users
}
