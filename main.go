package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var host = "localhost"
var port = "3306"
var user = "gouser"
var password = "gopasswd"

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("tps/*"))
}

//HandleError is a function...
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

//Registro is a function...
func Registro(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "registro.html", nil)
	HandleError(w, err)
}

//Login is a function...
func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "login.html", nil)
	HandleError(w, err)
}

//CreaServicio is a function...
func CreaServicio(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "crea.html", nil)
	HandleError(w, err)
}

//Index is a function...
func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	HandleError(w, err)
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/registro.aspx", Registro)
	router.POST("/registro.aspx", Registro)
	router.GET("/crea.php", CreaServicio)
	router.POST("/crea.php", CreaServicio)
	router.GET("/login.aspx", Login)
	router.POST("/login.aspx", Login)

	log.Fatal(http.ListenAndServe(":8800", router))
}
