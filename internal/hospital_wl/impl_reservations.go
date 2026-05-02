package hospital_wl

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implReservationsAPI struct{}

func NewReservationsApi() ReservationsAPI {
	return &implReservationsAPI{}
}

func (o implReservationsAPI) GetReservations(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		result := department.Reservations
		if result == nil {
			result = []Reservation{}
		}
		return nil, result, http.StatusOK
	})
}

func (o implReservationsAPI) CreateReservation(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		var reservation Reservation

		if err := c.ShouldBindJSON(&reservation); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if reservation.PatientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Patient ID is required",
			}, http.StatusBadRequest
		}

		if reservation.Id == "" || reservation.Id == "@new" {
			reservation.Id = uuid.NewString()
		}

		if reservation.Status == "" {
			reservation.Status = "pending"
		}

		if reservation.From.IsZero() {
			reservation.From = time.Now()
		}

		if reservation.To.IsZero() {
			reservation.To = time.Now().Add(24 * time.Hour)
		}

		conflictIdx := slices.IndexFunc(department.Reservations, func(r Reservation) bool {
			return r.Id == reservation.Id
		})
		if conflictIdx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Reservation already exists",
			}, http.StatusConflict
		}

		department.Reservations = append(department.Reservations, reservation)
		return department, reservation, http.StatusOK
	})
}

func (o implReservationsAPI) GetReservation(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		reservationId := c.Param("reservationId")
		if reservationId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Reservation ID is required",
			}, http.StatusBadRequest
		}

		idx := slices.IndexFunc(department.Reservations, func(r Reservation) bool {
			return r.Id == reservationId
		})
		if idx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Reservation not found",
			}, http.StatusNotFound
		}

		return nil, department.Reservations[idx], http.StatusOK
	})
}

func (o implReservationsAPI) UpdateReservation(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		var updated Reservation
		if err := c.ShouldBindJSON(&updated); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		reservationId := c.Param("reservationId")
		if reservationId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Reservation ID is required",
			}, http.StatusBadRequest
		}

		idx := slices.IndexFunc(department.Reservations, func(r Reservation) bool {
			return r.Id == reservationId
		})
		if idx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Reservation not found",
			}, http.StatusNotFound
		}

		existing := &department.Reservations[idx]

		if updated.PatientId != "" {
			existing.PatientId = updated.PatientId
		}
		if updated.PatientName != "" {
			existing.PatientName = updated.PatientName
		}
		if updated.Department != "" {
			existing.Department = updated.Department
		}
		if updated.Reason != "" {
			existing.Reason = updated.Reason
		}
		if !updated.From.IsZero() {
			existing.From = updated.From
		}
		if !updated.To.IsZero() {
			existing.To = updated.To
		}
		if updated.ContactInfo != "" {
			existing.ContactInfo = updated.ContactInfo
		}
		if updated.Status != "" {
			existing.Status = updated.Status
		}
		if updated.CancelReason != "" {
			existing.CancelReason = updated.CancelReason
		}
		if updated.Note != "" {
			existing.Note = updated.Note
		}
		if updated.RoomOrAmbulance != "" {
			existing.RoomOrAmbulance = updated.RoomOrAmbulance
		}

		return department, *existing, http.StatusOK
	})
}

func (o implReservationsAPI) DeleteReservation(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		reservationId := c.Param("reservationId")
		if reservationId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Reservation ID is required",
			}, http.StatusBadRequest
		}

		idx := slices.IndexFunc(department.Reservations, func(r Reservation) bool {
			return r.Id == reservationId
		})
		if idx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Reservation not found",
			}, http.StatusNotFound
		}

		department.Reservations = append(department.Reservations[:idx], department.Reservations[idx+1:]...)
		return department, nil, http.StatusNoContent
	})
}
