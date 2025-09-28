// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package models

import "gorm.io/gorm"

// Service represents a deployable unit within an environment.
type Service struct {
	gorm.Model
	Name          string `json:"name"`
	EnvironmentID uint   `json:"environment_id"`

	// --- Service Hierarchy ---
	ParentServiceID *uint     `json:"parent_service_id,omitempty"` // Pointer to allow null
	SubServices     []Service `json:"sub_services,omitempty" gorm:"foreignKey:ParentServiceID"`

	// --- Service Type ---
	// 'container', 'compose', 'database', etc.
	Type string `json:"type"`

	// --- Type-Specific Fields ---
	// For 'container'
	Image       string `json:"image,omitempty"`
	ContainerID string `json:"container_id,omitempty"`

	// For 'compose'
	ComposePath string `json:"compose_path,omitempty"`

	// For 'database' (future use)
	// DBType string `json:"db_type,omitempty"` 
	// DBConnectionString string `json:"-"` // Don't expose connection strings

	// --- Git Integration for CI/CD (Future Feature) ---
	GitRepoURL  string `json:"git_repo_url,omitempty"`
	GitBranch   string `json:"git_branch,omitempty"`
	WebhookID   string `json:"-"` // Don't expose webhook secret
}
