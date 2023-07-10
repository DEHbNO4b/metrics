package data

type Gauge struct {
	Name string
	Val  float64
}
type Counter struct {
	Name string
	Val  int
}

type MetStore struct {
	Gauges   map[string]float64
	Counters map[string]int
}

func NewMetStore() *MetStore {
	g := make(map[string]float64)
	c := make(map[string]int)
	ms := MetStore{Gauges: g, Counters: c}
	return &ms
}

func (ms *MetStore) SetGauge(g Gauge) error {
	ms.Gauges[g.Name] = g.Val
	return nil
}
func (ms *MetStore) SetCounter(c Counter) error {
	ms.Counters[c.Name] = c.Val
	return nil
}
func (ms *MetStore) GetGauge(name string) (Gauge, error) {
	g := Gauge{}
	g.Name = name
	g.Val = ms.Gauges[name]
	return g, nil
}
func (ms *MetStore) GetCounter(name string) (Counter, error) {
	c := Counter{}
	c.Name = name
	c.Val = ms.Counters[name]
	return c, nil
}
