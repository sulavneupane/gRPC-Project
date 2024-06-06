package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func main() {
	r := gin.Default()
	r.POST("/upload", uploadFile)
	r.Run(":9001")
}

func uploadFile(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		fmt.Println("error while uploading file", err)
		return
	}

	filePath := filepath.Join("uploads", f.Filename)
	if err := c.SaveUploadedFile(f, filePath); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"error": "saved successfully"})
}
