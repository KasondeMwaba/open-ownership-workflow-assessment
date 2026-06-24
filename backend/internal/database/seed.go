package database

import (
	"encoding/json"

	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/workflow"

	"gorm.io/gorm"
)

const demoPasswordHash = "$2a$10$BPMCXMVTIW07SipA06cqpuNmdOZhDNbS6LjtMAlwdxHi1d2GuZU5O"

func SeedDemoData(db *gorm.DB) error {
	if err := seedAccessControls(db); err != nil {
		return err
	}

	users := []models.User{
		{Name: "Amina Requester", Email: "requester@example.com", PasswordHash: demoPasswordHash, Role: workflow.Requester, IsActive: true},
		{Name: "Noah Reviewer", Email: "reviewer@example.com", PasswordHash: demoPasswordHash, Role: workflow.Reviewer, IsActive: true},
		{Name: "Sam Admin", Email: "admin@example.com", PasswordHash: demoPasswordHash, Role: workflow.Admin, IsActive: true},
	}
	for _, user := range users {
		if err := db.Where("email = ?", user.Email).Assign(models.User{IsActive: true}).FirstOrCreate(&user).Error; err != nil {
			return err
		}
	}

	var submissionCount int64
	if err := db.Model(&models.Submission{}).Count(&submissionCount).Error; err != nil {
		return err
	}
	if submissionCount > 0 {
		return nil
	}

	var requester models.User
	if err := db.Where("email = ?", "requester@example.com").First(&requester).Error; err != nil {
		return err
	}
	var reviewer models.User
	if err := db.Where("email = ?", "reviewer@example.com").First(&reviewer).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		submission := models.Submission{
			Title:      "Beneficial ownership declaration",
			Summary:    "Company declaration awaiting review with ownership and control details.",
			Data:       json.RawMessage(`{"company":"Northstar Minerals Ltd","jurisdiction":"Ghana","beneficialOwners":[{"name":"Amina Mensah","ownershipPercent":42.5},{"name":"Jon Bell","ownershipPercent":18}]}`),
			Status:     workflow.Submitted,
			OwnerID:    requester.ID,
			ReviewerID: &reviewer.ID,
			Version:    1,
		}
		if err := tx.Create(&submission).Error; err != nil {
			return err
		}

		event := models.AuditEvent{
			SubmissionID: submission.ID,
			ActorID:      requester.ID,
			ToStatus:     submission.Status,
			Comment:      "Seed submission created",
			Metadata:     json.RawMessage(`{}`),
		}
		return tx.Create(&event).Error
	})
}

func seedAccessControls(db *gorm.DB) error {
	permissions := []models.Permission{
		{Name: "users:manage", Description: "Create, update, enable, and disable user accounts"},
		{Name: "roles:manage", Description: "Create roles and assign permissions"},
		{Name: "permissions:manage", Description: "Create and update permission records"},
		{Name: "submissions:create", Description: "Create and edit owned draft submissions"},
		{Name: "submissions:review", Description: "Review submitted records and approve or reject them"},
		{Name: "dashboard:view", Description: "View dashboard summaries and workflow metrics"},
	}
	for _, permission := range permissions {
		if err := db.Where("name = ?", permission.Name).FirstOrCreate(&permission).Error; err != nil {
			return err
		}
	}

	rolePermissions := map[string][]string{
		string(workflow.Admin):     {"users:manage", "roles:manage", "permissions:manage", "submissions:create", "submissions:review", "dashboard:view"},
		string(workflow.Reviewer):  {"submissions:review", "dashboard:view"},
		string(workflow.Requester): {"submissions:create", "dashboard:view"},
	}
	descriptions := map[string]string{
		string(workflow.Admin):     "System administrator with full access",
		string(workflow.Reviewer):  "Reviewer who can process submitted records",
		string(workflow.Requester): "Requester who can create and manage own submissions",
	}

	for roleName, names := range rolePermissions {
		role := models.AccessRole{Name: roleName, Description: descriptions[roleName]}
		if err := db.Where("name = ?", roleName).FirstOrCreate(&role).Error; err != nil {
			return err
		}
		var permissions []models.Permission
		if err := db.Where("name IN ?", names).Find(&permissions).Error; err != nil {
			return err
		}
		if err := db.Model(&role).Association("Permissions").Replace(permissions); err != nil {
			return err
		}
	}
	return nil
}
