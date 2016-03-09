package handler

import(
	"io/ioutil"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/realistschuckle/gohaml"
)

type AssetHandler struct {
	files map[string][]byte
}

func NewAssetHandler() *AssetHandler {
	return &AssetHandler{
		files: make(map[string][]byte),
	}
}

func (a *AssetHandler) Serve(c *gin.Context) {
	
	path := fmt.Sprintf(".%v", c.Request.URL.Path)
	content, ok := a.files[path]

	if !ok {
		fmt.Println(path)
		content, err := ioutil.ReadFile(path)

		fmt.Printf("Found content: %v\n", string(content))

		if err != nil {
			fmt.Println(err)
			c.String(http.StatusNotFound, "Asset not found")
			return
		}

		a.files[path] = content
	}

	c.String(http.StatusOK, string(content) )
}


type PageHandler struct {
	templates map[string][]byte
}

func NewPageHandler() *PageHandler {
	return &PageHandler{
		templates: make(map[string][]byte),
	}
}

func (p *PageHandler) Serve(page string, c *gin.Context) {
	path := fmt.Sprintf("views/%v.haml", page)

	content, ok := p.templates[path]

	if !ok || true { //Always read file while debugging
		fmt.Println(path)
		newContent, err := ioutil.ReadFile(path)

		fmt.Printf("Found content: %v\n", string(newContent))

		if err != nil {
			fmt.Println(err)
			c.String(http.StatusNotFound, "View not found")
			return
		}

		p.templates[path] = newContent
		content = newContent
	}

	fmt.Printf("Using tempalte: %v\n", content)

	scope := make(map[string]interface{})
	scope["context"] = c
	engine, err := gohaml.NewEngine(string(content))
	if err != nil {
		fmt.Printf("err w engine: %v\n", err)
	}

	output := engine.Render(scope)

	c.String(http.StatusOK, output)
}
