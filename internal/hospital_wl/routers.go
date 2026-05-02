package hospital_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// NewRouter returns a new router.
func NewRouter(handleFunctions ApiHandleFunctions) *gin.Engine {
	return NewRouterWithGinEngine(gin.Default(), handleFunctions)
}

// NewRouterWithGinEngine adds routes to existing gin engine.
func NewRouterWithGinEngine(router *gin.Engine, handleFunctions ApiHandleFunctions) *gin.Engine {
	for _, route := range getRoutes(handleFunctions) {
		if route.HandlerFunc == nil {
			route.HandlerFunc = DefaultHandleFunc
		}
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

// DefaultHandleFunc is the default handler for not yet implemented routes.
func DefaultHandleFunc(c *gin.Context) {
	c.String(http.StatusNotImplemented, "501 not implemented")
}

// ApiHandleFunctions holds all API handler interfaces.
type ApiHandleFunctions struct {
	// Routes for the DepartmentsAPI part of the API
	DepartmentsAPI DepartmentsAPI
	// Routes for the ReservationsAPI part of the API
	ReservationsAPI ReservationsAPI
	// Routes for the StaysAPI part of the API
	StaysAPI StaysAPI
}

func getRoutes(handleFunctions ApiHandleFunctions) []Route {
	return []Route{
		// Departments
		{
			"CreateDepartment",
			http.MethodPost,
			"/api/medbed/department",
			handleFunctions.DepartmentsAPI.CreateDepartment,
		},
		{
			"DeleteDepartment",
			http.MethodDelete,
			"/api/medbed/department/:departmentId",
			handleFunctions.DepartmentsAPI.DeleteDepartment,
		},
		// Reservations
		{
			"GetReservations",
			http.MethodGet,
			"/api/medbed/:departmentId/reservations",
			handleFunctions.ReservationsAPI.GetReservations,
		},
		{
			"CreateReservation",
			http.MethodPost,
			"/api/medbed/:departmentId/reservations",
			handleFunctions.ReservationsAPI.CreateReservation,
		},
		{
			"GetReservation",
			http.MethodGet,
			"/api/medbed/:departmentId/reservations/:reservationId",
			handleFunctions.ReservationsAPI.GetReservation,
		},
		{
			"UpdateReservation",
			http.MethodPut,
			"/api/medbed/:departmentId/reservations/:reservationId",
			handleFunctions.ReservationsAPI.UpdateReservation,
		},
		{
			"DeleteReservation",
			http.MethodDelete,
			"/api/medbed/:departmentId/reservations/:reservationId",
			handleFunctions.ReservationsAPI.DeleteReservation,
		},
		// Stays
		{
			"GetStays",
			http.MethodGet,
			"/api/medbed/:departmentId/stays",
			handleFunctions.StaysAPI.GetStays,
		},
		{
			"CreateStay",
			http.MethodPost,
			"/api/medbed/:departmentId/stays",
			handleFunctions.StaysAPI.CreateStay,
		},
		{
			"GetStay",
			http.MethodGet,
			"/api/medbed/:departmentId/stays/:stayId",
			handleFunctions.StaysAPI.GetStay,
		},
		{
			"UpdateStay",
			http.MethodPut,
			"/api/medbed/:departmentId/stays/:stayId",
			handleFunctions.StaysAPI.UpdateStay,
		},
		{
			"DeleteStay",
			http.MethodDelete,
			"/api/medbed/:departmentId/stays/:stayId",
			handleFunctions.StaysAPI.DeleteStay,
		},
	}
}
