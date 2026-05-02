package hospital_wl

// Department represents a hospital department that manages reservations and stays.
type Department struct {
	// Unique identifier of the department
	Id string `json:"id"`

	// Human readable display name of the department
	Name string `json:"name"`

	// Short department code (e.g. KARD, ORL)
	Code string `json:"code"`

	Reservations []Reservation         `json:"reservations,omitempty"`
	Stays        []HospitalizationStay `json:"stays,omitempty"`
}
