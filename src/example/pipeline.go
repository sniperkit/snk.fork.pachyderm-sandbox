package example

import(
	"fmt"
	"strings"
	"bytes"
	"encoding/json"
	"io"
	"text/template"

	"golang.org/x/net/context"
        "github.com/golang/protobuf/jsonpb"	

	pps_client "github.com/pachyderm/pachyderm/src/client/pps"
)

func (e *Example) KickoffPipeline() ([]string, error) {
	raw_pipeline_json, err := e.loadPipeline()

	if err != nil {
		return nil, err
	}

	pipeline_reader := strings.NewReader(raw_pipeline_json)
	decoder := json.NewDecoder(pipeline_reader)

	var pipelineNames []string

	for {
	
		message := json.RawMessage{}

		if err = decoder.Decode(&message); err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("err decoding pipeline json %v\n", err)
				return nil, err
			}
		}

		fmt.Printf("message %v\n", string(message))
		var request pps_client.CreatePipelineRequest

		if err := jsonpb.UnmarshalString(string(message), &request); err != nil {
			fmt.Printf("err unmarshaling json %v\n", err)
			return nil, err
		}

		fmt.Printf("create pipeline request: %v\n", request)

		pipelineNames = append(pipelineNames, request.Pipeline.Name)

		if _, err := e.client.CreatePipeline(
			context.Background(),
			&request,
		); err != nil {
			return nil, err
		}
	}
	

	return pipelineNames, nil
}

func (e *Example) IsPipelineDone() {
	// Will have to poll for this
	// Check the right job state? and see if it completes?

	//e.destroyPipeline()
}

// Will be called once pipeline is done ...
// to keep things idempotent in a REPL
func (e *Example) destroyPipeline() {

}


func (e *Example) loadPipeline() (string, error) {
	raw_pipeline, err := e.rawFiles.FindOrPopulate(fmt.Sprintf("assets/examples/%v/pipeline.json", e.Name))

	if err != nil {
		return "", err
	}

	pipeline_template, err := template.New("pipeline").Parse(string(raw_pipeline))

	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	pipeline_template.Execute(&buffer, e)

	return buffer.String(), nil
}
