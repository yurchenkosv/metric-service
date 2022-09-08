package model

import (
	"encoding/json"
	"fmt"
)

type Gauge float64
type Counter int64

type Metrics struct {
	Metric []Metric
}

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func (g *Gauge) String() string {
	if g == nil {
		return ""
	}
	return fmt.Sprintf("%.6f", *g)
}

func (c *Counter) String() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("%d", *c)
}

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

func NewCounter(val int64) *Counter {
	cval := Counter(val)
	return &cval
}

func NewGauge(val float64) *Gauge {
	gval := Gauge(val)
	return &gval
}

func (c *Counter) MarshalJSON() ([]byte, error) {
	str := c.String()
	return json.Marshal(str)
}

func (g *Gauge) MarshalJSON() ([]byte, error) {
	str := g.String()
	return json.Marshal(str)
}

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
