package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (tpl   *template.Template)

func init() {
	tpl = template.Must(template.ParseGlob("./templates/*.html"))
}

type Page struct {
	Body  []byte
	Title string
	Data  interface{}
	Backend string
}
func LoadPage(file string) (*Page, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	var body []byte
	_, err = f.Read(body)
	if err != nil {
		return nil, err
	}
	return &Page{Body: body,Backend: "http://"+os.Getenv("BACKEND")+":"+os.Getenv("BACKPORT")},  nil
}

func RenderTemplate(w http.ResponseWriter, file string, page *Page) {
	err := tpl.ExecuteTemplate(w, file, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func PatientHandler(w http.ResponseWriter, r *http.Request) {
	file := "patient.html"
	filePath := "templates/" + file
	pageName := "Patient Page"
	page, err := LoadPage(filePath)
	if err != nil {
		page = &Page{Title: pageName}
	}
	page.Title = pageName
	RenderTemplate(w, file, page)
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	file := "index.html"
	filePath := "templates/" + file
	pageName := "Home Page"
	page, err := LoadPage(filePath)
	if err != nil {
		page = &Page{Title: pageName}
	}
	page.Title = pageName
	RenderTemplate(w, file, page)
}
func TriageHandler(w http.ResponseWriter, r *http.Request) {
	file := "triage.html"
	filePath := "templates/" + file
	pageName := "Triages Page"
	page, err := LoadPage(filePath)
	if err != nil {
		page = &Page{Title: pageName}
	}
	page.Title = pageName
	RenderTemplate(w, file, page)
}
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			return
		}
		next.ServeHTTP(w, r)
})
}
func main(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the .env file")
	}
	fmt.Printf("Server running on  %v : %v",os.Getenv("SERVER"),os.Getenv("PORT"))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(os.Getenv("TMP")))))
	http.Handle("/", middleware(http.HandlerFunc(IndexHandler)))
	http.Handle("/patient", middleware(http.HandlerFunc(PatientHandler)))
	http.Handle("/triage", middleware(http.HandlerFunc(TriageHandler)))

	err = http.ListenAndServe(string(os.Getenv("SERVER")+":"+os.Getenv("PORT")), nil)
	if err == http.ErrServerClosed {
		log.Fatal("Backend server closed")
	} else if err != nil {
		log.Fatalf("Backend server:Error occured %v", err.Error())
		os.Exit(1)
	}

}