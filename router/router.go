package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wisdom-oss/common-go/v3/middleware/gin/jwt"

	app "microservice/internal"
	internal "microservice/internal/router"
	v1Routes "microservice/routes/v1"
)

// Configure generates a new router and adds routes to the router
//
// The router can also be imported during tests, as long as the tests are in a
// separate package.
// If the tests are in the same package (e.g. routes defined in `v3` and tests
// also defined in `v3`) an import cycle exists.
func Configure() (*gin.Engine, error) {
	r, err := internal.GenerateRouter()
	if err != nil {
		return nil, err
	}

	authorize := jwt.ScopeRequirer{}
	authorize.Configure(app.ServiceName)

	v1 := r.Group("/v1")
	{
		v1.GET("/public/*filename", v1Routes.GetFile)
		v1.GET("/:bucket/*filename", authorize.RequireRead, v1Routes.GetFile)
		v1.PUT("/:bucket/*basepath", authorize.RequireWrite, v1Routes.Upload)
		v1.DELETE("/:bucket/*filename", authorize.RequireDelete, v1Routes.DeleteFile)
	}

	return r, nil
}
