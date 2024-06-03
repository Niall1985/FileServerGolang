package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// BasicAuthMiddleware is the authentication middleware
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract username and password from the Authorization header
		user, pass, ok := r.BasicAuth()
		// Check if the credentials are provided and valid
		if !ok || !validateCredentials(user, pass) {
			// Prompt for credentials if not valid
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		// Call the next handler if credentials are valid
		next.ServeHTTP(w, r)
	})
}

// validateCredentials validates the provided credentials
func validateCredentials(user, pass string) bool {
	// Replace with real validation logic
	return user == "admin" && pass == "password"
}

// LoggingMiddleware logs the details of each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("%s %s %s %s\n", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}

// UploadHandler handles file uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	savePath := filepath.Join("uploads", handler.Filename)
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

// DownloadHandler handles file downloads
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/download/")
	filePath := filepath.Join("uploads", fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}

// ListHandler lists all files in the upload directory
func ListHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("uploads")
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		fmt.Fprintln(w, file.Name())
	}
}

// DeleteHandler deletes a specified file
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/delete/")
	filePath := filepath.Join("uploads", fileName)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "File deleted successfully: %s\n", fileName)
}

// StartServer initializes the server and sets up the routes
func StartServer(addr string) {
	// Ensure the uploads directory exists
	os.MkdirAll("uploads", os.ModePerm)

	http.Handle("/upload", LoggingMiddleware(BasicAuthMiddleware(http.HandlerFunc(UploadHandler))))
	http.Handle("/download/", LoggingMiddleware(BasicAuthMiddleware(http.HandlerFunc(DownloadHandler))))
	http.Handle("/list", LoggingMiddleware(BasicAuthMiddleware(http.HandlerFunc(ListHandler))))
	http.Handle("/delete/", LoggingMiddleware(BasicAuthMiddleware(http.HandlerFunc(DeleteHandler))))

	fmt.Println("Starting server on " + addr)
	http.ListenAndServe(addr, nil)
}

func main() {
	StartServer(":8080")
}

// package fileserver

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strings"
// )

// // UploadHandler handles file uploads
// func UploadHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	file, handler, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, "Failed to get file", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	savePath := filepath.Join("uploads", handler.Filename)
// 	out, err := os.Create(savePath)
// 	if err != nil {
// 		http.Error(w, "Failed to save file", http.StatusInternalServerError)
// 		return
// 	}
// 	defer out.Close()

// 	_, err = io.Copy(out, file)
// 	if err != nil {
// 		http.Error(w, "Failed to save file", http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
// }

// // DownloadHandler handles file downloads
// func DownloadHandler(w http.ResponseWriter, r *http.Request) {
// 	fileName := strings.TrimPrefix(r.URL.Path, "/download/")
// 	filePath := filepath.Join("uploads", fileName)

// 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 		http.Error(w, "File not found", http.StatusNotFound)
// 		return
// 	}

// 	http.ServeFile(w, r, filePath)
// }

// // ListHandler lists all files in the upload directory
// func ListHandler(w http.ResponseWriter, r *http.Request) {
// 	files, err := os.ReadDir("uploads")
// 	if err != nil {
// 		http.Error(w, "Failed to list files", http.StatusInternalServerError)
// 		return
// 	}

// 	for _, file := range files {
// 		fmt.Fprintln(w, file.Name())
// 	}
// }

// // DeleteHandler deletes a specified file
// func DeleteHandler(w http.ResponseWriter, r *http.Request) {
// 	fileName := strings.TrimPrefix(r.URL.Path, "/delete/")
// 	filePath := filepath.Join("uploads", fileName)

// 	if err := os.Remove(filePath); err != nil {
// 		if os.IsNotExist(err) {
// 			http.Error(w, "File not found", http.StatusNotFound)
// 		} else {
// 			http.Error(w, "Failed to delete file", http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	fmt.Fprintf(w, "File deleted successfully: %s\n", fileName)
// }

// func StartServer(addr string) {
// 	// Ensure the uploads directory exists
// 	os.MkdirAll("uploads", os.ModePerm)

// 	http.HandleFunc("/upload", UploadHandler)
// 	http.HandleFunc("/download/", DownloadHandler)
// 	http.HandleFunc("/list", ListHandler)
// 	http.HandleFunc("/delete/", DeleteHandler)

// 	fmt.Println("Starting server on :8080")
// 	http.ListenAndServe(":8080", nil)
// }
