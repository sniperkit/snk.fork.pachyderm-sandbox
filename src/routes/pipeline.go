package routes

import(
	"encoding/json"
	"fmt"
	e "errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/pachyderm/sandbox/src/example"	
)

func check_pipeline_status(c *gin.Context) (ex *example.Example, errors []error){
	s := sessions.Default(c)

	ex, err := example.LoadFromCookie(s, APIClient, assetHandler)

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)
	}

	fmt.Printf("Loaded example: %v\n", ex)
	
	value := s.Get("pipelines")

	if value == nil {
		errors = append(errors, e.New("Couldnt find any pipelines in session"))
		return nil, errors
	}

	var pipelines []string

	err = json.Unmarshal(value.([]byte), &pipelines)

	if err != nil {
		errors = append(errors, err)
		return nil, errors
	}

	fmt.Printf("Pipelines: %v\n", pipelines)

	return ex, errors
}
