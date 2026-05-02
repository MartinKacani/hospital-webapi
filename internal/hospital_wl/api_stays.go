package hospital_wl

import (
	"github.com/gin-gonic/gin"
)

// StaysAPI defines CRUD operations for managing hospitalization stays.
type StaysAPI interface {
	// GetStays Get /api/medbed/:departmentId/stays
	// Returns list of hospitalization stays for the department
	GetStays(c *gin.Context)

	// CreateStay Post /api/medbed/:departmentId/stays
	// Creates a new hospitalization stay
	CreateStay(c *gin.Context)

	// GetStay Get /api/medbed/:departmentId/stays/:stayId
	// Returns details of a single stay
	GetStay(c *gin.Context)

	// UpdateStay Put /api/medbed/:departmentId/stays/:stayId
	// Updates a specific stay
	UpdateStay(c *gin.Context)

	// DeleteStay Delete /api/medbed/:departmentId/stays/:stayId
	// Cancels/ends a hospitalization stay
	DeleteStay(c *gin.Context)
}
