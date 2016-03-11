package main

import(
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/pachyderm/sandbox/src/asset"
)

var assetHandler = asset.NewAssetHandler()
var router = gin.New()

func init() {

	assets := router.Group("/assets")
	{
		assets.GET("/styles.css", assetHandler.Serve)
		assets.GET("/main.js", assetHandler.Serve)
	}

	router.HTMLRender = loadTemplates()

	router.GET("/", handle("main"))

	http.Handle("/", router)

	appengine.Main()

}

func handle(page string) ( func (*gin.Context) ){
	return func(c *gin.Context) {
		if gin.Mode() == "debug" {
			router.HTMLRender = loadTemplates()
		}
		
		c.HTML(http.StatusOK, page, gin.H{
			"title" : "thing",
		})

		log.Infof(c, fmt.Sprintf("Serving the [%v] page.", page))
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
