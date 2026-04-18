package hospital_wl

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type implHospitalWaitingListAPI struct {
}

func NewHospitalWaitingListApi() HospitalWaitingListAPI {
    return &implHospitalWaitingListAPI{}
}

func (o implHospitalWaitingListAPI) CreateWaitingListEntry(c *gin.Context) {
    c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implHospitalWaitingListAPI) DeleteWaitingListEntry(c *gin.Context) {
    c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implHospitalWaitingListAPI) GetWaitingListEntries(c *gin.Context) {
    c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implHospitalWaitingListAPI) GetWaitingListEntry(c *gin.Context) {
    c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implHospitalWaitingListAPI) UpdateWaitingListEntry(c *gin.Context) {
    c.AbortWithStatus(http.StatusNotImplemented)
}