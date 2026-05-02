package hospital_wl

import "time"

// Reservation represents a reservation for examination or planned hospitalization.
type Reservation struct {
	// Unique identifier of the reservation
	Id string `json:"id"`

	// Unique patient identifier
	PatientId string `json:"patientId"`

	// Full name of the patient
	PatientName string `json:"patientName"`

	// Target department ID
	Department string `json:"department"`

	// Reason for the reservation
	Reason string `json:"reason"`

	// Start date and time of the reservation
	From time.Time `json:"from"`

	// End date and time of the reservation
	To time.Time `json:"to"`

	// Contact information (patient or referring doctor)
	ContactInfo string `json:"contactInfo,omitempty"`

	// Status: pending | confirmed | cancelled
	Status string `json:"status"`

	// Reason for cancellation
	CancelReason string `json:"cancelReason,omitempty"`

	// Note from the hospital
	Note string `json:"note,omitempty"`

	// Assigned room or ambulance
	RoomOrAmbulance string `json:"roomOrAmbulance,omitempty"`
}
