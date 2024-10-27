// handlers/solved_handler.go

package handlers

import (
	"net/http"
	"strings"
	"time"

	"Flow-Chart-Block-Coding-Backend/judge"  // 프로젝트 경로에 맞게 수정
	"Flow-Chart-Block-Coding-Backend/models" // 프로젝트 경로에 맞게 수정

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CodeSubmission struct {
	Code      string `json:"code"`
	Username  string `json:"username"`
	ProblemID uint   `json:"problemId"`
}

func SolvedHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var submission CodeSubmission
		if err := c.ShouldBindJSON(&submission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "잘못된 요청 형식입니다",
			})
			return
		}

		// 사용자 확인
		var user models.User
		if err := db.Where("name = ?", submission.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "존재하지 않는 사용자입니다",
			})
			return
		}

		// 문제 확인
		var problem models.Problem
		if err := db.First(&problem, submission.ProblemID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "존재하지 않는 문제입니다",
			})
			return
		}

		// '/' 구분자로 테스트케이스 분리
		testCaseGroups := strings.Split(strings.TrimSpace(problem.TestcaseInput), "/")
		outputs := strings.Split(strings.TrimSpace(problem.TestcaseOutput), "/")

		// 테스트케이스 준비
		var testCases []judge.TestCase
		for i := 0; i < len(testCaseGroups); i++ {
			// 각 테스트케이스의 입력값들을 공백으로 분리
			inputs := strings.Fields(testCaseGroups[i])
			if len(inputs) < 2 {
				continue
			}

			testCase := judge.TestCase{
				ID:     i,
				Input:  []string{inputs[0], inputs[1]},
				Output: []string{strings.TrimSpace(outputs[i])},
			}
			testCases = append(testCases, testCase)
		}

		// judge 생성 및 테스트 실행
		judgeSystem := judge.NewJudge(2*time.Second, 4)
		results := judgeSystem.RunTestsParallel(submission.Code, testCases)

		// 결과 확인
		allPassed := true
		var failedMessage string
		for _, result := range results {
			if !result.Passed {
				allPassed = false
				if result.Error != nil {
					failedMessage = result.Error.Error()
				} else {
					failedMessage = result.Message
				}
				break
			}
		}

		if !allPassed {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": failedMessage,
			})
			return
		}

		// 이미 해결했는지 확인
		// var existingSolved models.Solved
		// err := db.Where("problem_id = ? AND user_id = ?",
		// 	submission.ProblemID, user.ID).First(&existingSolved).Error
		// if err == nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"success": true,
		// 		"message": "이미 해결한 문제입니다",
		// 	})
		// 	return
		// }

		// solved 테이블에 추가
		solved := models.Solved{
			ProblemID: submission.ProblemID,
			UserID:    user.ID,
			UserName:  user.Name,
			SolvedAt:  time.Now(),
		}

		if err := db.Create(&solved).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "채점 결과 저장에 실패했습니다",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "문제를 성공적으로 해결했습니다",
		})
	}
}
