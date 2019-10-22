package health

import "encoding/json"

type status string

const (
	up           status = "UP"
	down                = "DOWN"
	outOfService        = "OUT OF SERVICE"
	unknown             = "UNKNOWN"
)

// Health is a health status interface
type Health interface {
	// MarshalJSON is a custom JSON marshaller
	MarshalJSON() ([]byte, error)
	// AddInfo adds a info value to the Info map
	AddInfo(key string, value interface{}) Health
	// GetInfo returns a value from the info map
	GetInfo(key string) interface{}
	// IsUnknown returns true if Status is Unknown
	IsUnknown() bool
	// IsUp returns true if Status is Up
	IsUp() bool
	// IsDown returns true if Status is Down
	IsDown() bool
	// IsOutOfService returns true if Status is IsOutOfService
	IsOutOfService() bool
	// Down set the status to Down
	Down() Health
	// OutOfService set the status to OutOfService
	OutOfService() Health
	// Unknown set the status to Unknown
	Unknown() Health
	// Up set the status to Up
	Up() Health
}

// Health is a health status struct implementation
type health struct {
	status status
	info   map[string]interface{}
}

// MarshalJSON is a custom JSON marshaller
func (h health) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	for k, v := range h.info {
		data[k] = v
	}

	data["status"] = h.status

	return json.Marshal(data)
}

// NewHealth return a new Health with status Down
func NewHealth() Health {
	h := health{
		info: make(map[string]interface{}),
	}

	h.Unknown()
	return &h
}

// AddInfo adds a info value to the Info map
func (h *health) AddInfo(key string, value interface{}) Health {
	if h.info == nil {
		h.info = make(map[string]interface{})
	}

	h.info[key] = value

	return h
}

// GetInfo returns a value from the info map
func (h health) GetInfo(key string) interface{} {
	return h.info[key]
}

// IsUnknown returns true if Status is Unknown
func (h *health) IsUnknown() bool {
	return h.status == unknown
}

// IsUp returns true if Status is Up
func (h health) IsUp() bool {
	return h.status == up
}

// IsDown returns true if Status is Down
func (h health) IsDown() bool {
	return h.status == down
}

// IsOutOfService returns true if Status is IsOutOfService
func (h health) IsOutOfService() bool {
	return h.status == outOfService
}

// Down set the status to Down
func (h *health) Down() Health {
	h.status = down
	return h
}

// OutOfService set the status to OutOfService
func (h *health) OutOfService() Health {
	h.status = outOfService
	return h
}

// Unknown set the status to Unknown
func (h *health) Unknown() Health {
	h.status = unknown
	return h
}

// Up set the status to Up
func (h *health) Up() Health {
	h.status = up
	return h
}
