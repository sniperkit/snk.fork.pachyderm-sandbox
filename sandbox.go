package main

import(
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/renders/multitemplate"

	"github.com/pachyderm/sandbox/src/asset"
)

var assetHandler = asset.NewAssetHandler()
var router = gin.New()

func main() {

	assets := router.Group("/assets")
	{
		assets.GET("/styles.css", assetHandler.Serve)
		assets.GET("/main.js", assetHandler.Serve)
	}

	router.HTMLRender = loadTemplates()

	router.GET("/", handle("main"))

	router.Run(":9080")
	
}

func handle(page string) ( func (*gin.Context) ){
	return func(c *gin.Context) {
		if gin.Mode() == "debug" {
			router.HTMLRender = loadTemplates()
		}
		
		c.HTML(http.StatusOK, page, gin.H{
			"title" : "successfully deployed? yes",
		})
	}
}

func loadTemplates() multitemplate.Render {
	templates := multitemplate.New()
	templates.AddFromFiles(
		"main",
		"views/base.html",
		"views/main.html",
		"views/data.html",
		"views/copy.html",
		"views/code.html",
	)

	return templates
}
