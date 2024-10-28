// handlers/solved_handler.go

package handlers

import (
	"net/http"
	"strconv"
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
		var existingSolved models.Solved
		err := db.Where("problem_id = ? AND user_id = ?",
			submission.ProblemID, user.ID).First(&existingSolved).Error
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "이미 해결한 문제입니다",
			})
			return
		}

		// solved 테이블에 추가
		solved := models.Solved{
			ProblemID: submission.ProblemID,
			UserID:    user.ID,
			UserName:  user.Name,
			// SolvedAt:  time.Now(),
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

// handlers/solved_handler.go에 다음 두 함수를 추가합니다.

// GetUserSolvedProblems는 사용자가 해결한 문제 목록을 반환합니다.
func GetUserSolvedProblems(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Param("username") // URL에서 username 파라미터를 가져옴

		var solved []models.Solved
		if err := db.Where("user_name = ?", userName).Find(&solved).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "해결한 문제 목록을 조회하는데 실패했습니다",
			})
			return
		}

		// 문제 상세 정보를 포함하여 응답
		var problems []map[string]interface{}
		for _, s := range solved {
			var problem models.Problem
			if err := db.First(&problem, s.ProblemID).Error; err != nil {
				continue // 문제를 찾을 수 없는 경우 스킵
			}

			problems = append(problems, map[string]interface{}{
				"problemId": problem.ID,
				"title":     problem.Title,
				"solvedAt":  s.SolvedAt,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"userName":       userName,
				"solvedProblems": problems,
			},
		})
	}
}

func GetProblemSolvedUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// string을 uint로 변환
		problemID, err := strconv.ParseUint(c.Param("problem_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "잘못된 문제 ID 형식입니다",
			})
			return
		}

		// 문제가 존재하는지 확인
		var problem models.Problem
		if err := db.First(&problem, uint(problemID)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "존재하지 않는 문제입니다",
			})
			return
		}

		var solved []models.Solved
		if err := db.Where("problem_id = ?", uint(problemID)).Find(&solved).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "해결한 사용자 목록을 조회하는데 실패했습니다",
			})
			return
		}

		// 사용자 상세 정보를 포함하여 응답
		var users []map[string]interface{}
		for _, s := range solved {
			var user models.User
			if err := db.First(&user, s.UserID).Error; err != nil {
				continue // 사용자를 찾을 수 없는 경우 스킵
			}

			users = append(users, map[string]interface{}{
				"userId":   user.ID,
				"userName": user.Name,
				"classnum": user.Classnum,
				"solvedAt": s.SolvedAt,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"problemId":    uint(problemID),
				"problemTitle": problem.Title,
				"solvedUsers":  users,
			},
		})
	}
}

// GetProblemSolvedUsers는 특정 문제를 해결한 사용자 목록을 반환합니다.
// func GetProblemSolvedUsers(db *gorm.DB) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         // string을 uint로 변환
//         problemID, err := strconv.ParseUint(c.Param("problemId"), 10, 32)
//         if err != nil {
//             c.JSON(http.StatusBadRequest, gin.H{
//                 "success": false,
//                 "message": "잘못된 문제 ID 형식입니다",
//             })
//             return
//         }

//         // 문제가 존재하는지 확인
//         var problem models.Problem
//         if err := db.First(&problem, uint(problemID)).Error; err != nil {
//             c.JSON(http.StatusNotFound, gin.H{
//                 "success": false,
//                 "message": "존재하지 않는 문제입니다",
//             })
//             return
//         }

//         var solved []models.Solved
//         if err := db.Where("problem_id = ?", uint(problemID)).Find(&solved).Error; err != nil {
//             c.JSON(http.StatusInternalServerError, gin.H{
//                 "success": false,
//                 "message": "해결한 사용자 목록을 조회하는데 실패했습니다",
//             })
//             return
//         }

//         // 사용자 상세 정보를 포함하여 응답
//         var users []map[string]interface{}
//         for _, s := range solved {
//             var user models.User
//             if err := db.First(&user, s.UserID).Error; err != nil {
//                 continue // 사용자를 찾을 수 없는 경우 스킵
//             }

//             users = append(users, map[string]interface{}{
//                 "userId":   user.ID,
//                 "userName": user.Name,
//                 "classnum": user.Classnum,
//                 "solvedAt": s.SolvedAt,
//             })
//         }

//         c.JSON(http.StatusOK, gin.H{
//             "success": true,
//             "data": gin.H{
//                 "problemId":    uint(problemID),
//                 "problemTitle": problem.Title,
//                 "solvedUsers":  users,
//             },
//         })
//     }
// }
