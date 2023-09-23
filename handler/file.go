package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type File struct{}

func NewFile() *File {
	return &File{}
}

// UploadFile is a function that handles the upload of a single file
func (f *File) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	// Define the path where the file will be saved
	fileName := uuid.New().String()
	filePath := filepath.Join("uploads", fileName+".png")

	// Save the file to the defined path
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, struct {
		File string
	}{
		fileName,
	})
}

// GetFile is a function that retrieves a file from the server
func (f *File) GetFile(c *gin.Context) {
	fileName := c.Param("file")

	// Define the path of the file to be retrieved
	filePath := filepath.Join("uploads", fileName+".png")
	// Open the file
	fileData, err := os.Open(filePath)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}
	defer fileData.Close()

	// Read the first 512 bytes of the file to determine its content type
	fileHeader := make([]byte, 512)
	_, err = fileData.Read(fileHeader)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}
	fileContentType := http.DetectContentType(fileHeader)
	// Get the file info
	fileInfo, err := fileData.Stat()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}
	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", fileContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.File(filePath)
}

// DeleteFile is a function that deletes a file from the server and its metadata from the database
func (f *File) DeleteFile(c *gin.Context) {
	fileName := c.Param("file")

	// Define the path of the file to be deleted
	filePath := filepath.Join("uploads", fileName+".png")
	// Delete the file from the server
	err := os.Remove(filePath)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, struct{}{})
}
