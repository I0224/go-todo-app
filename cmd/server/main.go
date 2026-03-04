package main

import (
	"html/template"
	"log"
	"net/http"
	"todo-app/internal/controller"
	"todo-app/internal/db"
	"todo-app/internal/repository"
	"todo-app/internal/service"
)

func main() {
	// -----------------------------
	// 1) DB 初期化(接続, テーブル作成)
	// -----------------------------
	database, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// -----------------------------
	// 2) Repository / Service / Controller の組み立て
	// -----------------------------
	repo := &repository.TodoRepository{DB: database}
	svc := &service.TodoService{Repo: repo}

	// テンプレート読み込み(templates/index.html)
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	ctrl := &controller.TodoController{
		Service: svc,
		Tmpl:    tmpl,
	}

	// -----------------------------
	// 3) ルーティング(URL → Controller)
	// -----------------------------
	http.HandleFunc("/", ctrl.Index)
	http.HandleFunc("/add", ctrl.Add)
	http.HandleFunc("/delete", ctrl.Delete)
	http.HandleFunc("/toggle", ctrl.Toggle)
	http.HandleFunc("/update-title", ctrl.UpdateTitle)
	http.HandleFunc("/update-date", ctrl.UpdateDate)

	// 静的ファイル(CSS/JS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// -----------------------------
	// 4) サーバー起動
	// -----------------------------
	addr := ":8080"
	log.Println("server started:", "http://localhost"+addr)

	// ListenAndServe は基本的に「止まらない」ので、ここから下には通常来ない
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
