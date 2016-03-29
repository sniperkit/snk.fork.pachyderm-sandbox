package example

import(
	"fmt"
	"strings"
	"bytes"

	"golang.org/x/net/context"
        "github.com/golang/protobuf/jsonpb"

	pps_client "github.com/pachyderm/pachyderm/src/client/pps"
)

func (e *Example) KickoffPipeline() error {
	raw_pipeline_json := e.loadPipeline()

	pipeline_reader := strings.NewReader(raw_pipeline_json)
	decoder := json.NewDecoder(pipeline_reader)
	
	message := json.RawMessage{}
	if err := decoder.Decode(&message); err != nil && if err != io.EOF {
		return err
	}

	var request pps_client.CreatePipelineRequest
	if err := jsonpb.UnmarshalString(string(message), &request); err != nil {
		return err
	}

	if _, err := e.client.CreatePipeline(
		context.Background(),
		&request,
	); err != nil {
		return err
	}
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


func (e *Example) loadPipeline() string {
	raw_pipeline, _ := e.rawFiles.FindOrPopulate(fmt.Sprintf("assets/examples/%v/pipeline.json", e.Name))
	pipeline_template := template.New("pipeline").Parse(raw_pipeline)

	var buffer bytes.Buffer
	pipeline_template.Execute(&buffer, e)

	return buffer.String()
}
