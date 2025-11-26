package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vellalasantosh/wound_iq_api_new/internal/models"
)

func parsePagination(c *gin.Context) (int, int) {
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}
	return page, pageSize
}

// ListPatients GET /v1/patients
func (h *Handlers) ListPatients(c *gin.Context) {
	page, pageSize := parsePagination(c)
	offset := (page - 1) * pageSize

	rows, err := h.DB.Query(`SELECT id, full_name, date_of_birth, gender, medical_record_number, created_at, updated_at
                             FROM patients ORDER BY id DESC LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		h.Log.Sugar().Errorf("list patients: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch patients"})
		return
	}
	defer rows.Close()

	patients := []models.Patient{}
	for rows.Next() {
		var p models.Patient
		var dob sql.NullTime
		if err := rows.Scan(&p.ID, &p.FullName, &dob, &p.Gender, &p.MedicalRecordNumber, &p.CreatedAt, &p.UpdatedAt); err != nil {
			h.Log.Sugar().Errorf("scan patient: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read patients"})
			return
		}
		if dob.Valid {
			t := dob.Time
			p.DateOfBirth = &t
		}
		patients = append(patients, p)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      patients,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetPatient GET /v1/patients/:id
func (h *Handlers) GetPatient(c *gin.Context) {
	id := c.Param("id")
	var p models.Patient
	var dob sql.NullTime
	row := h.DB.QueryRow(`SELECT id, full_name, date_of_birth, gender, medical_record_number, created_at, updated_at 
                          FROM patients WHERE id=$1`, id)
	if err := row.Scan(&p.ID, &p.FullName, &dob, &p.Gender, &p.MedicalRecordNumber, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
			return
		}
		h.Log.Sugar().Errorf("get patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get patient"})
		return
	}
	if dob.Valid {
		t := dob.Time
		p.DateOfBirth = &t
	}
	c.JSON(http.StatusOK, p)
}

// CreatePatient POST /v1/patients
func (h *Handlers) CreatePatient(c *gin.Context) {
	var in struct {
		FullName            string `json:"full_name" binding:"required"`
		DateOfBirth         string `json:"date_of_birth"` // ISO-8601 expected
		Gender              string `json:"gender"`
		MedicalRecordNumber string `json:"medical_record_number"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dobParam interface{}
	if in.DateOfBirth != "" {
		t, err := time.Parse(time.RFC3339, in.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "date_of_birth must be ISO-8601 (RFC3339)"})
			return
		}
		dobParam = t
	} else {
		dobParam = nil
	}

	var newID int64
	err := h.DB.QueryRow(`SELECT add_patient($1, $2, $3, $4)`, in.FullName, dobParam, in.Gender, in.MedicalRecordNumber).Scan(&newID)
	if err != nil {
		h.Log.Sugar().Errorf("call add_patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create patient"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": newID})
}

// UpdatePatient PUT /v1/patients/:id
func (h *Handlers) UpdatePatient(c *gin.Context) {
	id := c.Param("id")
	var in struct {
		FullName            *string `json:"full_name"`
		DateOfBirth         *string `json:"date_of_birth"`
		Gender              *string `json:"gender"`
		MedicalRecordNumber *string `json:"medical_record_number"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.DB.Exec(`UPDATE patients SET full_name = COALESCE($1, full_name),
                       date_of_birth = COALESCE($2, date_of_birth),
                       gender = COALESCE($3, gender),
                       medical_record_number = COALESCE($4, medical_record_number),
                       updated_at = now()
                       WHERE id = $5`,
		in.FullName, nilIfEmptyPtr(in.DateOfBirth), in.Gender, in.MedicalRecordNumber, id)
	if err != nil {
		h.Log.Sugar().Errorf("update patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update patient"})
		return
	}
	c.Status(http.StatusNoContent)
}

// DeletePatient DELETE /v1/patients/:id
func (h *Handlers) DeletePatient(c *gin.Context) {
	id := c.Param("id")
	res, err := h.DB.Exec(`DELETE FROM patients WHERE id = $1`, id)
	if err != nil {
		h.Log.Sugar().Errorf("delete patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete patient"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func nilIfEmptyPtr(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil
	}
	return t
}
