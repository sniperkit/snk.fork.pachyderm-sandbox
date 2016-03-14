package main

import(
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/renders/multitemplate"

	"github.com/pachyderm/sandbox/src/asset"
	"github.com/pachyderm/pachyderm"
	"github.com/pachyderm/pachyderm/src/pfs/pfsutil"
)

var assetHandler = asset.NewAssetHandler()
var router = gin.New()
var APIClient *pachyderm.APIClient

func main() {
	client, err := pachyderm.NewAPIClient()
	APIClient = client
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

		repos, err := pfsutil.ListRepo(APIClient)

		if err != nil {
			fmt.Printf("ERR! %v\n", err)
			errors = append(errors, err)
		}

		c.HTML(http.StatusOK, page, gin.H{
			"title" : "REPL",
			"repos" : repos,
			"errors": errors,
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
