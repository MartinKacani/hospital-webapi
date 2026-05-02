package hospital_wl

import "time"

// HospitalizationStay represents an actual hospitalization stay with an assigned bed.
type HospitalizationStay struct {
	// Unique identifier of the stay
	Id string `json:"id"`

	// ID of the linked reservation (optional)
	ReservationId string `json:"reservationId,omitempty"`

	// Unique patient identifier
	PatientId string `json:"patientId"`

	// Full name of the patient
	PatientName string `json:"patientName"`

	// Department where the patient is staying
	Department string `json:"department"`

	// Room number
	RoomNumber string `json:"roomNumber"`

	// Bed number within the room
	BedNumber string `json:"bedNumber"`

	// Start of hospitalization
	From time.Time `json:"from"`

	// End of hospitalization
	To time.Time `json:"to"`

	// Status: planned | active | completed | cancelled
	Status string `json:"status"`

	// Reason for cancellation or early discharge
	CancelReason string `json:"cancelReason,omitempty"`

	// Clinical notes for the stay
	Notes string `json:"notes,omitempty"`
}
