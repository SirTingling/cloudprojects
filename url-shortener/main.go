package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	urlMap = make(map[string]string)
	mutex  = &sync.Mutex{}
)

func generateShortURL() string {
	b := make([]byte, 6) // 6 bytes = 48 bits
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:6]
}

func shortenURL(c *gin.Context) {
	var request struct {
		URL string `json:"url"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	short := generateShortURL()
	urlMap[short] = request.URL
	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8080/" + short})
}

func resolveURL(c *gin.Context) {
	short := c.Param("short")
	mutex.Lock()
	defer mutex.Unlock()

	original, exists := urlMap[short]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, original)
}

func main() {
	r := gin.Default()

	r.POST("/shorten", shortenURL)
	r.GET("/:short", resolveURL)

	r.Run(":8080")
}
