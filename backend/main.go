package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"example.com/patient"
	"example.com/triage"
)

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

	p:=patient.NewPatients()
	t:=triage.NewTriages()

	http.Handle("/patient", middleware(http.HandlerFunc(p.ServeHTTP)))
	http.Handle("/patient/", middleware(http.HandlerFunc(p.ServeHTTP)))
	http.Handle("/patients", middleware(http.HandlerFunc(p.ServeHTTP)))
	http.Handle("/triage", middleware(http.HandlerFunc(t.ServeHTTP)))
	http.Handle("/triage/", middleware(http.HandlerFunc(t.ServeHTTP)))
	http.Handle("/triages", middleware(http.HandlerFunc(t.ServeHTTP)))

	err = http.ListenAndServe(string(os.Getenv("SERVER")+":"+os.Getenv("PORT")), nil)
	if err == http.ErrServerClosed {
		log.Fatal("Backend server closed")
	} else if err != nil {
		log.Fatalf("Backend server:Error occured %v", err.Error())
		os.Exit(1)
	}

}