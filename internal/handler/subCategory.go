package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/model"
	"github.com/wafi04/otomaxv2/internal/services"
	"github.com/wafi04/otomaxv2/pkg/response"
)

type SubCategoryHandler struct {
	subCategoryService *services.SubCategoryService
}

func NewSubCategoryHandler(subCategoryService *services.SubCategoryService) *SubCategoryHandler {
	return &SubCategoryHandler{
		subCategoryService: subCategoryService,
	}
}

// Create SubCategory Handler
func (h *SubCategoryHandler) CreateSubCategory(c *gin.Context) {
	var req model.CreateSubcategory
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	subCategory, err := h.subCategoryService.CreateSubCategory(c.Request.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "already exists") {
			response.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
			return
		}
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create services", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "SubCategory created successfully", subCategory)
}

// Get All SubCategories Handler
func (h *SubCategoryHandler) GetAllSubCategories(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	search := c.DefaultQuery("search", "")
	status := c.DefaultQuery("status", "")

	// Validasi status parameter
	if status != "" && status != "active" && status != "inactive" {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid status parameter", "status must be 'active' or 'inactive'")
		return
	}

	// Hitung pagination
	paginationResult := response.CalculatePagination(&page, &limit)

	// Panggil service
	data, totalCount, err := h.subCategoryService.GetAllSubCategories(
		c.Request.Context(),
		paginationResult.Skip,
		paginationResult.Take,
		search,
		status,
	)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch subcategories", err.Error())
		return
	}

	// Buat paginated response
	responses := response.CreatePaginatedResponse(
		data,
		paginationResult.CurrentPage,
		paginationResult.ItemsPerPage,
		totalCount,
	)

	response.SuccessResponse(c, http.StatusOK, "SubCategories retrieved successfully", responses)
}

// Get SubCategories by Category ID Handler
func (h *SubCategoryHandler) GetSubCategoriesByCategoryID(c *gin.Context) {
	categoryIDStr := c.Param("categoryId")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	search := c.DefaultQuery("search", "")
	status := c.DefaultQuery("status", "")

	// Validasi status parameter
	if status != "" && status != "active" && status != "inactive" {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid status parameter", "status must be 'active' or 'inactive'")
		return
	}

	// Hitung pagination
	paginationResult := response.CalculatePagination(&page, &limit)

	// Panggil service
	data, totalCount, err := h.subCategoryService.GetSubCategoriesByCategoryID(
		c.Request.Context(),
		categoryID,
		paginationResult.Skip,
		paginationResult.Take,
		search,
		status,
	)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch subcategories", err.Error())
		return
	}

	// Buat paginated response
	responses := response.CreatePaginatedResponse(
		data,
		paginationResult.CurrentPage,
		paginationResult.ItemsPerPage,
		totalCount,
	)

	response.SuccessResponse(c, http.StatusOK, "SubCategories retrieved successfully", responses)
}

// Get SubCategory by ID Handler
func (h *SubCategoryHandler) GetSubCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid services ID", err.Error())
		return
	}

	subCategory, err := h.subCategoryService.GetSubCategoryByID(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.ErrorResponse(c, http.StatusNotFound, "SubCategory not found", err.Error())
			return
		}
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch services", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "SubCategory retrieved successfully", subCategory)
}

// Update SubCategory Handler
func (h *SubCategoryHandler) UpdateSubCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid services ID", err.Error())
		return
	}

	var req model.UpdateSubcategory
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	subCategory, err := h.subCategoryService.UpdateSubCategory(c.Request.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.ErrorResponse(c, http.StatusNotFound, "SubCategory not found", err.Error())
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			response.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
			return
		}
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update services", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "SubCategory updated successfully", subCategory)
}

// Delete SubCategory Handler (Soft Delete)
func (h *SubCategoryHandler) DeleteSubCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid services ID", err.Error())
		return
	}

	err = h.subCategoryService.DeleteSubCategory(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.ErrorResponse(c, http.StatusNotFound, "SubCategory not found", err.Error())
			return
		}
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete services", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "SubCategory deleted successfully", nil)
}