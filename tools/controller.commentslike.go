package Tools

import "database/sql"

func LikeComment(commentID, userID int) {
	db := OpenDB()
	defer db.Close()
	// Faire des vérifications avant d'insert on peut faire un update
	if !UserHasAlreadyLikedComment(commentID, userID) {
		_, err := db.Exec("INSERT INTO comment_likes (id_comment, id_user, like_type) VALUES (?, ?, ?)",
			commentID, userID, 1)
		if err != nil {
			HandleError(err, "Insert comment_likes")
			return
		}
	}
}

func DislikeComment(commentID, userID int) {
	db := OpenDB()
	defer db.Close()
	// Faire des vérifications avant d'insert on peut faire un update
	if !UserHasAlreadyDisliked(commentID, userID) {
		_, err := db.Exec("INSERT INTO comment_likes (id_comment, id_user, like_type) VALUES (?, ?, ?)",
			commentID, userID, -1)
		if err != nil {
			HandleError(err, "Insert comment_likes")
			return
		}
	}
}

// Retirer son like laisser en neutre
func RemoveLikeFromComment(commentID, userID int) {
	db := OpenDB()
	defer db.Close()

	_, err := db.Exec("DELETE FROM comment_likes WHERE id_comment = ? AND id_user = ? AND like_type = 1", commentID, userID)
	if err != nil {
		HandleError(err, "Deleting from post_likes")
		return
	}
}

func RemoveDislikeFromComment(commentID, userID int) {
	db := OpenDB()
	defer db.Close()

	_, err := db.Exec("DELETE FROM comment_likes WHERE id_comment = ? AND id_user = ? AND like_type = -1", commentID, userID)
	if err != nil {
		HandleError(err, "Deleting from post_likes")
		return
	}
}

// Vérifie si l'utilisateur a déjà aimé un post
func UserHasAlreadyLikedComment(postID, userID int) bool {
	db := OpenDB()
	defer db.Close()
	var exists bool
	err := db.QueryRow("SELECT 1 FROM comment_likes WHERE id_comment = ? AND id_user = ? AND like_type = 1",
		postID, userID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		HandleError(err, "UserHasAlreadyLikedComment")
	}
	return exists
}

func UserHasAlreadyDislikedComment(commentID, userID int) bool {
	db := OpenDB()
	defer db.Close()
	var exists bool
	err := db.QueryRow("SELECT 1 FROM comment_likes WHERE id_comment = ? AND id_user = ? AND like_type = -1",
		commentID, userID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		HandleError(err, "UserHasAlreadyLiked")
	}
	return exists
}

func GetLikesFromComment(commentID int) int {
	var Likes int
	db := OpenDB()
	defer db.Close()
	err := db.QueryRow("SELECT COUNT(*) FROM comment_likes WHERE id_comment = ? AND like_type = 1 ", commentID).Scan(&Likes)
	if err != nil {
		HandleError(err, "GetLikeFromPost")
	}
	return Likes
}

func GetDislikesFromComment(commentID int) int {
	var Dislikes int
	db := OpenDB()
	defer db.Close()
	err := db.QueryRow("SELECT COUNT(*) FROM comment_likes WHERE id_comment = ? AND like_type = -1 ", commentID).Scan(&Dislikes)
	if err != nil {
		HandleError(err, "GetLikeFromPost")
	}
	return Dislikes
}
