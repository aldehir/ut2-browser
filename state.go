package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aldehir/ut2u/pkg/query"
)

type State struct {
	Servers      map[int]*ServerState
	serversMutex sync.RWMutex

	MaxPendingRemovals int
	removals           int
}

type ServerState struct {
	Address  string
	Online   bool
	Failures int
	Details  *query.ServerDetails

	Created time.Time
	Updated time.Time
}

type ServerStateUpdate struct {
	Online  bool
	Details *query.ServerDetails
}

var ErrInvalidID = errors.New("invalid id")

func NewState() *State {
	return &State{
		Servers:            make(map[int]*ServerState),
		MaxPendingRemovals: 100,
	}
}

func (s *State) Get(id int) (ServerState, bool) {
	s.serversMutex.RLock()
	defer s.serversMutex.RUnlock()

	if state, exists := s.Servers[id]; exists {
		return *state, true
	}

	return ServerState{}, false
}

func (s *State) GetOrAdd(id int, address string) ServerState {
	s.serversMutex.Lock()
	defer s.serversMutex.Unlock()

	if state, exists := s.Servers[id]; exists {
		return *state
	}

	t := time.Now()

	state := &ServerState{
		Address: address,
		Created: t,
	}

	s.Servers[id] = state
	return *state
}

func (s *State) Update(id int, update ServerStateUpdate) error {
	s.serversMutex.Lock()
	defer s.serversMutex.Unlock()

	state, exists := s.Servers[id]
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

	_, exists := s.Servers[id]
	if !exists {
		return
	}

	delete(s.Servers, id)
	s.removals++

	if s.removals >= s.MaxPendingRemovals {
		s.gc()
	}
}

func (s *State) gc() {
	logger.Debug("performing state garbage collection")

	// Go maps are never sized down, so after several removals we create a new
	// map and clean up the old
	old := s.Servers
	s.Servers = make(map[int]*ServerState, len(old))
	for k, v := range old {
		s.Servers[k] = v
	}

	clear(old)
	s.removals = 0
}
