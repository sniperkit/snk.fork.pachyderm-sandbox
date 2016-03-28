package routes

import(
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/pachyderm/pachyderm/src/client"

	"github.com/pachyderm/sandbox/src/asset"
	"github.com/pachyderm/sandbox/src/example"
)

var assetHandler = asset.NewAssetHandler()
var router = gin.New()
var APIClient *client.APIClient

func init() {
	apiClient, _ := client.New()
	APIClient = apiClient
	// SJ: This feels wrong, am I missing a go-ism to solve the 'declared' compile error?

	assets := router.Group("/assets")
	{
		assets.GET("/styles.css", assetHandler.Serve)
		assets.GET("/main.js", assetHandler.Serve)
	}

	router.HTMLRender = loadTemplates()
}

func Route() {

	router.GET("/", handle("main", step1))

	router.POST("/", handle("main", step1submit))

	router.Run(":9080")
}

func handle(page string, customHandler func(*gin.Context) (*example.Example, []error) ) ( func (*gin.Context) ){
	return func(c *gin.Context) {
		if gin.Mode() == "debug" {
			router.HTMLRender = loadTemplates()
		}

		example, errors := customHandler(c)

		c.HTML(http.StatusOK, page, gin.H{
			"title" : example.Name + "Example",
			"errors": errors,
			"example": example,
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
