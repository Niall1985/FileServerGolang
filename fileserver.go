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


package fileserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Authentication middleware for upload
func uploadAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Basic Authentication
		user, pass, ok := r.BasicAuth()
		if !ok || !validateUploadCredentials(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Upload"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func validateUploadCredentials(user, pass string) bool {
	// Validate credentials for uploading (this is a simple example, replace with real validation)
	return user == "uploader" && pass == "uploadpassword"
}

// General logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
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

func StartServer(addr string) {
	// Ensure the uploads directory exists
	os.MkdirAll("uploads", os.ModePerm)

	// Apply middlewares and handlers
	http.Handle("/upload", loggingMiddleware(uploadAuthMiddleware(http.HandlerFunc(UploadHandler))))
	http.Handle("/download/", loggingMiddleware(http.HandlerFunc(DownloadHandler)))
	http.Handle("/list", loggingMiddleware(http.HandlerFunc(ListHandler)))
	http.Handle("/delete/", loggingMiddleware(http.HandlerFunc(DeleteHandler)))

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
