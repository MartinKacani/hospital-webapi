package hospital_wl

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implStaysAPI struct{}

func NewStaysApi() StaysAPI {
	return &implStaysAPI{}
}

func (o implStaysAPI) GetStays(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		changed := applyAutoStatus(department)
		result := department.Stays
		if result == nil {
			result = []HospitalizationStay{}
		}
		if changed {
			return department, result, http.StatusOK
		}
		return nil, result, http.StatusOK
	})
}

func (o implStaysAPI) CreateStay(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		var stay HospitalizationStay

		if err := c.ShouldBindJSON(&stay); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if stay.PatientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Patient ID is required",
			}, http.StatusBadRequest
		}

		if stay.RoomNumber == "" || stay.BedNumber == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Room number and bed number are required",
			}, http.StatusBadRequest
		}

		if stay.Id == "" || stay.Id == "@new" {
			stay.Id = uuid.NewString()
		}

		if stay.Status == "" {
			stay.Status = "planned"
		}

		if stay.From.IsZero() {
			stay.From = time.Now()
		}

		if stay.To.IsZero() {
			stay.To = time.Now().Add(24 * time.Hour)
		}

		conflictIdx := slices.IndexFunc(department.Stays, func(s HospitalizationStay) bool {
			return s.Id == stay.Id
		})
		if conflictIdx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Stay already exists",
			}, http.StatusConflict
		}

		department.Stays = append(department.Stays, stay)

		if stay.ReservationId != "" {
			resIdx := slices.IndexFunc(department.Reservations, func(r Reservation) bool {
				return r.Id == stay.ReservationId
			})
			if resIdx >= 0 && department.Reservations[resIdx].Status == "pending" {
				department.Reservations[resIdx].Status = "confirmed"
			}
		}

		return department, stay, http.StatusOK
	})
}

func (o implStaysAPI) GetStay(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		stayId := c.Param("stayId")
		if stayId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Stay ID is required",
			}, http.StatusBadRequest
		}

		changed := applyAutoStatus(department)

		idx := slices.IndexFunc(department.Stays, func(s HospitalizationStay) bool {
			return s.Id == stayId
		})
		if idx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Stay not found",
			}, http.StatusNotFound
		}

		if changed {
			return department, department.Stays[idx], http.StatusOK
		}
		return nil, department.Stays[idx], http.StatusOK
	})
}

func (o implStaysAPI) UpdateStay(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		var updated HospitalizationStay
		if err := c.ShouldBindJSON(&updated); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		stayId := c.Param("stayId")
		if stayId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Stay ID is required",
			}, http.StatusBadRequest
		}

		idx := slices.IndexFunc(department.Stays, func(s HospitalizationStay) bool {
			return s.Id == stayId
		})
		if idx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Stay not found",
			}, http.StatusNotFound
		}

		existing := &department.Stays[idx]

		if updated.ReservationId != "" {
			existing.ReservationId = updated.ReservationId
		}
		if updated.PatientId != "" {
			existing.PatientId = updated.PatientId
		}
		if updated.PatientName != "" {
			existing.PatientName = updated.PatientName
		}
		if updated.Department != "" {
			existing.Department = updated.Department
		}
		if updated.RoomNumber != "" {
			existing.RoomNumber = updated.RoomNumber
		}
		if updated.BedNumber != "" {
			existing.BedNumber = updated.BedNumber
		}
		if !updated.From.IsZero() {
			existing.From = updated.From
		}
		if !updated.To.IsZero() {
			existing.To = updated.To
		}
		if updated.Status != "" {
			existing.Status = updated.Status
		}
		if updated.CancelReason != "" {
			existing.CancelReason = updated.CancelReason
		}
		if updated.Notes != "" {
			existing.Notes = updated.Notes
		}

		return department, *existing, http.StatusOK
	})
}

func (o implStaysAPI) DeleteStay(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		stayId := c.Param("stayId")
		if stayId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Stay ID is required",
			}, http.StatusBadRequest
		}

		idx := slices.IndexFunc(department.Stays, func(s HospitalizationStay) bool {
			return s.Id == stayId
		})
		if idx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Stay not found",
			}, http.StatusNotFound
		}

		deletedStay := department.Stays[idx]
		department.Stays = append(department.Stays[:idx], department.Stays[idx+1:]...)

		if deletedStay.ReservationId != "" {
			resIdx := slices.IndexFunc(department.Reservations, func(r Reservation) bool {
				return r.Id == deletedStay.ReservationId
			})
			if resIdx >= 0 && department.Reservations[resIdx].Status == "confirmed" {
				department.Reservations[resIdx].Status = "pending"
			}
		}

		return department, nil, http.StatusNoContent
	})
}
