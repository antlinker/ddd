package ginhook

import (
	"github.com/antlinker/ddd"
	"github.com/antlinker/ddd/log"
	"github.com/gin-gonic/gin"
)

var (
	dddctxkey = "__ddd.ctx.key__"
)

// DDDContext ddd context 接入gin框架
func DDDContext(getuid func(c *gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceid := c.GetHeader("X-Request-Id")
		if getuid != nil {
			uid := getuid(c)

			// ctx := ddd.NewContext(c, uid, log.FromContext(c))
			ctx := ddd.NewTraceContext(c, traceid, uid, log.FromContext(c))
			c.Set(dddctxkey, ctx)
		}
	}
}

// FromContext  从gin.Context中获取 ddd.Context
func FromContext(c *gin.Context) ddd.Context {
	if ctmp, ok := c.Get(dddctxkey); ok {
		return ctmp.(ddd.Context)
	}
	return nil
}
