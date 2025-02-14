package method

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ocr() {

	router := gin.Default()
	// Endpoint
	router.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image file"})
			return
		}
		defer file.Close()

		tempFilePath := "temp_" + header.Filename
		outFile, err := os.Create(tempFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		defer outFile.Close()
		io.Copy(outFile, file)

		//handing over of the image to the python file
		jsonResponse, err := sendToPythonService(tempFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Delete the temporary file
		os.Remove(tempFilePath)

		c.Data(http.StatusOK, "application/json", jsonResponse)
	})

	router.Run(":8080")
}

func sendToPythonService(imagePath string) ([]byte, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy image data: %v", err)
	}

	writer.Close()

	req, err := http.NewRequest("POST", "http://localhost:8080/analyze", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to Python service failed: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return responseData, nil
}

func main() {
	ocr()
}
