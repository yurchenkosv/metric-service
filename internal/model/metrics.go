package model

import "fmt"

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

func (g Gauge) String() string {
	return fmt.Sprintf("%.2f", g)
}

func (c Counter) String() string {
	return fmt.Sprintf("%d", c)
}

func (m Metrics) String() string {
	var result string
	for _, v := range m.Metric {
		switch v.MType {
		case "counter":
			result = result + fmt.Sprintf("%s = %d\n", v.ID, v.Delta)
		case "gauge":
			result = result + fmt.Sprintf("%s = %.2f\n", v.ID, v.Value)
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
