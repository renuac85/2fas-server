package pairing

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MemoryStore keeps in memory pairing between extension and mobile.
//
// TODO: move cache to separate package.
// check ttlcache pkg.
type MemoryStore struct {
	mu            sync.Mutex
	extensionsMap map[string]Item
}

type Item struct {
	ExtensionID string
	Expires     time.Time
	PairingInfo PairingInfo
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		extensionsMap: make(map[string]Item),
	}
}

func (s *MemoryStore) AddExtension(_ context.Context, extensionID string) {
	s.setItem(extensionID, Item{ExtensionID: extensionID})
}

func (s *MemoryStore) ExtensionExists(_ context.Context, extensionID string) bool {
	_, ok := s.getItem(extensionID)
	return ok
}

func (s *MemoryStore) GetPairingInfo(ctx context.Context, extensionID string) (PairingInfo, error) {
	v, ok := s.getItem(extensionID)
	if !ok {
		return PairingInfo{}, errors.New("extension does not exists")
	}
	return v.PairingInfo, nil
}

func (s *MemoryStore) SetPairingInfo(ctx context.Context, extensionID string, pi PairingInfo) error {
	_, ok := s.getItem(extensionID)
	if !ok {
		return errors.New("extension does not exists")
	}
	s.setItem(extensionID, Item{
		ExtensionID: extensionID,
		Expires:     time.Time{},
		PairingInfo: pi,
	})
	return nil
}

func (s *MemoryStore) setItem(key string, item Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.extensionsMap[key] = item
}

func (s *MemoryStore) getItem(key string) (Item, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.extensionsMap[key]
	return v, ok
}

type WSMemoryStore struct {
	mu          sync.Mutex
	devicesConn map[string]*websocket.Conn
}

func NewWSMemoryStore() *WSMemoryStore {
	return &WSMemoryStore{
		devicesConn: map[string]*websocket.Conn{},
	}
}

func (s *WSMemoryStore) SetMobileConn(ctx context.Context, deviceID string, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.devicesConn[deviceID] = conn
}

func (s *WSMemoryStore) GetMobileConn(ctx context.Context, deviceID string) (*websocket.Conn, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.devicesConn[deviceID]
	if !ok {
		return nil, false
	}
	return v, true
}