package models

import "time"

type Clinician struct {
    ID        int64     `json:"id"`
    FullName  string    `json:"full_name"`
    Email     string    `json:"email,omitempty"`
    Role      string    `json:"role,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
