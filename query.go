package main

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/aldehir/ut2u/pkg/query"
)

type QueryEngine struct {
	Registry    *Registry
	State       *State
	Concurrency int
	Context     context.Context
	MaxFailures int

	client *query.Client

	stop     chan struct{}
	stopOnce sync.Once

	pending      []int
	pendingMutex sync.RWMutex
}

var ErrStopped = errors.New("stopped")

func (m *QueryEngine) init() (err error) {
	if m.Registry == nil {
		m.Registry = &Registry{}
	}

	if m.State == nil {
		m.State = &State{}
	}

	if m.Context == nil {
		m.Context = context.Background()
	}

	if m.Concurrency == 0 {
		m.Concurrency = 12
	}

	if m.MaxFailures == 0 {
		m.MaxFailures = 10
	}

	m.stop = make(chan struct{}, 1)

	m.client, err = query.NewClient()
	return
}

func (m *QueryEngine) Run() error {
	err := m.init()
	if err != nil {
		return err
	}

	m.client, err = query.NewClient()
	if err != nil {
		return err
	}

	var jobTokens = make(chan struct{}, m.Concurrency)

	ticker := time.NewTicker(10 * time.Millisecond)

	var runErr error

	for runErr == nil {
		select {
		case <-m.stop:
			runErr = ErrStopped
		case <-m.Context.Done():
			runErr = m.Context.Err()
		case <-ticker.C:
			// Find servers that are ready for query
			regs := m.Registry.Registrations()
			for _, r := range regs {
				state := m.State.GetOrAdd(r)

				if time.Since(state.Updated) >= m.randomize(r.Interval) && !m.isPending(r.ID) {
					m.markPending(r.ID)

					go func(r Registration, state ServerState) {
						jobTokens <- struct{}{}
						defer func() { <-jobTokens }()

						m.performQuery(r, state)
						m.clearPending(r.ID)
					}(r, state)
				}
			}
		}
	}

	ticker.Stop()
	close(m.stop)

	return runErr
}

// randomize adds a random factor to given time duration to prevent all queries
// from happening at the same time
func (m *QueryEngine) randomize(t time.Duration) time.Duration {
	return t + (time.Duration(rand.Intn(2000)) * time.Millisecond)
}

func (m *QueryEngine) performQuery(r Registration, state ServerState) {
	addr, err := net.ResolveUDPAddr("udp", r.Address)
	if err != nil {
		logger.Error("failed to resolve address", "addr", r.Address, "err", err)
		return
	}

	// Use query port
	addr.Port += 1

	opts := []query.QueryOption{
		query.WithRules(),
		query.WithPlayers(),
		query.WithTimeout(r.Timeout),
	}

	var offline bool

	details, err := m.client.Query(m.Context, addr, opts...)

	if err != nil {
		offline = true
		logger.Error("failed to query server", "addr", r.Address, "err", err)

		// If not persisting and we reach max failures, remove registration
		if !r.Persist && state.Failures+1 >= m.MaxFailures {
			m.Registry.Unregister(r.ID)
			m.State.Remove(r.ID)
		}
	} else {
		logger.Debug("query success", "addr", r.Address, "servername", details.Info.ServerName, "mapname", details.Info.MapName)
	}

	m.State.Update(r.ID, ServerStateUpdate{
		Online:  !offline,
		Details: &details,
	})
}

func (m *QueryEngine) markPending(id int) {
	m.pendingMutex.Lock()
	defer m.pendingMutex.Unlock()

	for _, other := range m.pending {
		if other == id {
			return
		}
	}

	m.pending = append(m.pending, id)
}

func (m *QueryEngine) isPending(id int) bool {
	m.pendingMutex.RLock()
	defer m.pendingMutex.RUnlock()

	for _, other := range m.pending {
		if other == id {
			return true
		}
	}

	return false
}

func (m *QueryEngine) clearPending(id int) {
	m.pendingMutex.Lock()
	defer m.pendingMutex.Unlock()

	for i := len(m.pending) - 1; i >= 0; i-- {
		if m.pending[i] == id {
			m.pending[i] = m.pending[len(m.pending)-1]
			m.pending = m.pending[:len(m.pending)-1]
		}
	}
}

func (m *QueryEngine) Stop() {
	m.stopOnce.Do(func() {
		m.stop <- struct{}{}
	})
}
