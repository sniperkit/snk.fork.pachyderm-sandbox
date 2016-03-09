package main

import(
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/pachyderm/sandbox/src/handler"
)

var assetHandler = handler.NewAssetHandler()
// var pageHandler = handler.NewPageHandler()

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("views/*")

	assets := router.Group("/assets")
	{
		assets.GET("/styles.css", assetHandler.Serve)
		assets.GET("/main.js", assetHandler.Serve)
	}

	router.GET("/", func (c *gin.Context) {
//		pageHandler.Serve("main", c)
		c.HTML(http.StatusOK, "main.tmpl", gin.H{
			"title" : "thing",
		})
	})

	router.Run(":5678")
}
