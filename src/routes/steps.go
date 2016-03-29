package routes

import(
	"fmt"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
	pfs_client "github.com/pachyderm/pachyderm/src/client/pfs"

	"github.com/pachyderm/sandbox/src/example"	
)

func step1(c *gin.Context) (ex *example.Example, errors []error){

	ex, err := example.New("fruit-stand", APIClient, assetHandler)

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)
	} else {
		// Silly ... but go compiler doesn't know I'm using it in a view
		fmt.Printf("Loaded %v\n", ex.Name)			
	}

	repos, err := pfs_client.ListRepo(APIClient)

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)
	} else {
		// Again ... silly ... but compiler doesn't know its used in a view
		fmt.Printf("Loaded %v repos", len(repos))
	}

	s := sessions.Default(c)
	s.Set("example_name", ex.Name)
	s.Set("repo_name", ex.Repo.Name)
	s.Save()

	return ex, errors
}


func step1submit(c *gin.Context) (ex *example.Example, errors []error) {
	
	ex, err := example.LoadFromCookie(sessions.Default(c), APIClient, assetHandler)

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)
	}

	pipelines, err := ex.KickoffPipeline()

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)		
	}

	s := sessions.Default(c)
	pipelineJSON, _ := json.Marshal(pipelines)

	s.Set("pipelines", pipelineJSON)
	s.Save()

	return ex, errors
}
