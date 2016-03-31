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
	"github.com/pachyderm/sandbox/src/util"
)


var ErrNoPipelinesInSession = errors.New("Could not find pipeline from session data")

type PipelineState int

const (
	PipelineNotFound PipelineState = iota
	PipelineWorking
	PipelineCompleted
)

func fullyQualifyName(repo string, input string) string {
	return fmt.Sprintf("%v-pipeline-%v", repo, util.GenerateUniqueToken())
}

func (ex *Example) fullyQualifyRequest(request *pps_client.CreatePipelineRequest, fqPipelineName string, chained bool) {
	// Pipeline names must always be unique
	request.Pipeline.Name = strings.Replace(request.Pipeline.Name, ex.Repo.DisplayName, fqPipelineName, -1)

	// Replace transform/stdin
	// inputs/repo/name

	replacement := fqPipelineName
	if !chained {
		replacement = ex.Repo.Name
	}

	for i, _ := range(request.Transform.Stdin) {
		request.Transform.Stdin[i] = strings.Replace(request.Transform.Stdin[i], ex.Repo.DisplayName, replacement, -1)
	}

	for i, _ := range(request.Inputs) {
		request.Inputs[i].Repo.Name = strings.Replace(request.Inputs[i].Repo.Name, ex.Repo.DisplayName, replacement, -1)
	}

}

func (e *Example) KickoffPipeline(manifest string) ([]string, error) {
	fmt.Printf("raw manifest:\n%v\n\n", manifest)

	pipeline_reader := strings.NewReader(manifest)
	decoder := json.NewDecoder(pipeline_reader)

	var pipelineNames []string
	chainedTransform := false
	var fqPipelineName string
	
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

		// Need to add a unique token at the end so that the user can 're-run' the transform
		fmt.Printf("pipeline name? %v\n", request.Pipeline.Name)
		fmt.Printf("create pipeline request: %v\n", request)

		if !chainedTransform {
			fqPipelineName = fullyQualifyName(e.Repo.Name, request.Pipeline.Name)
		}
		e.fullyQualifyRequest(&request, fqPipelineName, chainedTransform)

		fmt.Printf("create pipeline NORMALIZED request: %v\n", request)


		pipelineNames = append(pipelineNames, request.Pipeline.Name)

		if _, err := e.client.CreatePipeline(
			context.Background(),
			&request,
		); err != nil {
			return nil, err
		}

		chainedTransform = true
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

		fmt.Printf("reading raw pipelien from session: %v\n", pipeline)
		pipeline := strings.Split(pipeline, "-pipeline")[0]
		fmt.Printf("normalized pipeline from session: %v\n", pipeline)

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
