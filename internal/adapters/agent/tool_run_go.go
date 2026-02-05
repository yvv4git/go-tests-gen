package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yvv4git/go-tests-gen/internal/ports"
	"github.com/yvv4git/go-tests-gen/internal/usecases/tools"
)

type RunGoTestTool struct {
	entity *tools.ExecTools
}

func NewRunGoTestTool(entity *tools.ExecTools) *RunGoTestTool {
	return &RunGoTestTool{
		entity: entity,
	}
}

func (r *RunGoTestTool) Name() string {
	return "run_go_unit_test"
}

func (r *RunGoTestTool) Description() string {
	return `The run_go_unit_test utility that runs a unit test for a specific function in a Go file.

Input format (JSON):
{
    "path": "<absolute_path_to_source_file>",
    "content": "<function_name_to_test>"
}

Returns a success message.`
}

func (r *RunGoTestTool) Call(ctx context.Context, input string) (string, error) {
	var params struct {
		SourceFilePath string `json:"path"`
		TestFnName     string `json:"content"`
	}

	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return input, err
	}

	err := r.entity.RunGoUnitTest(ctx, &ports.ParamsRunGoUnitTest{
		SourceFilePath: params.SourceFilePath,
		TestFnName:     params.TestFnName,
	})
	if err != nil {
		return "", fmt.Errorf("run go unit test[%s]: %w", params.TestFnName, err)
	}

	response := ResponseMessage{
		Msg: "success",
	}

	encodedResponse, err := json.Marshal(&response)
	if err != nil {
		return "", err
	}

	return string(encodedResponse), nil
}
