// // package fileserver

// // import (
// // 	"fmt"
// // 	"io"
// // 	"net/http"
// // 	"os"
// // 	"path/filepath"
// // 	"strings"
// // )

// // // UploadHandler handles file uploads
// // func UploadHandler(w http.ResponseWriter, r *http.Request) {
// // 	if r.Method != http.MethodPost {
// // 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// // 		return
// // 	}

// // 	file, handler, err := r.FormFile("file")
// // 	if err != nil {
// // 		http.Error(w, "Failed to get file", http.StatusBadRequest)
// // 		return
// // 	}
// // 	defer file.Close()

// // 	savePath := filepath.Join("uploads", handler.Filename)
// // 	out, err := os.Create(savePath)
// // 	if err != nil {
// // 		http.Error(w, "Failed to save file", http.StatusInternalServerError)
// // 		return
// // 	}
// // 	defer out.Close()

// // 	_, err = io.Copy(out, file)
// // 	if err != nil {
// // 		http.Error(w, "Failed to save file", http.StatusInternalServerError)
// // 		return
// // 	}

// // 	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
// // }

// // // DownloadHandler handles file downloads
// // func DownloadHandler(w http.ResponseWriter, r *http.Request) {
// // 	fileName := strings.TrimPrefix(r.URL.Path, "/download/")
// // 	filePath := filepath.Join("uploads", fileName)

// // 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
// // 		http.Error(w, "File not found", http.StatusNotFound)
// // 		return
// // 	}

// // 	http.ServeFile(w, r, filePath)
// // }

// // // ListHandler lists all files in the upload directory
// // func ListHandler(w http.ResponseWriter, r *http.Request) {
// // 	files, err := os.ReadDir("uploads")
// // 	if err != nil {
// // 		http.Error(w, "Failed to list files", http.StatusInternalServerError)
// // 		return
// // 	}

// // 	for _, file := range files {
// // 		fmt.Fprintln(w, file.Name())
// // 	}
// // }

// // // DeleteHandler deletes a specified file
// // func DeleteHandler(w http.ResponseWriter, r *http.Request) {
// // 	fileName := strings.TrimPrefix(r.URL.Path, "/delete/")
// // 	filePath := filepath.Join("uploads", fileName)

// // 	if err := os.Remove(filePath); err != nil {
// // 		if os.IsNotExist(err) {
// // 			http.Error(w, "File not found", http.StatusNotFound)
// // 		} else {
// // 			http.Error(w, "Failed to delete file", http.StatusInternalServerError)
// // 		}
// // 		return
// // 	}

// // 	fmt.Fprintf(w, "File deleted successfully: %s\n", fileName)
// // }

// // func StartServer(addr string) {
// // 	// Ensure the uploads directory exists
// // 	os.MkdirAll("uploads", os.ModePerm)

// // 	http.HandleFunc("/upload", UploadHandler)
// // 	http.HandleFunc("/download/", DownloadHandler)
// // 	http.HandleFunc("/list", ListHandler)
// // 	http.HandleFunc("/delete/", DeleteHandler)

// // 	fmt.Println("Starting server on :8080")
// // 	http.ListenAndServe(":8080", nil)
// // }

// // //test

// package fileserver

// import (
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/rand"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strings"
// )

// // Encryption key (must be 32 bytes for AES-256)
// var encryptionKey = []byte("a very very very very secret key!!!!") // Change this to your own key

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

// 	// Read the file content
// 	data, err := io.ReadAll(file)
// 	if err != nil {
// 		http.Error(w, "Failed to read file", http.StatusInternalServerError)
// 		return
// 	}

// 	// Encrypt the file content
// 	encryptedData, err := encrypt(data)
// 	if err != nil {
// 		http.Error(w, "Failed to encrypt file", http.StatusInternalServerError)
// 		return
// 	}

// 	savePath := filepath.Join("uploads", handler.Filename)
// 	err = os.WriteFile(savePath, encryptedData, 0644)
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

// 	encryptedData, err := os.ReadFile(filePath)
// 	if err != nil {
// 		http.Error(w, "Failed to read file", http.StatusInternalServerError)
// 		return
// 	}

// 	// Decrypt the file content
// 	decryptedData, err := decrypt(encryptedData)
// 	if err != nil {
// 		http.Error(w, "Failed to decrypt file", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
// 	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
// 	w.Write(decryptedData)
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

// // StartServer starts the file server
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

// // Encryption and Decryption functions

// func encrypt(data []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(encryptionKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
// 		return nil, err
// 	}

// 	return gcm.Seal(nonce, nonce, data, nil), nil
// }

// func decrypt(data []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(encryptionKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nonceSize := gcm.NonceSize()
// 	if len(data) < nonceSize {
// 		return nil, fmt.Errorf("ciphertext too short")
// 	}

// 	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
// 	return gcm.Open(nil, nonce, ciphertext, nil)
// }

package fileserver

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Encryption key (must be 16, 24, or 32 bytes long)
var encryptionKey = []byte("example key 1234")

// UploadHandler handles file uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		log.Printf("Failed to get file: %v", err)
		return
	}
	defer file.Close()

	savePath := filepath.Join("uploads", handler.Filename)
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		log.Printf("Failed to save file: %v", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		log.Printf("Failed to save file: %v", err)
		return
	}

	// Encrypt the file
	encryptedPath := savePath + ".enc"
	if err := encryptFile(savePath, encryptedPath, encryptionKey); err != nil {
		http.Error(w, "Failed to encrypt file", http.StatusInternalServerError)
		log.Printf("Failed to encrypt file: %v", err)
		return
	}

	// Remove the original file after encryption
	if err := os.Remove(savePath); err != nil {
		log.Printf("Failed to remove original file: %v", err)
	}

	fmt.Fprintf(w, "File uploaded and encrypted successfully: %s\n", handler.Filename)
}

// EncryptFile encrypts the input file and saves the result to the output file
func encryptFile(inputFile, outputFile string, key []byte) error {
	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(inputBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("failed to read IV: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], inputBytes)

	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()

	if _, err := output.Write(ciphertext); err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}

	return nil
}

// DownloadHandler handles file downloads
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/download/")
	filePath := filepath.Join("uploads", fileName+".enc")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		log.Printf("File not found: %s", filePath)
		return
	}

	http.ServeFile(w, r, filePath)
}

// ListHandler lists all files in the upload directory
func ListHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("uploads")
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		log.Printf("Failed to list files: %v", err)
		return
	}

	for _, file := range files {
		fmt.Fprintln(w, file.Name())
	}
}

// DeleteHandler deletes a specified file
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/delete/")
	filePath := filepath.Join("uploads", fileName+".enc")

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete file", http.StatusInternalServerError)
			log.Printf("Failed to delete file: %v", err)
		}
		return
	}

	fmt.Fprintf(w, "File deleted successfully: %s\n", fileName)
}

// StartServer starts the file server
func StartServer(addr string) {
	// Ensure the uploads directory exists
	os.MkdirAll("uploads", os.ModePerm)

	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/download/", DownloadHandler)
	http.HandleFunc("/list", ListHandler)
	http.HandleFunc("/delete/", DeleteHandler)

	fmt.Println("Starting server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
