package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yvv4git/go-tests-gen/internal/ports"
	"github.com/yvv4git/go-tests-gen/internal/usecases/tools"
)

type FileCat struct {
	entity *tools.FSTools
}

func NewFileCat(entity *tools.FSTools) *FileCat {
	return &FileCat{
		entity: entity,
	}
}

func (f FileCat) Name() string {
	return "file_cat"
}

func (f FileCat) Description() string {
	return `The file_cat utility reads the file using the specified file path. Returns the entire contents of the file.`
}

func (f FileCat) Call(ctx context.Context, input string) (string, error) {
	var params struct {
		Path string `json:"path"`
	}

	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return input, err
	}

	result, err := f.entity.FileCat(ctx, &ports.ParamsFileCat{FilePath: params.Path})
	if err != nil {
		return "", fmt.Errorf("call tool file cat [%s]: %w", params.Path, err)
	}

	encodedResult, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("marshal result: %w", err)
	}

	return string(encodedResult), nil
}
