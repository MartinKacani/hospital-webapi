package hospital_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/MartinKacani/hospital-webapi/internal/db_service"
)

type departmentUpdater = func(
	ctx *gin.Context,
	department *Department,
) (updatedDepartment *Department, responseContent interface{}, status int)

func updateDepartmentFunc(ctx *gin.Context, updater departmentUpdater) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Department])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	departmentId := ctx.Param("departmentId")

	department, err := db.FindDocument(ctx, departmentId)

	switch err {
	case nil:
		// continue
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Department not found",
				"error":   err.Error(),
			},
		)
		return
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to load department from database",
				"error":   err.Error(),
			})
		return
	}

	updatedDepartment, responseObject, status := updater(ctx, department)

	if updatedDepartment != nil {
		err = db.UpdateDocument(ctx, departmentId, updatedDepartment)
	} else {
		err = nil
	}

	switch err {
	case nil:
		if responseObject != nil {
			ctx.JSON(status, responseObject)
		} else {
			ctx.AbortWithStatus(status)
		}
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Department was deleted while processing the request",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to update department in database",
				"error":   err.Error(),
			})
	}
}
