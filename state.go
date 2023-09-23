package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aldehir/ut2u/pkg/query"
)

type State struct {
	servers      map[int]*ServerState
	serversMutex sync.RWMutex

	MaxPendingRemovals int
	removals           int
}

type ServerState struct {
	Registration Registration
	Online       bool
	Failures     int
	Details      *query.ServerDetails

	Created time.Time
	Updated time.Time
}

func (s ServerState) Filled() float64 {
	if s.Details == nil {
		return 0
	}

	if s.Details.Info.MaxPlayers == 0 {
		return 0
	}

	return float64(s.Details.Info.CurrentPlayers) / float64(s.Details.Info.MaxPlayers)
}

type ServerStateUpdate struct {
	Online  bool
	Details *query.ServerDetails
}

var ErrInvalidID = errors.New("invalid id")

func NewState() *State {
	return &State{
		servers:            make(map[int]*ServerState),
		MaxPendingRemovals: 100,
	}
}

func (s *State) Servers() []ServerState {
	s.serversMutex.RLock()
	defer s.serversMutex.RUnlock()

	servers := make([]ServerState, 0, len(s.servers))

	for _, server := range s.servers {
		servers = append(servers, *server)
	}

	return servers
}

func (s *State) Get(id int) (ServerState, bool) {
	s.serversMutex.RLock()
	defer s.serversMutex.RUnlock()

	if state, exists := s.servers[id]; exists {
		return *state, true
	}

	return ServerState{}, false
}

func (s *State) GetOrAdd(reg Registration) ServerState {
	s.serversMutex.Lock()
	defer s.serversMutex.Unlock()

	if state, exists := s.servers[reg.ID]; exists {
		return *state
	}

	t := time.Now()

	state := &ServerState{
		Registration: reg,
		Created:      t,
	}

	s.servers[reg.ID] = state
	return *state
}

func (s *State) Update(id int, update ServerStateUpdate) error {
	s.serversMutex.Lock()
	defer s.serversMutex.Unlock()

	state, exists := s.servers[id]
	if !exists {
		return fmt.Errorf("%w: %d", ErrInvalidID, id)
	}

	state.Online = update.Online
	state.Details = update.Details
	state.Updated = time.Now()

	if state.Online {
		state.Failures = 0
	} else {
		state.Failures++
	}

	return nil
}

func (s *State) Remove(id int) {
	s.serversMutex.Lock()
	defer s.serversMutex.Unlock()

	_, exists := s.servers[id]
	if !exists {
		return
	}

	delete(s.servers, id)
	s.removals++

	if s.removals >= s.MaxPendingRemovals {
		s.gc()
	}
}

func (s *State) gc() {
	logger.Debug("performing state garbage collection")

	// Go maps are never sized down, so after several removals we create a new
	// map and clean up the old
	old := s.servers
	s.servers = make(map[int]*ServerState, len(old))
	for k, v := range old {
		s.servers[k] = v
	}

	clear(old)
	s.removals = 0
}
