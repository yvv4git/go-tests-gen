package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yvv4git/go-tests-gen/internal/ports"
	"github.com/yvv4git/go-tests-gen/internal/usecases/tools"
)

type FileUpdateTool struct {
	entity *tools.FSTools
}

func NewFileUpdateTool(entity *tools.FSTools) *FileUpdateTool {
	return &FileUpdateTool{
		entity: entity,
	}
}

func (f *FileUpdateTool) Name() string {
	return "file_update"
}

func (f *FileUpdateTool) Description() string {
	return `The file_update utility that creates or replaces the contents of a file.

Input format (JSON):
{
    "path": "<absolute_path_to_file>",
    "content": "<file_content>"
}

Returns a success message.`
}

func (f *FileUpdateTool) Call(ctx context.Context, input string) (string, error) {
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

	response := ResponseMessage{
		Msg: "success",
	}

	encodedResponse, err := json.Marshal(&response)
	if err != nil {
		return "", err
	}

	return string(encodedResponse), nil
}
