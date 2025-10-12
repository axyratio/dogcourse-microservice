package config

import (
	"log"

	"gorm.io/gorm"
	"auth-service/internal/models"
)

func SeedRoles(db *gorm.DB) {
	roles := []models.Role{
		{ID: 1, RoleName: "user"},
		{ID: 2, RoleName: "trainer"},
		{ID: 3, RoleName: "admin"},
	}

	for _, role := range roles {
		var existing models.Role
		if err := db.First(&existing, role.ID).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to insert role %s: %v", role.RoleName, err)
			} else {
				log.Printf("Role %s inserted", role.RoleName)
			}
		}
	}
}