package models

import "time"

type Patient struct {
    ID                  int64      `json:"id"`
    FullName            string     `json:"full_name"`
    DateOfBirth         *time.Time `json:"date_of_birth,omitempty"`
    Gender              string     `json:"gender,omitempty"`
    MedicalRecordNumber string     `json:"medical_record_number,omitempty"`
    CreatedAt           time.Time  `json:"created_at"`
    UpdatedAt           time.Time  `json:"updated_at"`
}
