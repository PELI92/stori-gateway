package proxy

import "github.com/gin-gonic/gin"

type Proxy interface {
	Handle(c *gin.Context)
}
