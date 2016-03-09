package main

import(
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/renders/multitemplate"

	"github.com/pachyderm/sandbox/src/handler"
)

var assetHandler = handler.NewAssetHandler()
// var pageHandler = handler.NewPageHandler()

func main() {
	router := gin.Default()


	assets := router.Group("/assets")
	{
		assets.GET("/styles.css", assetHandler.Serve)
		assets.GET("/main.js", assetHandler.Serve)
	}

	templates := multitemplate.New()
	templates.AddFromFiles(
		"main",
		"views/base.html",
		"views/main.html",
		"views/data.html",
	)

	router.HTMLRender = templates

	router.GET("/", func (c *gin.Context) {
		c.HTML(http.StatusOK, "main", gin.H{
			"title" : "thing",
		})
	})

	router.Run(":5678")
}
