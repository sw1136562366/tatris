// Copyright 2022 Tatris Project Authors. Licensed under Apache-2.0.

// Package utils contains basic utilities for Tatris
package utils

import (
	"crypto/rand"
	"fmt"

	"github.com/tatris-io/tatris/internal/common/log/logger"
)

// TODO: distributed ID
func GenerateID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		logger.Errorf("generate ID fail: %s", err.Error())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
