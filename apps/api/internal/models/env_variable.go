// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package models

import (
	"docker-manager/api/internal/crypto"
	"gorm.io/gorm"
)

// EnvironmentVariable represents a key-value pair for an environment.
// The value is encrypted at rest in the database.
type EnvironmentVariable struct {
	gorm.Model
	Key           string `json:"key" gorm:"uniqueIndex:idx_env_key"`
	Value         string `json:"value"`
	EnvironmentID uint   `json:"environment_id" gorm:"uniqueIndex:idx_env_key"`
}

// BeforeSave is a GORM hook that encrypts the Value before saving it to the database.
func (ev *EnvironmentVariable) BeforeSave(tx *gorm.DB) (err error) {
	encryptedValue, err := crypto.Encrypt(ev.Value)
	if err != nil {
		return err
	}
	ev.Value = encryptedValue
	return nil
}

// AfterFind is a GORM hook that decrypts the Value after retrieving it from the database.
func (ev *EnvironmentVariable) AfterFind(tx *gorm.DB) (err error) {
	decryptedValue, err := crypto.Decrypt(ev.Value)
	if err != nil {
		// If decryption fails, it might be because the value is not encrypted (e.g. during creation)
		// We can choose to ignore the error or handle it. For now, we'll just return the encrypted value.
		return nil 
	}
	ev.Value = decryptedValue
	return nil
}
