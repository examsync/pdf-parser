package controllers

import (
	"net/http"

	"github.com/examsync/pdf-parser/internal/services"
	"github.com/labstack/echo/v5"
)

// ParsedPDFController handles HTTP requests for parsed PDFs.
type ParsedPDFController struct {
	service *services.ParsedPDFService
}

// NewParsedPDFController creates a new instance of ParsedPDFController.
func NewParsedPDFController(service *services.ParsedPDFService) *ParsedPDFController {
	return &ParsedPDFController{service: service}
}

// GetPDFs handles the HTTP request to get all parsed PDFs.
func (c *ParsedPDFController) GetPDFs(ctx *echo.Context) error {
	pdfs, err := c.service.GetPDFs()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, pdfs)
}
