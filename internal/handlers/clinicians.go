package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vellalasantosh/wound_iq_api_new/internal/models"
)

// ListClinicians GET /v1/clinicians
func (h *Handlers) ListClinicians(c *gin.Context) {
	page, pageSize := parsePagination(c)
	offset := (page - 1) * pageSize

	rows, err := h.DB.Query(`SELECT id, full_name, email, role, created_at, updated_at
                             FROM clinicians ORDER BY id DESC LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		h.Log.Sugar().Errorf("list clinicians: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch clinicians"})
		return
	}
	defer rows.Close()

	out := []models.Clinician{}
	for rows.Next() {
		var cl models.Clinician
		if err := rows.Scan(&cl.ID, &cl.FullName, &cl.Email, &cl.Role, &cl.CreatedAt, &cl.UpdatedAt); err != nil {
			h.Log.Sugar().Errorf("scan clinician: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read clinicians"})
			return
		}
		out = append(out, cl)
	}
	c.JSON(http.StatusOK, gin.H{"data": out, "page": page, "page_size": pageSize})
}

// GetClinician GET /v1/clinicians/:id
func (h *Handlers) GetClinician(c *gin.Context) {
	id := c.Param("id")
	var cl models.Clinician
	row := h.DB.QueryRow(`SELECT id, full_name, email, role, created_at, updated_at FROM clinicians WHERE id=$1`, id)
	if err := row.Scan(&cl.ID, &cl.FullName, &cl.Email, &cl.Role, &cl.CreatedAt, &cl.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "clinician not found"})
			return
		}
		h.Log.Sugar().Errorf("get clinician: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get clinician"})
		return
	}
	c.JSON(http.StatusOK, cl)
}

// CreateClinician POST /v1/clinicians
func (h *Handlers) CreateClinician(c *gin.Context) {
	var in struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newID int64
	err := h.DB.QueryRow(`INSERT INTO clinicians (full_name, email, role, created_at, updated_at)
                          VALUES ($1, $2, $3, now(), now()) RETURNING id`, in.FullName, in.Email, in.Role).Scan(&newID)
	if err != nil {
		h.Log.Sugar().Errorf("create clinician: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create clinician"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": newID})
}

// UpdateClinician PUT /v1/clinicians/:id
func (h *Handlers) UpdateClinician(c *gin.Context) {
	id := c.Param("id")
	var in struct {
		FullName *string `json:"full_name"`
		Email    *string `json:"email"`
		Role     *string `json:"role"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.DB.Exec(`UPDATE clinicians SET full_name = COALESCE($1, full_name),
                       email = COALESCE($2, email),
                       role = COALESCE($3, role),
                       updated_at = now()
                       WHERE id = $4`,
		in.FullName, in.Email, in.Role, id)
	if err != nil {
		h.Log.Sugar().Errorf("update clinician: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update clinician"})
		return
	}
	c.Status(http.StatusNoContent)
}

// DeleteClinician DELETE /v1/clinicians/:id
func (h *Handlers) DeleteClinician(c *gin.Context) {
	id := c.Param("id")
	res, err := h.DB.Exec(`DELETE FROM clinicians WHERE id = $1`, id)
	if err != nil {
		h.Log.Sugar().Errorf("delete clinician: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete clinician"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "clinician not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
