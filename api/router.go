package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(e *Env) *gin.Engine {
	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.GET("/slugs", e.GetSlugs)
	r.POST("/slugs", e.AddSlug)
	r.GET("slug/:id", e.GetSlugById)
	r.DELETE("slug/:id", e.DeleteSlugById)

	r.GET("/internal/panic", func(c *gin.Context) {
		// panic with a string -- the custom middleware could save this to a database or report it to the user
		panic("panic emulation")
	})
	return r
}
