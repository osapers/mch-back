package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getMyProjects(c *gin.Context) {
	projects, err := s.us.MyProjects(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, jsonResp(projects, nil))
}

func (s *Server) generateProject(c *gin.Context) {
	err := s.us.GenerateProject(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, jsonResp(map[string]string{"status": "ok"}, nil))
}

func (s *Server) searchProject(c *gin.Context) {
	projects, err := s.us.SearchProjects(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, jsonResp(projects, nil))
}

func (s *Server) deleteMyProject(c *gin.Context) {
	projectID := c.Param("id")

	err := s.us.DeleteProject(c.Request.Context(), projectID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, jsonResp(map[string]string{"status": "ok"}, nil))
}

func (s *Server) getProjectCandidates(c *gin.Context) {
	projectID := c.Param("id")

	candidates, err := s.us.GetProjectCandidates(c.Request.Context(), projectID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, jsonResp(candidates, nil))
}
