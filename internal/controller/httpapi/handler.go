package httpapi

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) addHandlers(r *gin.Engine) {
	apiv1 := r.Group("/api/v1")
	{
		events := apiv1.Group("/events", s.tokenMiddleware())
		{
			events.GET("", s.listEventsHandler)
			events.POST("/:eventID/participate", s.participateInEvent)
		}

		auth := apiv1.Group("/auth")
		{
			auth.POST("", s.authorizeUser)
		}

		user := apiv1.Group("/users", s.tokenMiddleware())
		{
			user.PUT("/me", s.updateUser)
			user.GET("/me", s.getUser)
			projects := user.Group("/projects")
			{
				projects.GET("/my", s.getMyProjects)
				projects.POST("/generate", s.generateProject)
				projects.GET("/search", s.searchProject)
				projects.POST("/:id/view", s.userViewProject)
				projects.POST("/:id/apply", s.userApplyProject)
				projects.DELETE("/:id", s.deleteMyProject)
			}
		}

		projects := apiv1.Group("/projects")
		{
			projects.GET("/:id/candidates", s.getProjectCandidates)
		}

		tags := apiv1.Group("/tags", s.tokenMiddleware())
		{
			tags.POST("/search", s.searchTags)
		}
	}

	if !s.cfg.envName.IsProd() {
		testdata := r.Group("/testdata")
		{
			testdata.GET("/events", s.loadTestEvents)
		}
	}
}
