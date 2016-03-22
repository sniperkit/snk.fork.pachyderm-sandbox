package main

import(
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/pachyderm/pachyderm/src/client"
	pfs_client "github.com/pachyderm/pachyderm/src/client/pfs"

	"github.com/pachyderm/sandbox/src/asset"
	"github.com/pachyderm/sandbox/src/example"
)

var assetHandler = asset.NewAssetHandler()
var router = gin.New()
var APIClient *client.APIClient

func main() {
	apiClient, err := client.New()
	APIClient = apiClient
	// SJ: This feels wrong, am I missing a go-ism to solve the 'declared' compile error?

	if err != nil {
		fmt.Printf("Error connecting to pachd %v\n", err)
	}

	assets := router.Group("/assets")
	{
		assets.GET("/styles.css", assetHandler.Serve)
		assets.GET("/main.js", assetHandler.Serve)
	}

	router.HTMLRender = loadTemplates()

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
			"repos" : repos,
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
