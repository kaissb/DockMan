// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package models

import "gorm.io/gorm"

// ComposeProject represents a docker-compose project managed by the application.
type ComposeProject struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
	Path string `json:"path"`
}
