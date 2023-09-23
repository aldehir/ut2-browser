package main

import (
	"sync"
	"time"
)

type Registration struct {
	ID           int
	Group        string
	Address      string
	Interval     time.Duration
	Timeout      time.Duration
	Persist      bool
	ExternalLink string
}

type Registry struct {
	registrations      []Registration
	registrationsMutex sync.RWMutex

	nextID int
}

func (r *Registry) Register(reg Registration) {
	r.registrationsMutex.Lock()
	defer r.registrationsMutex.Unlock()

	// Update registration if one already exists for the given address
	for i := 0; i < len(r.registrations); i++ {
		if r.registrations[i].Address == reg.Address {
			id := r.registrations[i].ID
			r.registrations[i] = reg
			r.registrations[i].ID = id

			logger.Debug("update registration", "id", r.registrations[i].ID, "addr", r.registrations[i].Address)
			return
		}
	}

	// Assign an ID for tracking state
	reg.ID = r.nextID
	r.nextID++
	r.registrations = append(r.registrations, reg)

	logger.Debug("register server", "id", reg.ID, "addr", reg.Address)
}

func (r *Registry) Unregister(id int) {
	r.registrationsMutex.Lock()
	defer r.registrationsMutex.Unlock()

	for i := len(r.registrations) - 1; i >= 0; i-- {
		if r.registrations[i].ID == id {
			logger.Debug("remove registration", "id", r.registrations[i].ID, "addr", r.registrations[i].Address)
			r.registrations[i] = r.registrations[len(r.registrations)-1]
			r.registrations = r.registrations[:len(r.registrations)-1]
		}
	}
}

func (r *Registry) Registrations() []Registration {
	r.registrationsMutex.RLock()
	defer r.registrationsMutex.RUnlock()

	// Create a copy to ensure that the data does not get modified underneath
	result := make([]Registration, len(r.registrations))
	copy(result, r.registrations)
	return result
}
