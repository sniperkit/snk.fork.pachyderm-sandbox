package main

import(
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/pachyderm/sandbox/src/handler"
)

var assetHandler = handler.NewAssetHandler()
// var pageHandler = handler.NewPageHandler()

func main() {
	router := gin.Default()


	assets := router.Group("/assets")
	{
		fmt.Printf("I SEE AN ASSET REQUEST")
		assets.GET("/styles.css", assetHandler.Serve)
		assets.GET("/main.js", assetHandler.Serve)
	}


	router.LoadHTMLGlob("views/*")

	router.GET("/", func (c *gin.Context) {
//		pageHandler.Serve("main", c)
		c.HTML(http.StatusOK, "main.tmpl", gin.H{
			"title" : "thing",
		})
	})

	router.Run(":5678")
}
