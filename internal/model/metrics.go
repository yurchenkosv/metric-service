package model

import (
	"encoding/json"
	"fmt"
)

type Gauge float64 // custom type for Gauge based on float64
type Counter int64 // custom type for Counter based on int64

// Metrics struct represents set of Metric
type Metrics struct {
	Metric []Metric // slice of Metric
}

// Metric - struc that represents single metric
type Metric struct {
	ID    string   `json:"id"`              // metric name
	MType string   `json:"type"`            // parameter, that could be gauge or counter
	Delta *Counter `json:"delta,omitempty"` // metric value in case of Counter
	Value *Gauge   `json:"value,omitempty"` // metric value in case of Gauge
	Hash  string   `json:"hash,omitempty"`  // hash-function value
}

// String is string representation of Gauge
func (g *Gauge) String() string {
	if g == nil {
		return ""
	}
	return fmt.Sprintf("%.3f", *g)
}

// String is string representation of Counter
func (c *Counter) String() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("%d", *c)
}

// String method for string representation of Metrics.
// Useful when we generate page with all metrics
func (m Metrics) String() string {
	var result string
	for _, v := range m.Metric {
		switch v.MType {
		case "counter":
			result = result + fmt.Sprintf("%s = %s\n", v.ID, v.Delta.String())
		case "gauge":
			result = result + fmt.Sprintf("%s = %s\n", v.ID, v.Value.String())
		}
	}
	return result
}

// NewCounter is function - helper to convert int64 to *Counter
func NewCounter(val int64) *Counter {
	cval := Counter(val)
	return &cval
}

// NewGauge is function - helper to convert float64 to *Gauge
func NewGauge(val float64) *Gauge {
	gval := Gauge(val)
	return &gval
}

// MarshalJSON custom marshaller for  Counter type
func (c *Counter) MarshalJSON() ([]byte, error) {
	str := c.String()
	return json.Marshal(str)
}

// MarshalJSON custom marshaller for Gauge type.
func (g *Gauge) MarshalJSON() ([]byte, error) {
	str := g.String()
	return json.Marshal(str)
}

// MarshalJSON custom unmarshaller for Metric struct.
// We need this because of custom types, Counter and Gauge inside Metric.
func (m *Metric) MarshalJSON() ([]byte, error) {
	type Alias Metric
	var (
		value *float64
		delta *int64
	)
	if m.Value != nil {
		val := float64(*m.Value)
		value = &val
		delta = nil
	} else if m.Delta != nil {
		val := int64(*m.Delta)
		delta = &val
		value = nil
	}
	return json.Marshal(&struct {
		Delta *int64   `json:"delta,omitempty"`
		Value *float64 `json:"value,omitempty"`
		*Alias
	}{
		Delta: delta,
		Value: value,
		Alias: (*Alias)(m),
	})
}
