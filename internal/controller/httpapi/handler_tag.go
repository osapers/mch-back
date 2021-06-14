package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type searchTags struct {
	Query string `json:"query"`
}

func (s *Server) searchTags(c *gin.Context) {
	var req searchTags

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	limit := 20

	t, err := s.ts.Search(c.Request.Context(), req.Query, limit)
	c.JSON(http.StatusOK, jsonResp(t, err))
}
