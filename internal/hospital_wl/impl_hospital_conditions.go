package hospital_wl

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type implHospitalConditionsAPI struct {
}

func NewHospitalConditionsApi() HospitalConditionsAPI {
    return &implHospitalConditionsAPI{}
}

func (o implHospitalConditionsAPI) GetConditions(c *gin.Context) {
    updateHospitalFunc(c, func(
        c *gin.Context,
        hospital *Hospital,
    ) (updatedHospital *Hospital, responseContent interface{}, status int) {
        result := hospital.PredefinedConditions
        if result == nil {
            result = []Condition{}
        }
        return nil, result, http.StatusOK
    })
}