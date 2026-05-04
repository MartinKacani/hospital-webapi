package hospital_wl

import "time"

// applyAutoStatus updates stays with time-driven status transitions:
//   - planned → active  when now >= from
//   - active  → completed when now > to
//
// cancelled and completed are terminal and are never changed.
// Returns true if any stay was modified so the caller can persist the department.
func applyAutoStatus(department *Department) bool {
	now := time.Now()
	changed := false
	for i := range department.Stays {
		s := &department.Stays[i]
		switch s.Status {
		case "planned":
			if !now.Before(s.From) {
				s.Status = "active"
				changed = true
			}
		case "active":
			if now.After(s.To) {
				s.Status = "completed"
				changed = true
			}
		}
	}
	return changed
}
