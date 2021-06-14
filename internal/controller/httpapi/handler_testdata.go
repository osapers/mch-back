package httpapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) loadTestEvents(c *gin.Context) {
	count, _ := strconv.Atoi(c.DefaultQuery("count", "5"))
	err := s.es.LoadTestEvents(count)
	c.JSON(http.StatusOK, jsonResp(nil, err))
}
