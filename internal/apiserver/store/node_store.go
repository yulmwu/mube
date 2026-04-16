package store

import (
	"fmt"
	"sort"
	"sync"
	"time"

	v1 "mube/pkg/api/v1"
)

type RegisteredNode struct {
	Name string
	IP   string
	Port int
}

type NodeStore interface {
	UpsertHeartbeat(req v1.NodeHeartbeatRequest, now time.Time) (v1.NodeStatus, error)
	List(now time.Time) []v1.NodeStatus
}

type MemoryNodeStore struct {
	mu              sync.RWMutex
	registered      map[string]RegisteredNode
	registeredAt    map[string]time.Time
	nodes           map[string]v1.NodeStatus
	notReadyTimeout time.Duration
}

func NewMemoryNodeStore(registered []RegisteredNode, notReadyTimeout time.Duration, now time.Time) *MemoryNodeStore {
	now = now.UTC()
	registeredMap := make(map[string]RegisteredNode, len(registered))
	registeredAtMap := make(map[string]time.Time, len(registered))
	nodes := make(map[string]v1.NodeStatus, len(registered))
	for _, n := range registered {
		registeredMap[n.Name] = n
		registeredAtMap[n.Name] = now
		nodes[n.Name] = v1.NodeStatus{
			Name:  n.Name,
			IP:    n.IP,
			Port:  n.Port,
			State: v1.NodeStateUnknown,
		}
	}

	return &MemoryNodeStore{
		registered:      registeredMap,
		registeredAt:    registeredAtMap,
		nodes:           nodes,
		notReadyTimeout: notReadyTimeout,
	}
}

func (s *MemoryNodeStore) UpsertHeartbeat(req v1.NodeHeartbeatRequest, now time.Time) (v1.NodeStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	registered, ok := s.registered[req.Name]
	if !ok {
		return v1.NodeStatus{}, fmt.Errorf("node %q is not registered in apiserver config", req.Name)
	}

	node := v1.NodeStatus{
		Name:          req.Name,
		IP:            registered.IP,
		Port:          registered.Port,
		State:         v1.NodeStateReady,
		Runtime:       req.Runtime,
		Version:       req.Version,
		Capacity:      req.Capacity,
		LastHeartbeat: now.UTC(),
	}
	s.nodes[req.Name] = node

	return node, nil
}

func (s *MemoryNodeStore) List(now time.Time) []v1.NodeStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now = now.UTC()
	items := make([]v1.NodeStatus, 0, len(s.nodes))
	for _, n := range s.nodes {
		n.State = s.computeState(n, now)
		items = append(items, n)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items
}

func (s *MemoryNodeStore) computeState(node v1.NodeStatus, now time.Time) string {
	if node.LastHeartbeat.IsZero() {
		registeredAt, ok := s.registeredAt[node.Name]
		if !ok {
			return v1.NodeStateUnknown
		}

		if now.Sub(registeredAt) <= s.notReadyTimeout {
			return v1.NodeStateUnknown
		}

		return v1.NodeStateNotReady
	}

	if now.Sub(node.LastHeartbeat) <= s.notReadyTimeout {
		return v1.NodeStateReady
	}

	return v1.NodeStateNotReady
}
