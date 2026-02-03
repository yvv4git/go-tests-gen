package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yvv4git/go-tests-gen/internal/ports"
	"github.com/yvv4git/go-tests-gen/internal/usecases/tools"
)

type FileUpdate struct {
	entity *tools.FSTools
}

func NewFileUpdate(entity *tools.FSTools) *FileUpdate {
	return &FileUpdate{
		entity: entity,
	}
}

func (f *FileUpdate) Name() string {
	return "file_update"
}

func (f *FileUpdate) Description() string {
	return `The file_update utility that replaces the contents of the specified file with the contents specified in the arguments.`
}

func (f *FileUpdate) Call(ctx context.Context, input string) (string, error) {
	var params struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return input, err
	}

	_, err := f.entity.FileUpdate(ctx, &ports.ParamsFileUpdate{
		FilePath: params.Path,
		Content:  params.Content,
	})
	if err != nil {
		return "", fmt.Errorf("update file [%s]: %w", params.Path, err)
	}

	var response struct {
		Msg string `json:"msg"`
	}

	response.Msg = "success"

	encodedResponse, err := json.Marshal(&response)
	if err != nil {
		return "", err
	}

	return string(encodedResponse), nil
}
