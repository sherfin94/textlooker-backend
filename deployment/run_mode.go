package deployment

import "github.com/gin-gonic/gin"

type RunMode uint8

const Production, Development, Test = 1, 2, 3

func IsTest(context *gin.Context) bool {
	runMode, _ := context.Get("run_mode")
	return runMode.(RunMode) == Test
}
