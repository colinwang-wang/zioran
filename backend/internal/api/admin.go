package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/response"
)

type AdminHandler struct {
	courseSvc *service.CourseService
}

func NewAdminHandler(courseSvc *service.CourseService) *AdminHandler {
	return &AdminHandler{courseSvc: courseSvc}
}

// Courses

func (h *AdminHandler) CourseList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	categoryID, _ := strconv.Atoi(c.Query("categoryId"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	result, err := h.courseSvc.AdminList(c.Request.Context(), page, pageSize, categoryID, status, keyword)
	if err != nil {
		response.Error(c, err.(*errcode.Error))
		return
	}
	response.Success(c, result)
}

func (h *AdminHandler) CourseCreate(c *gin.Context) {
	var req model.AdminCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.courseSvc.AdminCreate(c.Request.Context(), &req)
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

func (h *AdminHandler) CourseUpdate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.AdminCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, svcErr := h.courseSvc.AdminUpdate(c.Request.Context(), id, &req)
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

func (h *AdminHandler) CourseDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.courseSvc.AdminDelete(c.Request.Context(), id); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) CourseUpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.AdminCourseStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.courseSvc.AdminUpdateStatus(c.Request.Context(), id, req.Status); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// Categories

func (h *AdminHandler) CategoryList(c *gin.Context) {
	result, err := h.courseSvc.AdminCategoryList(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminHandler) CategoryCreate(c *gin.Context) {
	var req model.AdminCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.courseSvc.AdminCategoryCreate(c.Request.Context(), &req)
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

func (h *AdminHandler) CategoryUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.AdminCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, svcErr := h.courseSvc.AdminCategoryUpdate(c.Request.Context(), id, &req)
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

func (h *AdminHandler) CategoryDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.courseSvc.AdminCategoryDelete(c.Request.Context(), id); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// Tags

func (h *AdminHandler) TagList(c *gin.Context) {
	result, err := h.courseSvc.AdminTagList(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminHandler) TagCreate(c *gin.Context) {
	var req model.AdminTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.courseSvc.AdminTagCreate(c.Request.Context(), &req)
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

func (h *AdminHandler) TagUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.AdminTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, svcErr := h.courseSvc.AdminTagUpdate(c.Request.Context(), id, &req)
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

func (h *AdminHandler) TagDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.courseSvc.AdminTagDelete(c.Request.Context(), id); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}
