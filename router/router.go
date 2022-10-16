package router

import (
	"net/http"

	"github.com/anudeep652/mongoapi/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// router.HandleFunc("/api/movies", controller.GetAllMovies).Methods("GET")
	// router.HandleFunc("/api/movie", controller.CreateMovie).Methods("POST")
	// router.HandleFunc("/api/movie/{id}", controller.MarkAsWatched).Methods("PUT")
	// router.HandleFunc("/api/movie/{id}", controller.GetAMovie).Methods("GET")
	// router.HandleFunc("/api/movie/{id}", controller.DeleteAMovie).Methods("DELETE")
	// router.HandleFunc("/api/deleteallmovie", controller.DeleteAllMovie).Methods("DELETE")

	// components

	router.HandleFunc("/components", controller.GetAll).Methods("GET").Subrouter().Use(loggingMiddleware)
	router.HandleFunc("/component", controller.CreateComponent).Methods("POST") //.subrouter().use() for middleware
	router.HandleFunc("/component/{id}", controller.DeleteComponent).Methods("DELETE")
	router.HandleFunc("/component/{componentName}/{batchName}", controller.DeleteABatch).Methods("DELETE")
	router.HandleFunc("/component/{componentName}", controller.UpdateBatch).Methods("PUT")
	router.HandleFunc("/component/{componentName}/{batchName}", controller.UpdateProcess).Methods("PUT")

	return router

}

// middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		next.ServeHTTP(response, request)
	})
}
