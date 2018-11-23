package mgorepo

import (
	"github.com/antlinker/ddd"
	"github.com/antlinker/ddd/ginhook"
	"github.com/gin-gonic/gin"
)

// MgoContext ddd context 接入gin框架
func MgoContext(d ddd.Domain) gin.HandlerFunc {
	return func(c *gin.Context) {
		if ctx := ginhook.FromContext(c); ctx != nil {
			SetContextDB(ctx, d)
			c.Next()
			ReleaseDB(ctx)
			return
		}
		panic("需要使用插件ginhook.DDDContext后，在使用该插件")

	}
}
