package models

import "time"

type Assessment struct {
    ID          int64     `json:"id"`
    PatientID   int64     `json:"patient_id"`
    ClinicianID int64     `json:"clinician_id"`
    WoundID     *int64    `json:"wound_id,omitempty"`
    Notes       string    `json:"notes,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
