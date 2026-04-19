package hospital_wl

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "slices"
)

type implHospitalWaitingListAPI struct {
}

func NewHospitalWaitingListApi() HospitalWaitingListAPI {
    return &implHospitalWaitingListAPI{}
}

func (o implHospitalWaitingListAPI) CreateWaitingListEntry(c *gin.Context) {
    updateHospitalFunc(c, func(c *gin.Context, hospital *Hospital) (*Hospital,  interface{},  int){
        var entry WaitingListEntry

        if err := c.ShouldBindJSON(&entry); err != nil {
            return nil, gin.H{
                "status": http.StatusBadRequest,
                "message": "Invalid request body",
                "error": err.Error(),
            }, http.StatusBadRequest
        }

        if entry.PatientId == "" {
            return nil, gin.H{
                "status": http.StatusBadRequest,
                "message": "Patient ID is required",
            }, http.StatusBadRequest
        }

        if entry.Id == "" || entry.Id == "@new" {
            entry.Id = uuid.NewString()
        }

        conflictIndx := slices.IndexFunc( hospital.WaitingList, func(waiting WaitingListEntry) bool {
            return entry.Id == waiting.Id || entry.PatientId == waiting.PatientId
        })

        if conflictIndx >= 0 {
            return nil, gin.H{
                "status": http.StatusConflict,
                "message": "Entry already exists",
            }, http.StatusConflict
        }

        hospital.WaitingList = append(hospital.WaitingList, entry)
        hospital.reconcileWaitingList()
        // entry was copied by value return reconciled value from the list
        entryIndx := slices.IndexFunc( hospital.WaitingList, func(waiting WaitingListEntry) bool {
            return entry.Id == waiting.Id
        })
        if entryIndx < 0 {
            return nil, gin.H{
                "status": http.StatusInternalServerError,
                "message": "Failed to save entry",
            }, http.StatusInternalServerError
        }
        return hospital, hospital.WaitingList[entryIndx], http.StatusOK
    })
}

func (o implHospitalWaitingListAPI) DeleteWaitingListEntry(c *gin.Context) {
    updateHospitalFunc(c, func(c *gin.Context, hospital *Hospital) (*Hospital, interface{}, int) {
        entryId := c.Param("entryId")

        if entryId == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Entry ID is required",
            }, http.StatusBadRequest
        }

        entryIndx := slices.IndexFunc(hospital.WaitingList, func(waiting WaitingListEntry) bool {
            return entryId == waiting.Id
        })

        if entryIndx < 0 {
            return nil, gin.H{
                "status":  http.StatusNotFound,
                "message": "Entry not found",
            }, http.StatusNotFound
        }

        hospital.WaitingList = append(hospital.WaitingList[:entryIndx], hospital.WaitingList[entryIndx+1:]...)
        hospital.reconcileWaitingList()
        return hospital, nil, http.StatusNoContent
    })
}

func (o implHospitalWaitingListAPI) GetWaitingListEntries(c *gin.Context) {
    updateHospitalFunc(c, func(c *gin.Context, hospital *Hospital) (*Hospital, interface{}, int) {
        result := hospital.WaitingList
        if result == nil {
            result = []WaitingListEntry{}
        }
        // return nil ambulance - no need to update it in db
        return nil, result, http.StatusOK
    })
}

func (o implHospitalWaitingListAPI) GetWaitingListEntry(c *gin.Context) {
    updateHospitalFunc(c, func(c *gin.Context, hospital *Hospital) (*Hospital, interface{}, int) {
        entryId := c.Param("entryId")

        if entryId == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Entry ID is required",
            }, http.StatusBadRequest
        }

        entryIndx := slices.IndexFunc(hospital.WaitingList, func(waiting WaitingListEntry) bool {
            return entryId == waiting.Id
        })

        if entryIndx < 0 {
            return nil, gin.H{
                "status":  http.StatusNotFound,
                "message": "Entry not found",
            }, http.StatusNotFound
        }

        // return nil ambulance - no need to update it in db
        return nil, hospital.WaitingList[entryIndx], http.StatusOK
    })
}

func (o implHospitalWaitingListAPI) UpdateWaitingListEntry(c *gin.Context) {
    updateHospitalFunc(c, func(c *gin.Context, hospital *Hospital) (*Hospital, interface{}, int) {
        var entry WaitingListEntry

        if err := c.ShouldBindJSON(&entry); err != nil {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Invalid request body",
                "error":   err.Error(),
            }, http.StatusBadRequest
        }

        entryId := c.Param("entryId")

        if entryId == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Entry ID is required",
            }, http.StatusBadRequest
        }

        entryIndx := slices.IndexFunc(hospital.WaitingList, func(waiting WaitingListEntry) bool {
            return entryId == waiting.Id
        })

        if entryIndx < 0 {
            return nil, gin.H{
                "status":  http.StatusNotFound,
                "message": "Entry not found",
            }, http.StatusNotFound
        }

        if entry.PatientId != "" {
            hospital.WaitingList[entryIndx].PatientId = entry.PatientId
        }

        if entry.Id != "" {
            hospital.WaitingList[entryIndx].Id = entry.Id
        }

        if entry.WaitingSince.After(time.Time{}) {
            hospital.WaitingList[entryIndx].WaitingSince = entry.WaitingSince
        }

        if entry.EstimatedDurationMinutes > 0 {
            hospital.WaitingList[entryIndx].EstimatedDurationMinutes = entry.EstimatedDurationMinutes
        }

        hospital.reconcileWaitingList()
        return hospital, hospital.WaitingList[entryIndx], http.StatusOK
    })
}