package handlers

import (
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
)

// GetPatientHistory GET /v1/patients/:id/history
func (h *Handlers) GetPatientHistory(c *gin.Context) {
    id := c.Param("id")
    var result sql.NullString
    err := h.DB.QueryRow(`SELECT get_patient_wound_history($1)`, id).Scan(&result)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "history not found"})
            return
        }
        h.Log.Sugar().Errorf("get_patient_wound_history: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch history"})
        return
    }
    if !result.Valid {
        c.JSON(http.StatusOK, gin.H{"data": nil})
        return
    }
    c.Data(http.StatusOK, "application/json", []byte(result.String))
}
