package Tools

import (
	"database/sql"
	"log"

	Models "forum/models"

	"golang.org/x/crypto/bcrypt"
)

func GetUsers() []Models.User {
	var user Models.User
	var users []Models.User
	db := OpenDB()
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		HandleError(err, "Fetching users database.")
		return nil
	}
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Username)
		if err != nil {
			HandleError(err, "Fetching users database.")
			return users
		}
		users = append(users, user)
	}
	db.Close()
	return users
}

func GetUserByUsername(username string) int {
	var user_id int

	db := OpenDB()
	row := db.QueryRow("SELECT id_user FROM users WHERE username = ?", username)

	err := row.Scan(&user_id)
	if err != nil {
		return user_id
	}

	db.Close()
	return user_id
}

func IfUserExist(username, email string) (Models.User, bool) {
	db := OpenDB()

	users := GetUsers()

	for _, user := range users {
		if user.Username == username || user.Email == email {
			return user, true
		}
	}
	db.Close()
	return Models.User{}, false
}

func CreateUser(email, password, username string) {
	db := OpenDB()

	stmt, err := db.Prepare("INSERT INTO users(email, password, username) values(?,?,?)")
	if err != nil {
		HandleError(err, "preparing insertion of user")
		return
	}
	res, err := stmt.Exec(email, password, username)
	if err != nil {
		HandleError(err, "Excecuting insertion of user")
		return
	}
	res.RowsAffected()
	log.Printf("email:%s username:%s ; user created\n", email, username)
	db.Close()
}

func CheckCredentials(email, password string) (Models.User, bool) {
	users := GetUsers()

	var hashedPassword string

	for _, user := range users {

		hashedPassword = user.Password
		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

		if user.Email == email && err == nil {
			return user, true
		}
	}

	return Models.User{}, false
}

// TODO: Delete very soon
func ReadUser(rows *sql.Rows) []Models.User {
	var user Models.User
	var users []Models.User
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Username)
		if err != nil {
			HandleError(err, "Scanning tables users from db")
			return users
		}
		users = append(users, user)
	}
	return users
}
