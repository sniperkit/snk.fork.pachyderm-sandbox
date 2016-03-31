package routes

import(
	"fmt"
	"encoding/json"
	e "errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
	pfs_client "github.com/pachyderm/pachyderm/src/client/pfs"
	"github.com/segmentio/analytics-go"

	"github.com/pachyderm/sandbox/src/model/example"	
	"github.com/pachyderm/sandbox/src/session"	
)

func step1(c *gin.Context) (ex *example.Example, errors []error){
	s := sessions.Default(c)
	session.TagUserSession(s)

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

	session.Reset(s)
	s.Set("example_name", ex.Name)
	s.Set("repo_name", ex.Repo.Name)
	s.Save()

	return ex, errors
}

func step1submit(c *gin.Context) (ex *example.Example, errors []error) {
	s = sessions.Default(c)
	ex, err := example.LoadFromCookie(s, APIClient, assetHandler)

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)
	}

	ex.StepNumber, err = strconv.Atoi(c.PostForm("current_step"))

	if err != nil {
		errors = append(errors, err)
	}

	code := c.PostForm("code")
	if len(code) == 0 {
		errors = append(errors, e.New("No code submitted") )
		return nil, errors
	}

	ex.Code = code

	user, err := session.getUserToken(s)
	userPresent := (err != nil)

	if userPresent {
		analyticsClient.Track(&analytics.Track{
			Event:  "Submitted Code",
			UserId: user,
			Properties: map[string]interface{}{
				"code": code,
			},
		})
	}

	pipelines, err := ex.KickoffPipeline(code)

	if userPresent {
		analyticsClient.Track(&analytics.Track{
			Event:  "Kicked off pipelines",
			UserId: user,
			Properties: map[string]interface{}{
				"pipelines": pipelines,
			},
		})
	}

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)		
	}

	pipelineJSON, _ := json.Marshal(pipelines)

	fmt.Printf("Kicked off pipelines %v\n", string(pipelineJSON))

	s.Set("pipelines", pipelineJSON)
	s.Save()

	return ex, errors
}
