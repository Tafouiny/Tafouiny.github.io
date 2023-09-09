package Tools

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	Models "forum/models"
)

func OpenDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Determine if the session has expired or not
func IsExpired(expiry time.Time) bool {
	return time.Now().After(expiry)
}

// Print the error while achieving a given task.
func HandleError(err error, task string) {
	if err != nil {
		log.Printf("Error While %s | more=> %v\n", task, err)
	}
}

func IsAuthenticated(r *http.Request) (Models.Session, bool) {
	c, err := r.Cookie(sessionCookieName)
	if err == nil {
		sessionToken := c.Value
		userSession, exists := IfSessionExist(sessionToken)
		if exists && !IsExpired(userSession.Expiry) {
			return userSession, true
		}
	}
	return Models.Session{}, false
}

func OnlySpace(s string) bool {
	for _, v := range s {
		if !unicode.IsSpace(v) {
			return false
		}
	}
	return true
}

func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}
