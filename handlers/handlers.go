package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func enableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	
	// Allow requests from any origin
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	// Allow specified HTTP methods
	
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	
	// Allow specified headers
	
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
	
	// Continue with the next handler
	
	next.ServeHTTP(w, r)
	
	})
	
}

func NewRouter() *mux.Router {

    router := mux.NewRouter()

    router.Use(enableCORS)
    
    router.HandleFunc("/sendJSON", sendJSONFileContent).Methods("GET")
    router.HandleFunc("/SearchDawn/{keywords}", getKeywords).Methods("GET")
    router.HandleFunc("/SearchTribune/{location}", getLocation).Methods("GET")
    router.HandleFunc("/SearchNews/{timeframe}", getTimeFrame).Methods("GET")
    router.HandleFunc("/getData/{initialData}", getInitialData).Methods("GET")
    
    router.HandleFunc("/", PostData).Methods("POST")
    router.HandleFunc("/searchData", searchDataHandler).Methods("POST")
    router.HandleFunc("/PostedData", PostData).Methods("POST")
    router.HandleFunc("/keywords", GetKeywords).Methods("POST")
    router.HandleFunc("/login", Login).Methods("POST")
    router.HandleFunc("/signup", Signup).Methods("POST")

    return router
}
