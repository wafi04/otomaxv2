package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/model"
	"github.com/wafi04/otomaxv2/internal/services"
	"github.com/wafi04/otomaxv2/pkg/response"
)

// Response struktur untuk standardize response

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input model.CreateCategory
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	err := h.categoryService.CreateCategory(c.Request.Context(), input)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "Category created successfully", nil)
}

func (h *CategoryHandler) GetCategoryByCode(c *gin.Context) {
	codeParam := c.Param("code")

	cat, err := h.categoryService.GetCategoryByCode(c.Request.Context(), codeParam)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Category not found", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", cat)
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	// Menggunakan pagination utility
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	search := c.DefaultQuery("search", "")
	filterType := c.DefaultQuery("type", "")
	active := c.DefaultQuery("status", "")

	// Hitung pagination
	paginationResult := response.CalculatePagination(&page, &limit)

	// Panggil service dengan response parameters
	data, totalCount, err := h.categoryService.GetAllCategories(
		c.Request.Context(),
		paginationResult.Skip,
		paginationResult.Take,
		search,
		filterType,
		active,
	)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch categories", err.Error())
		return
	}

	// Buat paginated response
	responses := response.CreatePaginatedResponse(
		data,
		paginationResult.CurrentPage,
		paginationResult.ItemsPerPage,
		totalCount,
	)

	response.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", responses)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid ID parameter", err.Error())
		return
	}

	var input model.CreateCategory
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	err = h.categoryService.UpdateCategory(c.Request.Context(), id, input)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update category", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Category updated successfully", nil)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid ID parameter", err.Error())
		return
	}

	err = h.categoryService.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete category", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}