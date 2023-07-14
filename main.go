package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// //////////////////////////
type NeedEnhance struct {
	Count int `json:"count" binding:"required"`
	Retry int `json:"retry"`
}

type PreSignEnhance struct {
	Count int           `json:"count"`
	Urls  []EnhanceUrls `json:"urls"`
}
type EnhanceUrls struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// //////////////////////////
type NeedAnalyze struct {
	Count int `json:"count" binding:"required"`
	Retry int `json:"retry"`
}

type PreSignAnalyze struct {
	Count    int           `json:"count"`
	UrlJsons []AnalyzeUrls `json:"urljsons"`
}
type AnalyzeUrls struct {
	OriginalUrl        string `json:"originalurl"`
	OriginalOutputJson string `json:"originalouputjson"`
	InputUrl           string `json:"inputurl"`
	OutputJson         string `json:"outputjson"`
}

// ///////////////////////////
type NeedAnalyzeJson struct {
	Index int `json:"index"`
	Retry int `json:"retry"`
}
type AnalyzeJson struct {
	OriginalAnalyzeJsonData string `json:"originalAnalyzejson"`
	AnalyzeJsonData         string `json:"analyzejson"`
}

// //////////////////////////
type NeedEqualize struct {
	Count int `json:"count" binding:"required"`
	Retry int `json:"retry"`
}

type PreSignEqualize struct {
	Count int            `json:"count"`
	Urls  []EqualizeUrls `json:"urls"`
}
type EqualizeUrls struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// //////////////////////////
type NeedSTT struct {
	Index      int  `json:"index"`
	IsOriginal bool `json:"isOriginal"`
}
type STTStatus struct {
	Result string `json:"result"`
}

// //////////////////////////
type UploadExcel struct {
	FileName string `json:"fileName"`
	FileData string `json:"fileData"`
}

func setRouter(router *gin.Engine) {

	router.POST("/presignEnhance", func(c *gin.Context) {

		//print(c.Request.Header)
		//print(c.Request.Body)

		var requestBody NeedEnhance
		c.Bind(&requestBody)

		res := create_PreSignEnhance(requestBody.Count)

		c.JSON(http.StatusOK, res)
	})

	router.POST("/presignAnalyze", func(c *gin.Context) {

		//print(c.Request.Header)
		//print(c.Request.Body)

		var requestBody NeedAnalyze
		c.Bind(&requestBody)

		var res PreSignAnalyze

		if requestBody.Retry == 0 {
			res = create_PreSignAnalyze(requestBody.Count)
		} else {
			res = create_PreSignAnalyzeRetry(requestBody.Count, requestBody.Retry)
		}

		c.JSON(http.StatusOK, res)
	})

	router.POST("/getAnalyzeJson", func(c *gin.Context) {

		//print(c.Request.Header)
		//print(c.Request.Body)

		var requestBody NeedAnalyzeJson
		c.Bind(&requestBody)

		var res AnalyzeJson

		if requestBody.Retry == 0 {
			res = create_AnalyzeJson(requestBody.Index)
		} else {
			//res = create_PreSignAnalyzeRetry(requestBody.Count, requestBody.Retry)
		}

		c.JSON(http.StatusOK, res)
	})

	router.POST("/presignEqualize", func(c *gin.Context) {

		//print(c.Request.Header)
		//print(c.Request.Body)

		var requestBody NeedEqualize
		c.Bind(&requestBody)

		var res PreSignEqualize

		if requestBody.Retry == 0 {
			res = create_PreSignEqualize(requestBody.Count)
		} else {
			//res = create_PreSignAnalyzeRetry(requestBody.Count, requestBody.Retry)
		}

		c.JSON(http.StatusOK, res)
	})

	router.POST("/startStt", func(c *gin.Context) {

		//print(c.Request.Header)
		//print(c.Request.Body)

		var requestBody NeedSTT
		c.Bind(&requestBody)

		var jsondata = create_transcribe(requestBody.Index, requestBody.IsOriginal)

		c.JSON(http.StatusOK, jsondata)
	})

	router.POST("/getStt", func(c *gin.Context) {

		//print(c.Request.Header)
		//print(c.Request.Body)

		var requestBody NeedSTT
		c.Bind(&requestBody)

		var jsondata = get_transcribe(requestBody.Index, requestBody.IsOriginal)

		c.JSON(http.StatusOK, jsondata)
	})

	router.POST("/cleanUpSTT", func(c *gin.Context) {

		//print(c.Request.Header)
		//print(c.Request.Body)

		var requestBody NeedSTT
		c.Bind(&requestBody)

		var jsondata = delete_transcribe(requestBody.Index, requestBody.IsOriginal)

		c.JSON(http.StatusOK, jsondata)
	})

	router.POST("/test_cleanup", func(c *gin.Context) {

		//print(c.Request.Header)
		//print(c.Request.Body)

		var requestBody NeedSTT
		c.Bind(&requestBody)

		test_delete_all(requestBody.Index,requestBody.IsOriginal)

		c.Done()
	})

	router.POST("/uploadExcel", func(c *gin.Context) {

		var requestBody UploadExcel
		c.Bind(&requestBody)

		upload_excel(requestBody.FileName, requestBody.FileData)

		c.JSON(http.StatusOK, "Excel upload ok")
	})

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	setRouter(router)
	_ = router.Run(":8080")
}
