package handlers

import (
    "net/http"
    "encoding/json"
    "main/models"
    "fmt"
)


// Function to allow user to post some data
func PostData(w http.ResponseWriter, r *http.Request) {
	// Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content=Type", "application/json")
	fmt.Println("Post Data")

	// what if body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}

	var temp models.InputVars
	_ = json.NewDecoder(r.Body).Decode(&temp)

	json.NewEncoder(w).Encode(temp)
	fmt.Println("Keywords :", temp.Keywords)
	fmt.Println("Location :", temp.Location)
	fmt.Println("Time frame : ", temp.Timeframe)
}
