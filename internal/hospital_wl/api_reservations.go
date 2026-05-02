package hospital_wl

import (
	"github.com/gin-gonic/gin"
)

// ReservationsAPI defines CRUD operations for managing reservations.
type ReservationsAPI interface {
	// GetReservations Get /api/medbed/:departmentId/reservations
	// Returns list of reservations for the department
	GetReservations(c *gin.Context)

	// CreateReservation Post /api/medbed/:departmentId/reservations
	// Creates a new reservation
	CreateReservation(c *gin.Context)

	// GetReservation Get /api/medbed/:departmentId/reservations/:reservationId
	// Returns details of a single reservation
	GetReservation(c *gin.Context)

	// UpdateReservation Put /api/medbed/:departmentId/reservations/:reservationId
	// Updates a specific reservation
	UpdateReservation(c *gin.Context)

	// DeleteReservation Delete /api/medbed/:departmentId/reservations/:reservationId
	// Cancels/deletes a reservation
	DeleteReservation(c *gin.Context)
}
