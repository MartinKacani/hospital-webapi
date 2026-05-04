package hospital_wl

import (
	"context"
	"log"

	"github.com/MartinKacani/hospital-webapi/internal/db_service"
)

var defaultDepartments = []Department{
	{Id: "cardiology", Name: "Kardiologické oddelenie", Code: "KARD"},
	{Id: "neurology", Name: "Neurológia", Code: "NEUR"},
	{Id: "surgery", Name: "Chirurgia", Code: "CHIR"},
	{Id: "rheumatology", Name: "Reumatológia", Code: "REUM"},
	{Id: "orthopedics", Name: "Ortopédia", Code: "ORT"},
}

func SeedDepartments(ctx context.Context, db db_service.DbService[Department]) {
	for _, dept := range defaultDepartments {
		_, err := db.FindDocument(ctx, dept.Id)
		if err == db_service.ErrNotFound {
			d := dept
			d.Reservations = []Reservation{}
			d.Stays = []HospitalizationStay{}
			if err := db.CreateDocument(ctx, d.Id, &d); err != nil {
				log.Printf("Failed to seed department %s: %v", d.Id, err)
			} else {
				log.Printf("Seeded department: %s", d.Id)
			}
		}
	}
}
