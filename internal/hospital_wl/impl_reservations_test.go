package hospital_wl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MartinKacani/hospital-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MedBedSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[Department]
}

func TestMedBedSuite(t *testing.T) {
	suite.Run(t, new(MedBedSuite))
}

type DbServiceMock[DocType interface{}] struct {
	mock.Mock
}

func (this *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := this.Called(ctx, id)
	return args.Get(0).(*DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := this.Called(ctx, id)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := this.Called(ctx)
	return args.Error(0)
}

func (suite *MedBedSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[Department]{}

	// Compile time Assert that the mock is of type db_service.DbService[Department]
	var _ db_service.DbService[Department] = suite.dbServiceMock

	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Department{
				Id:   "test-department",
				Name: "Test Department",
				Code: "TEST",
				Reservations: []Reservation{
					{
						Id:          "test-reservation",
						PatientId:   "test-patient",
						PatientName: "Test Patient",
						Department:  "test-department",
						Reason:      "Test reason",
						From:        time.Now(),
						To:          time.Now().Add(24 * time.Hour),
						Status:      "pending",
					},
				},
				Stays: []HospitalizationStay{},
			},
			nil,
		)
}

func (suite *MedBedSuite) Test_UpdateReservation_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
		"id": "test-reservation",
		"patientId": "test-patient",
		"patientName": "Updated Name",
		"department": "test-department",
		"reason": "Updated reason",
		"status": "confirmed"
	}`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "departmentId", Value: "test-department"},
		{Key: "reservationId", Value: "test-reservation"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/api/medbed/test-department/reservations/test-reservation", strings.NewReader(json))

	sut := implReservationsAPI{}

	// ACT
	sut.UpdateReservation(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-department", mock.Anything)
	suite.Equal(http.StatusOK, recorder.Code)
}

func (suite *MedBedSuite) Test_GetReservations_ReturnsAllReservations() {
	// ARRANGE
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "departmentId", Value: "test-department"},
	}
	ctx.Request = httptest.NewRequest("GET", "/api/medbed/test-department/reservations", nil)

	sut := implReservationsAPI{}

	// ACT
	sut.GetReservations(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	suite.dbServiceMock.AssertNotCalled(suite.T(), "UpdateDocument")
}

func (suite *MedBedSuite) Test_DeleteReservation_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "departmentId", Value: "test-department"},
		{Key: "reservationId", Value: "test-reservation"},
	}
	ctx.Request = httptest.NewRequest("DELETE", "/api/medbed/test-department/reservations/test-reservation", nil)

	sut := implReservationsAPI{}

	// ACT
	sut.DeleteReservation(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-department", mock.Anything)
	suite.Equal(http.StatusNoContent, recorder.Code)
}
