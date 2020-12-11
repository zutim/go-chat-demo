package utils

import uuid "github.com/satori/go.uuid"

// UUID return unique id
func UUID() string {
	return uuid.NewV4().String()
}
