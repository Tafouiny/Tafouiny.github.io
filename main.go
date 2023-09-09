package main

import (
	"log"
	"net/http"

	Tools "forum/tools"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	Tools.InitDb()

	for _, route := range Tools.Routes {
		http.HandleFunc(route.Path, Tools.MiddlewareError(route.Path, route.Handler))
	}

	/*******************/
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./views/static/css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./views/static/js/"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./views/static/assets/"))))
	/*******************/
	log.Println("Server started and listenning on port", Tools.Port)
	log.Println("http://localhost" + Tools.Port)
	log.Fatal(http.ListenAndServe(Tools.Port, nil))
}
