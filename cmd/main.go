package main

import (
	"database/sql"
	"net/http"
	"newsfeed/platform/newsfeed"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qkgo/yin"
)

func main() {
	db, _ := sql.Open("sqlite3", "./newsfeed.db")

	feed := newsfeed.NewFeed(db)

	// feed.Add(newsfeed.Item{
	// 	Content: "hello",
	// })

	// items := feed.Get()
	// fmt.Println(items)

	r := chi.NewRouter()

	r.Use(yin.SimpleLogger)

	r.Get("/posts", func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		items := feed.Get()
		res.SendJSON(items)
	})

	r.Post("/posts", func(w http.ResponseWriter, r *http.Request) {
		res, req := yin.Event(w, r)
		body := map[string]string{}
		req.BindBody(&body)
		item := newsfeed.Item{
			Content: body["content"],
		}
		feed.Add(item)
		res.SendStatus(204)
	})

	http.ListenAndServe(":3000", r)

}
