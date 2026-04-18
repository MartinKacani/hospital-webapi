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
    c.AbortWithStatus(http.StatusNotImplemented)
}