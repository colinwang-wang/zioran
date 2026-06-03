package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/response"
)

type CourseHandler struct {
	courseSvc *service.CourseService
}

func NewCourseHandler(courseSvc *service.CourseService) *CourseHandler {
	return &CourseHandler{courseSvc: courseSvc}
}

func (h *CourseHandler) List(c *gin.Context) {
	var req model.CourseListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.courseSvc.List(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err.(*errcode.Error))
		return
	}
	response.Success(c, result)
}

func (h *CourseHandler) Detail(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.Error(c, errcode.ErrParam)
		return
	}
	var userID int64
	if id, exists := c.Get("user_id"); exists {
		userID = id.(int64)
	}
	result, err := h.courseSvc.Detail(c.Request.Context(), slug, userID)
	if err != nil {
		if e, ok := err.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *CourseHandler) Latest(c *gin.Context) {
	result, err := h.courseSvc.Latest(c.Request.Context())
	if err != nil {
		response.Error(c, err.(*errcode.Error))
		return
	}
	response.Success(c, result)
}

func (h *CourseHandler) Like(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, svcErr := h.courseSvc.ToggleLike(c.Request.Context(), userID.(int64), courseID)
	if svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *CourseHandler) Categories(c *gin.Context) {
	result, err := h.courseSvc.Categories(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *CourseHandler) Tags(c *gin.Context) {
	result, err := h.courseSvc.Tags(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *CourseHandler) Search(c *gin.Context) {
	var req model.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.courseSvc.Search(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err.(*errcode.Error))
		return
	}
	response.Success(c, result)
}
