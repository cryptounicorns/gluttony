package memory

import (
	"sync"

	"github.com/corpix/loggers"

	"github.com/corpix/stores/errors"
)

const (
	Name = "memory"
)

type Memory struct {
	Config Config

	locker *sync.RWMutex
	store  map[string]interface{}
}

func (s *Memory) Name() string {
	return Name
}

func (s *Memory) Set(key string, value interface{}) error {
	s.locker.Lock()
	defer s.locker.Unlock()

	s.store[key] = value

	return nil
}

func (s *Memory) Get(key string) (interface{}, error) {
	s.locker.RLock()
	defer s.locker.RUnlock()

	v, ok := s.store[key]
	if !ok {
		return nil, errors.NewErrKeyNotFound(key)
	}

	return v, nil
}

func (s *Memory) Remove(key string) error {
	s.locker.Lock()
	defer s.locker.Unlock()

	_, ok := s.store[key]
	if !ok {
		return errors.NewErrKeyNotFound(key)
	}

	delete(s.store, key)

	return nil
}

func (s *Memory) Keys() ([]string, error) {
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

func (s *Memory) Values() ([]interface{}, error) {
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

func (s *Memory) Map() (map[string]interface{}, error) {
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

func (s *Memory) Iter(fn func(key string, value interface{}) bool) error {
	s.locker.RLock()
	defer s.locker.RUnlock()

	for k, v := range s.store {
		if fn(k, v) == false {
			break
		}
	}

	return nil
}

func (s *Memory) Close() error {
	return nil
}

func New(c Config, l loggers.Logger) (*Memory, error) {
	return &Memory{
		Config: c,
		locker: &sync.RWMutex{},
		store:  map[string]interface{}{},
	}, nil
}
