package pipeline

import(
	"encoding/json"
	"fmt"
	e "errors"

	"github.com/gin-gonic/contrib/sessions"	
)

func LoadPipelinesFromSession(session sessions.Session) ([]string, error) {
	value := session.Get("pipelines")

	if value == nil {
		return nil, e.New("Couldnt find any pipelines in session")
	}

	fmt.Printf("raw pipeline data: %v\n", string(value.([]byte)) )

	var pipelines []string

	err := json.Unmarshal(value.([]byte), &pipelines)

	if err != nil {
		return nil, err
	}

	return pipelines, nil
}

