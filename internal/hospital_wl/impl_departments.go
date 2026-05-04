package hospital_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/MartinKacani/hospital-webapi/internal/db_service"
)

type implDepartmentsAPI struct{}

func NewDepartmentsApi() DepartmentsAPI {
	return &implDepartmentsAPI{}
}

func (o implDepartmentsAPI) GetDepartments(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service not found"})
		return
	}
	db, ok := value.(db_service.DbService[Department])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service context is not of expected type"})
		return
	}
	departments, err := db.FindAllDocuments(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "Bad Gateway", "message": "Failed to load departments", "error": err.Error()})
		return
	}
	result := make([]Department, 0, len(departments))
	for _, d := range departments {
		result = append(result, Department{Id: d.Id, Name: d.Name, Code: d.Code})
	}
	c.JSON(http.StatusOK, result)
}

func (o implDepartmentsAPI) CreateDepartment(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service not found",
		})
		return
	}
	db, ok := value.(db_service.DbService[Department])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service context is not of expected type",
		})
		return
	}

	var department Department
	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if department.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Department ID is required",
		})
		return
	}

	if department.Reservations == nil {
		department.Reservations = []Reservation{}
	}
	if department.Stays == nil {
		department.Stays = []HospitalizationStay{}
	}

	if err := db.CreateDocument(c, department.Id, &department); err != nil {
		switch err {
		case db_service.ErrConflict:
			c.JSON(http.StatusConflict, gin.H{
				"status":  http.StatusConflict,
				"message": "Department already exists",
			})
		default:
			c.JSON(http.StatusBadGateway, gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create department in database",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, department)
}

func (o implDepartmentsAPI) DeleteDepartment(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service not found",
		})
		return
	}
	db, ok := value.(db_service.DbService[Department])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service context is not of expected type",
		})
		return
	}

	departmentId := c.Param("departmentId")
	if departmentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Department ID is required",
		})
		return
	}

	if err := db.DeleteDocument(c, departmentId); err != nil {
		switch err {
		case db_service.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "Not Found",
				"message": "Department not found",
			})
		default:
			c.JSON(http.StatusBadGateway, gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete department from database",
				"error":   err.Error(),
			})
		}
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
