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

	router.GET("/", handle("main"))

	router.Run(":9080")
	
}

func handle(page string) ( func (*gin.Context) ){
	return func(c *gin.Context) {
		if gin.Mode() == "debug" {
			router.HTMLRender = loadTemplates()
		}

		var errors []error

		example, err := example.New("fruit-stand", APIClient, assetHandler)

		if err != nil {
			fmt.Printf("ERR! %v\n", err)
			errors = append(errors, err)
		} else {
			// Silly ... but go compiler doesn't know I'm using it in a view
			fmt.Printf("Loaded %v\n", example.Name)			
		}

		repos, err := pfs_client.ListRepo(APIClient)

		if err != nil {
			fmt.Printf("ERR! %v\n", err)
			errors = append(errors, err)
		}

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
