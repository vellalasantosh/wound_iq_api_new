package handlers

import (
    "database/sql"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/vellalasantosh/wound_iq_api/internal/models"
)

// ListAssessments GET /v1/assessments?patient_id=&clinician_id=&date_from=&date_to=&page=&page_size=
func (h *Handlers) ListAssessments(c *gin.Context) {
    page, pageSize := parsePagination(c)
    offset := (page - 1) * pageSize

    var args []interface{}
    where := []string{}
    idx := 1

    if v := c.Query("patient_id"); v != "" {
        where = append(where, "patient_id = $"+strconv.Itoa(idx))
        args = append(args, v)
        idx++
    }
    if v := c.Query("clinician_id"); v != "" {
        where = append(where, "clinician_id = $"+strconv.Itoa(idx))
        args = append(args, v)
        idx++
    }
    if v := c.Query("date_from"); v != "" {
        if t, err := time.Parse(time.RFC3339, v); err == nil {
            where = append(where, "created_at >= $"+strconv.Itoa(idx))
            args = append(args, t)
            idx++
        }
    }
    if v := c.Query("date_to"); v != "" {
        if t, err := time.Parse(time.RFC3339, v); err == nil {
            where = append(where, "created_at <= $"+strconv.Itoa(idx))
            args = append(args, t)
            idx++
        }
    }

    base := `SELECT id, patient_id, clinician_id, wound_id, notes, created_at, updated_at FROM assessments`
    if len(where) > 0 {
        base += " WHERE " + strings.Join(where, " AND ")
    }
    base += " ORDER BY id DESC LIMIT $" + strconv.Itoa(idx) + " OFFSET $" + strconv.Itoa(idx+1)
    args = append(args, pageSize, offset)

    rows, err := h.DB.Query(base, args...)
    if err != nil {
        h.Log.Sugar().Errorf("list assessments: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch assessments"})
        return
    }
    defer rows.Close()

    out := []models.Assessment{}
    for rows.Next() {
        var a models.Assessment
        var woundID sql.NullInt64
        if err := rows.Scan(&a.ID, &a.PatientID, &a.ClinicianID, &woundID, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
            h.Log.Sugar().Errorf("scan assessment: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read assessments"})
            return
        }
        if woundID.Valid {
            v := woundID.Int64
            a.WoundID = &v
        }
        out = append(out, a)
    }
    c.JSON(http.StatusOK, gin.H{"data": out, "page": page, "page_size": pageSize})
}

// GetAssessment GET /v1/assessments/:id
func (h *Handlers) GetAssessment(c *gin.Context) {
    id := c.Param("id")
    var a models.Assessment
    var woundID sql.NullInt64
    row := h.DB.QueryRow(`SELECT id, patient_id, clinician_id, wound_id, notes, created_at, updated_at FROM assessments WHERE id=$1`, id)
    if err := row.Scan(&a.ID, &a.PatientID, &a.ClinicianID, &woundID, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
            return
        }
        h.Log.Sugar().Errorf("get assessment: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get assessment"})
        return
    }
    if woundID.Valid {
        v := woundID.Int64
        a.WoundID = &v
    }
    c.JSON(http.StatusOK, a)
}

// CreateAssessment POST /v1/assessments
// Optionally calls add_full_assessment if available
func (h *Handlers) CreateAssessment(c *gin.Context) {
    var in struct {
        PatientID   int64  `json:"patient_id" binding:"required"`
        ClinicianID int64  `json:"clinician_id" binding:"required"`
        WoundID     *int64 `json:"wound_id"`
        Notes       string `json:"notes"`
    }
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // If DB has add_full_assessment function that accepts JSON or fields, adapt accordingly.
    var newID int64
    err := h.DB.QueryRow(`INSERT INTO assessments (patient_id, clinician_id, wound_id, notes, created_at, updated_at)
                          VALUES ($1, $2, $3, $4, now(), now()) RETURNING id`,
        in.PatientID, in.ClinicianID, in.WoundID, in.Notes).Scan(&newID)
    if err != nil {
        h.Log.Sugar().Errorf("create assessment: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create assessment"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"id": newID})
}

// UpdateAssessment PUT /v1/assessments/:id
func (h *Handlers) UpdateAssessment(c *gin.Context) {
    id := c.Param("id")
    var in struct {
        PatientID   *int64  `json:"patient_id"`
        ClinicianID *int64  `json:"clinician_id"`
        WoundID     *int64  `json:"wound_id"`
        Notes       *string `json:"notes"`
    }
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err := h.DB.Exec(`UPDATE assessments SET
                       patient_id = COALESCE($1, patient_id),
                       clinician_id = COALESCE($2, clinician_id),
                       wound_id = COALESCE($3, wound_id),
                       notes = COALESCE($4, notes),
                       updated_at = now()
                       WHERE id = $5`,
        in.PatientID, in.ClinicianID, in.WoundID, in.Notes, id)
    if err != nil {
        h.Log.Sugar().Errorf("update assessment: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update assessment"})
        return
    }
    c.Status(http.StatusNoContent)
}

// DeleteAssessment DELETE /v1/assessments/:id
func (h *Handlers) DeleteAssessment(c *gin.Context) {
    id := c.Param("id")
    res, err := h.DB.Exec(`DELETE FROM assessments WHERE id = $1`, id)
    if err != nil {
        h.Log.Sugar().Errorf("delete assessment: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete assessment"})
        return
    }
    rows, _ := res.RowsAffected()
    if rows == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
        return
    }
    c.Status(http.StatusNoContent)
}

// GetAssessmentFull uses DB function get_assessment_full(assessment_id)
func (h *Handlers) GetAssessmentFull(c *gin.Context) {
    id := c.Param("id")
    var fullJSON sql.NullString
    err := h.DB.QueryRow(`SELECT get_assessment_full($1)`, id).Scan(&fullJSON)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
            return
        }
        h.Log.Sugar().Errorf("get_assessment_full: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch full assessment"})
        return
    }
    if !fullJSON.Valid || fullJSON.String == "" {
        c.JSON(http.StatusOK, gin.H{"data": nil})
        return
    }
    c.Data(http.StatusOK, "application/json", []byte(fullJSON.String))
}
