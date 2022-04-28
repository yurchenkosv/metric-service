package types

import "fmt"

type Gauge float64
type Counter int64

type URLServer struct {
	host   string
	port   string
	schema string
}

func (s *URLServer) Build() string {
	if s.schema != "" {
		return fmt.Sprintf("%s://%s:%s", s.schema, s.host, s.port)
	}
	return fmt.Sprintf("%s:%s", s.host, s.port)
}

func (s *URLServer) SetHost(host string) *URLServer {
	s.host = host
	return s
}

func (s *URLServer) SetPort(port string) *URLServer {
	s.port = port
	return s
}

func (s *URLServer) SetSchema(schema string) *URLServer {
	s.schema = schema
	return s
}

type MemMetrics struct {
	Alloc         Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	Mallocs       Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge
	Sys           Gauge
	TotalAlloc    Gauge
	PollCount     Counter
	RandomValue   Gauge
	GaugeMetrics  map[string]Gauge
}
