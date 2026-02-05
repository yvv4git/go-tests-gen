package ports

import "context"

type (
	ParamsRunGoUnitTest struct {
		SourceFilePath string
		TestFnName     string
	}

	ResultRunGoUnitTest struct {
	}
)

type Executor interface {
	RunGoUnitTest(ctx context.Context, params *ParamsRunGoUnitTest) (*ResultRunGoUnitTest, error)
}
