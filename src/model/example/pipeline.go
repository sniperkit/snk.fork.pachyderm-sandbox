package example

import(
	"fmt"
	"strings"
	"bytes"
	"encoding/json"
	"io"
	"text/template"
	"errors"

	"golang.org/x/net/context"
        "github.com/golang/protobuf/jsonpb"	
	"github.com/gin-gonic/contrib/sessions"	
	pps_client "github.com/pachyderm/pachyderm/src/client/pps"

	"github.com/pachyderm/sandbox/src/model/pipeline"
)


var ErrNoPipelinesInSession = errors.New("Could not find pipeline from session data")

type PipelineState int

const (
	PipelineNotFound PipelineState = iota
	PipelineWorking
	PipelineCompleted
)

func (e *Example) KickoffPipeline(manifest string) ([]string, error) {
	fmt.Printf("raw manifest:\n%v\n\n", manifest)

	// Replace all instances of 
	fmt.Printf("REPLACING %v -> %v\n", e.Repo.DisplayName, e.Repo.Name)
	manifest = strings.Replace(manifest, e.Repo.DisplayName, e.Repo.Name, -1)

	fmt.Printf("normalized manifest:\n%v\n\n", manifest)

	pipeline_reader := strings.NewReader(manifest)
	decoder := json.NewDecoder(pipeline_reader)

	var pipelineNames []string

	for {
	
		message := json.RawMessage{}

		if err := decoder.Decode(&message); err != nil {
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

func (e *Example) getJobStates(session sessions.Session) (states map[string]pps_client.JobState, err error){
	states = make(map[string]pps_client.JobState)

	pipelines, err := pipeline.LoadPipelinesFromSession(session)

	if err != nil {
		return nil, ErrNoPipelinesInSession
	}


	for _, pipeline := range(pipelines) {

		jobInfos, err := e.client.ListJob(
			context.Background(),
			&pps_client.ListJobRequest{
				Pipeline: &pps_client.Pipeline{
					Name: pipeline,
				},
			},
		)

		if err != nil {
			return nil, err
		}
		
		var state pps_client.JobState

		for _, jobInfo := range(jobInfos.JobInfo) {
			state = jobInfo.State

			if state != pps_client.JobState_JOB_STATE_SUCCESS {
				break
			}
		}

		states[pipeline] = state
	}

	return states, nil
}

func (e *Example) IsPipelineDone(session sessions.Session) (status bool, states map[string]string, err error) {

	rawStates, err := e.getJobStates(session)

	states = make(map[string]string)
	status = true

	for pipeline, state := range(rawStates) {

		thisJobDone := (state == pps_client.JobState_JOB_STATE_SUCCESS)
		status = status && thisJobDone
		states[pipeline] = pps_client.JobState_name[int32(state)]
	}

	//e.destroyPipeline()
	return status, states, nil
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
