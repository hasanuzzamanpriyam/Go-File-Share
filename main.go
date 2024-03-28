package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Check if the request Content-Type is multipart/form-data
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB maximum file size
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error parsing form:", err)
		return
	}

	// Retrieve the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error retrieving file:", err)
		return
	}
	defer file.Close()

	// Create a new file in the uploads directory
	f, err := os.OpenFile("uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	// Copy the data to the new created file
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error copying file data:", err)
		return
	}

	fmt.Fprintf(w, "Successfully uploaded file %s", handler.Filename)
}

func main() {
	// serve static files from the static directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// handle file uploads
	http.HandleFunc("/upload", uploadFile)

	// start the server
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
