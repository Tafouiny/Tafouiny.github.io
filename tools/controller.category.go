package Tools

import (
	"log"

	Models "forum/models"
)

func GetCategories() []Models.Category {
	var category Models.Category
	var categories []Models.Category
	db := OpenDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		HandleError(err, "Fetching categories database.")
		return nil
	}
	for rows.Next() {
		err := rows.Scan(&category.ID, &category.Content)
		if err != nil {
			HandleError(err, "Fetching categories database.")
			return categories
		}
		categories = append(categories, category)
	}
	return categories
}

func CreateCategory(Content string) {
	db := OpenDB()
	stmt, err := db.Prepare("INSERT INTO categories(content) values(?)")
	if err != nil {
		HandleError(err, "preparing insertion of category")
		return
	}
	res, err := stmt.Exec(Content)
	if err != nil {
		HandleError(err, "Excecuting insertion of category")
		return
	}
	res.RowsAffected()
	log.Printf("category:%s created\n", Content)
	db.Close()
}
