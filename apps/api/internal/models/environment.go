// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package models

import "gorm.io/gorm"

// Environment represents a deployment environment within a project (e.g., dev, staging, prod).
type Environment struct {
	gorm.Model
	Name      string    `json:"name"`
	ProjectID uint      `json:"project_id"`
	Services  []Service `json:"services,omitempty"`
	Variables []EnvironmentVariable `json:"variables,omitempty"`
}
