package client

import "github.com/gin-gonic/gin"

func RunClient() {
	r := gin.Default()
	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
