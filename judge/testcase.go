// judge/testcase.go
package judge

// TestCase 구조체 정의
type TestCase struct {
	ID     int
	Input  []string
	Output []string
}

// TestResult 구조체 정의
type TestResult struct {
	TestCaseID int
	Passed     bool
	Message    string
	Error      error
}
