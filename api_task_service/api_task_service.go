package api_task_service

import "github.com/gin-gonic/gin"

// @Summary 删除任务
// @Produce  json
// @tags 任务管理
// @Param id path int true "文章ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} string "请求错误"
// @Failure 500 {object} string "内部错误"
// @Router /api/v1/task/{id} [delete]
func Delete(c *gin.Context) {
	return
}
