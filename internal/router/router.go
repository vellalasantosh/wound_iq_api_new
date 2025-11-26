package router

import (
    "database/sql"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "github.com/vellalasantosh/wound_iq_api/internal/config"
    "github.com/vellalasantosh/wound_iq_api/internal/handlers"
)

func New(db *sql.DB, log *zap.Logger, cfg *config.Config) *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(gin.LoggerWithWriter(gin.DefaultWriter))
    r.Use(corsMiddleware())

    h := handlers.NewHandlers(db, log, cfg)

    v1 := r.Group("/v1")
    {
        // Patients
        v1.GET("/patients", h.ListPatients)
        v1.GET("/patients/:id", h.GetPatient)
        v1.POST("/patients", h.CreatePatient)
        v1.PUT("/patients/:id", h.UpdatePatient)
        v1.DELETE("/patients/:id", h.DeletePatient)

        // Clinicians
        v1.GET("/clinicians", h.ListClinicians)
        v1.GET("/clinicians/:id", h.GetClinician)
        v1.POST("/clinicians", h.CreateClinician)
        v1.PUT("/clinicians/:id", h.UpdateClinician)
        v1.DELETE("/clinicians/:id", h.DeleteClinician)

        // Assessments
        v1.GET("/assessments", h.ListAssessments)
        v1.GET("/assessments/:id", h.GetAssessment)
        v1.POST("/assessments", h.CreateAssessment)
        v1.PUT("/assessments/:id", h.UpdateAssessment)
        v1.DELETE("/assessments/:id", h.DeleteAssessment)

        // Reports
        v1.GET("/patients/:id/history", h.GetPatientHistory)
        v1.GET("/assessments/:id/full", h.GetAssessmentFull)
    }

    return r
}

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}
