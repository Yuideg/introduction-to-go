package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"

	"github.com/Yuideg/restaurntdb/delivery/http/handler"
	"github.com/Yuideg/restaurntdb/menu/repository"
	"github.com/Yuideg/restaurntdb/menu/service"
)

func main() {

	dbconn, err := sql.Open("postgres", "postgres://postgres:23782378@localhost/restaurantdb?sslmode=disable")

	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		panic(err)
	}
	

	tmpl := template.Must(template.ParseGlob("../../ui/templates/*"))

	categoryRepo := repository.NewCategoryRepositoryImpl(dbconn)
	categoryServ := service.NewCategoryServiceImpl(categoryRepo)
	adminCatgHandler := handler.NewAdminCategoryHandler(tmpl, categoryServ)
	menuHandler := handler.NewMenuHandler(tmpl, categoryServ)
     mux:=http.NewServeMux()
	fs := http.FileServer(http.Dir("ui/asset"))
	mux.Handle("ui/asset/", http.StripPrefix("ui/asset/", fs))
	mux.HandleFunc("/", menuHandler.Index)
	mux.HandleFunc("/about", menuHandler.About)
	mux.HandleFunc("/contact", menuHandler.Contact)
	mux.HandleFunc("/menu", menuHandler.Menu)
	mux.HandleFunc("/admin", menuHandler.Admin)

	mux.HandleFunc("/admin/categories/update", adminCatgHandler.AdminCategoriesUpdate)
	mux.HandleFunc("/admin/categories/delete", adminCatgHandler.AdminCategoriesDelete)

	http.ListenAndServe(":8181", mux)

}
