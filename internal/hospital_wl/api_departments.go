package hospital_wl

import (
	"github.com/gin-gonic/gin"
)

// DepartmentsAPI defines operations for managing hospital departments.
type DepartmentsAPI interface {
	// CreateDepartment Post /api/medbed/department
	// Creates a new department
	CreateDepartment(c *gin.Context)

	// DeleteDepartment Delete /api/medbed/department/:departmentId
	// Deletes a specific department
	DeleteDepartment(c *gin.Context)
}
