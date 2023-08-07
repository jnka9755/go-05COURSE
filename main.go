package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jnka9755/go-05COURSE/internal/course"
	"github.com/jnka9755/go-05COURSE/package/boostrap"
	"github.com/joho/godotenv"
)

func main() {

	router := mux.NewRouter()

	_ = godotenv.Load()

	log := boostrap.InitLooger()

	db, err := boostrap.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		log.Fatal("paginator limit defauly is required")
	}

	courseRepository := course.NewRepository(log, db)
	courseBusiness := course.NewBusiness(log, courseRepository)
	courseController := course.MakeEndpoints(courseBusiness, course.Config{LimPageDef: pagLimDef})

	router.HandleFunc("/courses", courseController.Create).Methods("POST")
	router.HandleFunc("/courses", courseController.GetAll).Methods("GET")
	router.HandleFunc("/courses/{id}", courseController.Get).Methods("GET")
	router.HandleFunc("/courses/{id}", courseController.Delete).Methods("DELETE")
	router.HandleFunc("/courses/{id}", courseController.Update).Methods("PATCH")

	port := os.Getenv("PORT")

	address := fmt.Sprintf("127.0.0.1:%s", port)

	server := http.Server{
		Handler:      http.TimeoutHandler(router, time.Second*5, "Timeout!"),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	error := server.ListenAndServe()

	if err != nil {
		log.Fatal(error)
	}
}
