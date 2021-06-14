package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/osapers/mch-back/internal/service/user"
	"github.com/osapers/mch-back/internal/types"
)

type authorizeUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authorizeUserRes struct {
	Token           string `json:"token"`
	IsWrongPassword bool   `json:"is_wrong_password"`
}

func (s *Server) authorizeUser(c *gin.Context) {
	var req authorizeUserReq

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	if req.Email == "" || req.Password == "" {
		err = fmt.Errorf("email and password are required")
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	u, err := s.us.Authorize(c.Request.Context(), req.Email, req.Password)
	if errors.As(err, &user.WrongPasswordError) {
		res := &authorizeUserRes{IsWrongPassword: true}
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(res, nil))
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	res := &authorizeUserRes{Token: newToken(u.ID, s.cfg.secret)}
	c.JSON(http.StatusOK, jsonResp(res, nil))
}

func (s *Server) participateInEvent(c *gin.Context) {
	eventID := c.Param("eventID")

	err := s.us.ParticipateInEvent(c.Request.Context(), eventID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getUser(c *gin.Context) {
	u, err := s.us.GetMe(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, jsonResp(u.ToJson(), nil))
}

func (s *Server) updateUser(c *gin.Context) {
	var req types.User

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	u, err := s.us.Update(c.Request.Context(), &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, jsonResp(u.ToJson(), nil))
}

func (s *Server) userViewProject(c *gin.Context) {
	projectID := c.Param("id")

	err := s.us.ViewProject(c.Request.Context(), projectID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, jsonResp(map[string]string{"status": "ok"}, nil))
}

func (s *Server) userApplyProject(c *gin.Context) {
	projectID := c.Param("id")

	err := s.us.ApplyToProject(c.Request.Context(), projectID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, jsonResp(map[string]string{"status": "ok"}, nil))
}
