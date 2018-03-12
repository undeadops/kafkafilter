package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/undeadops/kafkafilter/v1/topics"
)

func main() {
	r := gin.Default()

	fmt.Println("Hello World")
	v1 := r.Group("/api/v1")
	topics.TopicsRegister(v1.Group("/topic"))

	r.Run(":5000") // listen and serve on 0.0.0.0:5000
}
