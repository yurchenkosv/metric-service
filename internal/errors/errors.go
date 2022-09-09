// Package errors for custom application errors.
package errors

import "fmt"

type NoSuchMetricError struct {
	MetricName string
}

type BatchProcessingError struct {
}

type HealthCheckError struct {
	HealthcheckType string
}

type MetricNotFoundError struct {
	MetricName string
}

type NoEncryptionKeyFoundError struct {
}

func (e NoSuchMetricError) Error() string {
	return fmt.Sprintf("no such metric %s", e.MetricName)
}

func (e BatchProcessingError) Error() string {
	return "batch process error"
}

func (e HealthCheckError) Error() string {
	return fmt.Sprintf("%s healthcheck failed", e.HealthcheckType)
}

func (e MetricNotFoundError) Error() string {
	return fmt.Sprintf("metric with name %s not found", e.MetricName)
}

func (e NoEncryptionKeyFoundError) Error() string {
	return "encryption key for hash signing not found"
}
