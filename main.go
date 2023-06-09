package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


type PresignURLs struct {
	Count int    `json:"count"`
	Urls  []UrlData `json:"urls"`
}
type UrlData struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}


func setRouter(router *gin.Engine) {

	router.POST("/presignURL", func(c *gin.Context) {

		print(c.Request.Header)
		print(c.Request.Body)

		
		res := create_presignURL(8)

		c.JSON(http.StatusOK, res)
	})
}


func main() {
	router := gin.Default()
	setRouter(router)
	_ = router.Run(":8080")
}
