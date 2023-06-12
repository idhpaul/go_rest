package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


type PresignUrlCount struct {
	Count int 	`json:"count" binding:"required"`
}

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

		var requestBody PresignUrlCount
		c.Bind(&requestBody)

		res := create_presignURL(requestBody.Count)

		c.JSON(http.StatusOK, res)
	})
}


func main() {
	router := gin.Default()
	setRouter(router)
	_ = router.Run(":8080")
}
