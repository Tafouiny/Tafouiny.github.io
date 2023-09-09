package Tools

import (
	"log"
	"time"

	Models "forum/models"
)

func GetSessions() ([]Models.Session, error) {
	var sessions []Models.Session
	db := OpenDB()
	defer db.Close()

	rows, err := db.Query("SELECT token, username, expiry FROM sessions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session Models.Session
		if err := rows.Scan(&session.Token, &session.Username, &session.Expiry); err != nil {
			return sessions, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func IfSessionExist(token string) (Models.Session, bool) {
	sessions, err := GetSessions()
	if err != nil {
		HandleError(err, "Fetching users sessions.")
		return Models.Session{}, false
	}

	for _, session := range sessions {
		if session.Token == token {
			return session, true
		}
	}
	return Models.Session{}, false
}

func UserHasAlreadyASession(username string) (Models.Session, bool) {
	sessions, err := GetSessions()
	if err != nil {
		HandleError(err, "Fetching users sessions.")
		return Models.Session{}, false
	}

	for _, session := range sessions {
		if session.Username == username {
			return session, true
		}
	}
	return Models.Session{}, false
}

func CreateSession(token, username string, expiry time.Time) {
	db := OpenDB()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO sessions(token, username, expiry) values(?,?,?)")
	if err != nil {
		HandleError(err, "preparing insertion of session")
		return

	}
	res, err := stmt.Exec(token, username, expiry)
	if err != nil {
		HandleError(err, "Excecuting insertion of an session")
		return
	}
	res.RowsAffected()
	log.Printf("Session with token as +> %s for username +> %s created and expiry in %v\n", token, username, expiry)
}

func DeleteSession(token string) {
	db := OpenDB()
	defer db.Close()
	stmt, err := db.Prepare("delete from sessions where token=?")
	HandleError(err, "preparing delete session")
	if err != nil {
		return
	}
	res, err := stmt.Exec(token)
	if err != nil {
		HandleError(err, "executing delete session")
		return
	}
	res.RowsAffected()
}

func UpdateSession(token string, newExpiry time.Time) {
	db := OpenDB()
	defer db.Close()
	stmt, err := db.Prepare("UPDATE sessions set expiry=? where token=?")
	if err != nil {
		HandleError(err, "preparing update session")
		return
	}
	res, err := stmt.Exec(newExpiry, token)
	if err != nil {
		HandleError(err, "executing update session")
		return
	}
	res.RowsAffected()
}

func GetIDUserFromSession(session Models.Session) int {
	users := GetUsers()
	for _, user := range users {
		if user.Username == session.Username {
			return user.ID
		}
	}
	return 0
}
