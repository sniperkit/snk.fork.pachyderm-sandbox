package routes

import(
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/segmentio/analytics-go"

	"github.com/pachyderm/sandbox/src/session"	
	"github.com/pachyderm/sandbox/src/model/example"	
	"github.com/pachyderm/sandbox/src/model/pipeline"
	"github.com/pachyderm/sandbox/src/model/repo"
)


func check_pipeline_status(c *gin.Context) {
	errors := make([]error,0)

	if gin.Mode() == "debug" {
		router.HTMLRender = loadTemplates()
	}

	s := sessions.Default(c)

	ex, err := example.LoadFromCookie(s, APIClient, assetHandler)

	if err != nil {
		fmt.Printf("ERR! %v\n", err)
		errors = append(errors, err)
	}

	status, states, err := ex.IsPipelineDone(s)

	if err != nil {
		errors = append(errors, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"errors": errors,
		"status": status,
		"states": states,
	})

	if status {

		user, err := session.GetUserToken(s)
		userPresent := (err != nil)

		if userPresent {
			analyticsClient.Track(&analytics.Track{
				Event:  "Transform Completed",
				UserId: user,
				Properties: map[string]interface{}{
					"status": status,
					"states": states,
				},
				Context: map[string]interface{}{
					"integrations" : map[string]interface{}{
						"All": true,
					},
				},
			})
		}
	}

}

func list_output_repos(c *gin.Context) {
	errors := make([]error,0)

	if gin.Mode() == "debug" {
		router.HTMLRender = loadTemplates()
	}

	s := sessions.Default(c)

	pipelines, err := pipeline.LoadPipelinesFromSession(s)

	if  err != nil {
		errors = append(errors, err)
	}

	var repos []*repo.SandboxRepo

	for _, pipeline := range(pipelines) {

		// Pipeline name == output repo name
		r, err := repo.Load(APIClient, pipeline)

		if err != nil {
			errors = append(errors, err)
		}

		repos = append(repos, r)
	}

	if len(errors) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"errors": errors,
		})
		
	} else {
		fmt.Printf("got repos? %v\n", repos)
		c.HTML(http.StatusOK, "pipeline_output", gin.H{
			"repos": repos,
		})
	}

	user, err := session.GetUserToken(s)
	userPresent := (err != nil)

	if userPresent {
		analyticsClient.Track(&analytics.Track{
			Event:  "Loaded output repo",
			UserId: user,
			Properties: map[string]interface{}{
				"repos": repos,
			},
			Context: map[string]interface{}{
				"integrations" : map[string]interface{}{
					"All": true,
				},
			},
		})
	}
	

}
