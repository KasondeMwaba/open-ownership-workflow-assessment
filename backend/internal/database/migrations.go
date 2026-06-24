package database

import (
	"openownership-workflow/backend/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&models.Permission{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&models.AccessRole{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&models.User{}); err != nil {
			return err
		}
		if tx.Migrator().HasConstraint(&models.User{}, "chk_users_role") {
			if err := tx.Migrator().DropConstraint(&models.User{}, "chk_users_role"); err != nil {
				return err
			}
		}
		if err := tx.AutoMigrate(&models.Submission{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&models.AuditEvent{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&models.SystemAuditEvent{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&models.SessionAuditEvent{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&models.ActivityAuditEvent{}); err != nil {
			return err
		}
		return nil
	})
}
