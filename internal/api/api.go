package api

import "github.com/gin-gonic/gin"

func StartApi() error {

	r := gin.Default()
	r.GET("/update", triggerUpdate)
	r.POST("/update", triggerUpdate)

	return r.Run()
}
