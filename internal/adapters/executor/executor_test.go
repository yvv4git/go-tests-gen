package executor

import (
	"context"
	"path/filepath"
	"testing"
)

func TestExecutor_RunGoUnitTest(t *testing.T) {
	// Путь к тестовому файлу
	sourceFilePath := filepath.Join("examples", "sum.go")
	testFuncName := "TestSum"

	// Создаём executor и запускаем тест
	exec := NewExecutor()
	err := exec.runGoUnitTest(context.Background(), sourceFilePath, testFuncName)

	if err != nil {
		t.Fatalf("RunGoUnitTest failed: %v", err)
	}
}
