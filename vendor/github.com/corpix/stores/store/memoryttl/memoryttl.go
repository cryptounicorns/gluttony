package memoryttl

import (
	"sync"
	"time"

	"github.com/corpix/loggers"

	"github.com/corpix/stores/errors"
)

const (
	Name = "memoryttl"
)

type MemoryTTL struct {
	Config Config

	log       loggers.Logger
	locker    *sync.RWMutex
	store     map[string]interface{}
	timeouted map[string]time.Time
	done      chan struct{}
}

func (s *MemoryTTL) Name() string {
	return Name
}

func (s *MemoryTTL) Set(key string, value interface{}) error {
	s.log.Debugf("Set key '%s' with value '%#v'", key, value)

	s.locker.Lock()
	defer s.locker.Unlock()

	s.store[key] = value
	s.timeouted[key] = time.Now().Add(s.Config.TTL.Duration())

	return nil
}

func (s *MemoryTTL) Get(key string) (interface{}, error) {
	s.locker.RLock()
	defer s.locker.RUnlock()

	v, ok := s.store[key]
	if !ok {
		return nil, errors.NewErrKeyNotFound(key)
	}

	return v, nil
}

func (s *MemoryTTL) Remove(key string) error {
	s.locker.Lock()
	defer s.locker.Unlock()

	return s.remove(key)
}

func (s *MemoryTTL) remove(key string) error {
	_, ok := s.store[key]
	if !ok {
		return errors.NewErrKeyNotFound(key)
	}

	delete(s.store, key)
	delete(s.timeouted, key)

	return nil
}

func (s *MemoryTTL) Keys() ([]string, error) {
	s.locker.RLock()
	defer s.locker.RUnlock()

	var (
		res = make([]string, len(s.store))
		n   = 0
	)

	for k, _ := range s.store {
		res[n] = k
		n++
	}

	return res, nil
}

func (s *MemoryTTL) Values() ([]interface{}, error) {
	s.locker.RLock()
	defer s.locker.RUnlock()

	var (
		res = make([]interface{}, len(s.store))
		n   = 0
	)

	for _, v := range s.store {
		res[n] = v
		n++
	}

	return res, nil
}

func (s *MemoryTTL) Map() (map[string]interface{}, error) {
	s.locker.RLock()
	defer s.locker.RUnlock()

	var (
		res = make(map[string]interface{}, len(s.store))
	)

	for k, v := range s.store {
		res[k] = v
	}

	return res, nil
}

func (s *MemoryTTL) Iter(fn func(key string, value interface{}) bool) error {
	s.locker.RLock()
	defer s.locker.RUnlock()

	for k, v := range s.store {
		if fn(k, v) == false {
			break
		}
	}

	return nil
}

func (s *MemoryTTL) Close() error {
	close(s.done)

	return nil
}

func (s *MemoryTTL) cancellationLoop() {
	var (
		resolution = s.Config.Resolution
	)

	for {
		select {
		case <-s.done:
			return
		case <-time.After(resolution.Duration()):
			s.locker.Lock()
			for k, v := range s.timeouted {
				if time.Now().After(v) {
					s.remove(k)
				}
			}
			s.locker.Unlock()
		}
	}
}

func New(c Config, l loggers.Logger) (*MemoryTTL, error) {
	var (
		s = &MemoryTTL{
			Config:    c,
			log:       l,
			locker:    &sync.RWMutex{},
			store:     map[string]interface{}{},
			timeouted: map[string]time.Time{},
			done:      make(chan struct{}),
		}
	)

	go s.cancellationLoop()

	return s, nil
}
