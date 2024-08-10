package handlers

import (
    "net/http"
    "os"
    "io"
  
)

func sendJSONFileContent(w http.ResponseWriter, r *http.Request) {
    // Open the JSON file
    file, err := os.Open("../know-gr/keywords.json")
    if err != nil {
        // Handle error if unable to open the file
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer file.Close()

    // Set the appropriate content type
    w.Header().Set("Content-Type", "application/json")

    // Write the file content to the response writer
    _, err = io.Copy(w, file)
    if err != nil {
        // Handle error if unable to copy file content to response
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

