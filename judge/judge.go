// judge/judge.go
package judge

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
)

// TestCase 구조체 정의
type Judge struct {
	timeout time.Duration
	workers int
}

func NewJudge(timeout time.Duration, workers int) *Judge {
	return &Judge{
		timeout: timeout,
		workers: workers,
	}
}

func (j *Judge) RunTestsParallel(code string, testCases []TestCase) []TestResult {
	results := make([]TestResult, len(testCases))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, j.workers)

	for i, tc := range testCases {
		wg.Add(1)
		go func(i int, tc TestCase) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := j.runSingleTest(code, tc)
			results[i] = result
		}(i, tc)
	}

	wg.Wait()
	return results
}

func (j *Judge) runSingleTest(code string, tc TestCase) TestResult {
	result := TestResult{
		TestCaseID: tc.ID,
		Passed:     false,
	}

	vm := goja.New()
	done := make(chan bool)
	outputs := make([]string, 0)
	currentInputIndex := 0

	go func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				result.Error = fmt.Errorf("런타임 에러: %v", r)
			}
		}()

		// console.log 함수 정의
		console := map[string]interface{}{
			"log": func(call goja.FunctionCall) goja.Value {
				var output strings.Builder
				for i, arg := range call.Arguments {
					if i > 0 {
						output.WriteString(" ")
					}
					output.WriteString(fmt.Sprint(arg))
				}
				outputs = append(outputs, output.String())
				return goja.Undefined()
			},
		}
		vm.Set("console", console)

		// prompt 함수 정의
		vm.Set("prompt", func(call goja.FunctionCall) goja.Value {
			if currentInputIndex >= len(tc.Input) {
				panic("입력 초과")
			}
			input := tc.Input[currentInputIndex]
			currentInputIndex++
			return vm.ToValue(input)
		})

		// 코드 실행
		_, err := vm.RunString(code)
		if err != nil {
			result.Error = fmt.Errorf("실행 오류: %v", err)
			return
		}

		// 출력 개수 확인
		if len(outputs) != len(tc.Output) {
			result.Message = fmt.Sprintf("출력 개수 불일치\n예상: %d개\n실제: %d개",
				len(tc.Output), len(outputs))
			return
		}

		// 출력 비교
		for i, expectedOutput := range tc.Output {
			expected := strings.TrimSpace(expectedOutput)
			actual := strings.TrimSpace(outputs[i])
			if actual != expected {
				result.Message = fmt.Sprintf("출력 불일치 (출력 #%d)\n예상: %s\n실제: %s",
					i+1, expected, actual)
				return
			}
		}

		result.Passed = true
		result.Message = "테스트 통과"
	}()

	select {
	case <-done:
		return result
	case <-time.After(j.timeout):
		result.Message = "시간 초과"
		return result
	}
}
