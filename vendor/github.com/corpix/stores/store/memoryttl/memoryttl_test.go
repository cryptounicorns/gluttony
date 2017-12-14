package memoryttl

import (
	"sync"
	"testing"

	"github.com/corpix/time"
	"github.com/stretchr/testify/assert"

	"github.com/corpix/stores/errors"
)

type logger struct{}

func (logger) Write([]byte) (int, error)     { return 0, nil }
func (logger) Debugf(string, ...interface{}) {}
func (logger) Printf(string, ...interface{}) {}
func (logger) Errorf(string, ...interface{}) {}
func (logger) Fatalf(string, ...interface{}) {}
func (logger) Debug(...interface{})          {}
func (logger) Print(...interface{})          {}
func (logger) Error(...interface{})          {}
func (logger) Fatal(...interface{})          {}

func TestMemoryTTLSetGet(t *testing.T) {
	var (
		samples = []struct {
			name  string
			input map[string]interface{}
		}{
			{
				name:  "empty",
				input: map[string]interface{}{},
			},
			{
				name: "single",
				input: map[string]interface{}{
					"foo": "bar",
				},
			},
			{
				name: "multiple",
				input: map[string]interface{}{
					"foo": "bar",
					"bar": 1,
					"baz": nil,
				},
			},
		}

		log = logger{}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					store *MemoryTTL
					vv    interface{}
					err   error
				)

				store, err = New(
					Config{
						TTL:        1 * time.Second,
						Resolution: 1 * time.Second,
					},
					log,
				)
				if err != nil {
					t.Error(err)
					return
				}
				defer store.Close()

				for k, v := range sample.input {
					err = store.Set(k, v)
					assert.Equal(t, nil, err)
				}

				for k, v := range sample.input {
					vv, err = store.Get(k)
					assert.Equal(t, nil, err)
					assert.EqualValues(t, v, vv)
				}
			},
		)
	}
}

func TestMemoryTTLSetConcurrent(t *testing.T) {
	var (
		samples = []struct {
			name  string
			input map[string]interface{}
		}{
			{
				name:  "empty",
				input: map[string]interface{}{},
			},
			{
				name: "single",
				input: map[string]interface{}{
					"foo": "bar",
				},
			},
			{
				name: "multiple",
				input: map[string]interface{}{
					"foo": "bar",
					"bar": 1,
					"baz": nil,
				},
			},
		}

		log = logger{}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					wg    = &sync.WaitGroup{}
					store *MemoryTTL
					vv    interface{}
					err   error
				)

				store, err = New(
					Config{
						TTL:        1 * time.Second,
						Resolution: 1 * time.Second,
					},
					log,
				)
				if err != nil {
					t.Error(err)
					return
				}
				defer store.Close()

				for k, v := range sample.input {
					wg.Add(1)
					go func(k string, v interface{}) {
						err := store.Set(k, v)
						assert.Equal(t, nil, err)
						wg.Done()
					}(k, v)
				}

				wg.Wait()

				for k, v := range sample.input {
					vv, err = store.Get(k)
					assert.Equal(t, nil, err)
					assert.EqualValues(t, v, vv)
				}
			},
		)
	}
}

func TestMemoryTTLSetRemoveGet(t *testing.T) {
	var (
		samples = []struct {
			name  string
			input map[string]interface{}
		}{
			{
				name:  "empty",
				input: map[string]interface{}{},
			},
			{
				name: "single",
				input: map[string]interface{}{
					"foo": "bar",
				},
			},
			{
				name: "multiple",
				input: map[string]interface{}{
					"foo": "bar",
					"bar": 1,
					"baz": nil,
				},
			},
		}

		log = logger{}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					store *MemoryTTL
					vv    interface{}
					err   error
				)

				store, err = New(
					Config{
						TTL:        1 * time.Second,
						Resolution: 1 * time.Second,
					},
					log,
				)
				if err != nil {
					t.Error(err)
					return
				}
				defer store.Close()

				for k, v := range sample.input {
					err = store.Set(k, v)
					assert.Equal(t, nil, err)

					err = store.Remove(k)
					assert.Equal(t, nil, err)
				}

				for k, _ := range sample.input {
					vv, err = store.Get(k)

					switch err.(type) {
					case *errors.ErrKeyNotFound:
					default:
						t.Error(err)
					}

					assert.EqualValues(t, nil, vv)
				}
			},
		)
	}
}

func TestMemoryTTLKeys(t *testing.T) {
	var (
		samples = []struct {
			name   string
			input  map[string]interface{}
			output []string
		}{
			{
				name:   "empty",
				input:  map[string]interface{}{},
				output: []string{},
			},
			{
				name: "single",
				input: map[string]interface{}{
					"foo": "bar",
				},
				output: []string{"foo"},
			},
			{
				name: "multiple",
				input: map[string]interface{}{
					"foo": "bar",
					"bar": 1,
					"baz": nil,
				},
				output: []string{"foo", "bar", "baz"},
			},
		}

		log = logger{}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					store *MemoryTTL
					keys  []string
					err   error
				)

				store, err = New(
					Config{
						TTL:        1 * time.Second,
						Resolution: 1 * time.Second,
					},
					log,
				)
				if err != nil {
					t.Error(err)
					return
				}
				defer store.Close()

				for k, v := range sample.input {
					err = store.Set(k, v)
					assert.Equal(t, nil, err)
				}

				keys, err = store.Keys()
				assert.Equal(t, nil, err)
				assert.Equal(t, len(sample.output), len(keys))
				for _, v := range keys {
					assert.Contains(t, sample.output, v)
				}
			},
		)
	}
}

func TestMemoryTTLValues(t *testing.T) {
	var (
		samples = []struct {
			name   string
			input  map[string]interface{}
			output []interface{}
		}{
			{
				name:   "empty",
				input:  map[string]interface{}{},
				output: []interface{}{},
			},
			{
				name: "single",
				input: map[string]interface{}{
					"foo": "bar",
				},
				output: []interface{}{"bar"},
			},
			{
				name: "multiple",
				input: map[string]interface{}{
					"foo": "bar",
					"bar": 1,
					"baz": nil,
				},
				output: []interface{}{"bar", 1, nil},
			},
		}

		log = logger{}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					store  *MemoryTTL
					values []interface{}
					err    error
				)

				store, err = New(
					Config{
						TTL:        1 * time.Second,
						Resolution: 1 * time.Second,
					},
					log,
				)
				if err != nil {
					t.Error(err)
					return
				}
				defer store.Close()

				for k, v := range sample.input {
					err = store.Set(k, v)
					assert.Equal(t, nil, err)
				}

				values, err = store.Values()
				assert.Equal(t, nil, err)
				assert.Equal(t, len(sample.output), len(values))
				for _, v := range values {
					assert.Contains(t, sample.output, v)
				}
			},
		)
	}
}

func TestMemoryTTLMap(t *testing.T) {
	var (
		samples = []struct {
			name   string
			input  map[string]interface{}
			output []interface{}
		}{
			{
				name:  "empty",
				input: map[string]interface{}{},
			},
			{
				name: "single",
				input: map[string]interface{}{
					"foo": "bar",
				},
			},
			{
				name: "multiple",
				input: map[string]interface{}{
					"foo": "bar",
					"bar": 1,
					"baz": nil,
				},
			},
		}

		log = logger{}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					store *MemoryTTL
					m     map[string]interface{}
					err   error
				)

				store, err = New(
					Config{
						TTL:        1 * time.Second,
						Resolution: 1 * time.Second,
					},
					log,
				)
				if err != nil {
					t.Error(err)
					return
				}
				defer store.Close()

				for k, v := range sample.input {
					err = store.Set(k, v)
					assert.Equal(t, nil, err)
				}

				m, err = store.Map()
				assert.Equal(t, nil, err)
				assert.EqualValues(t, sample.input, m)
			},
		)
	}
}

func TestMemoryTTLExpire(t *testing.T) {
	var (
		samples = []struct {
			name   string
			input  map[string]interface{}
			output []interface{}
		}{
			{
				name:  "empty",
				input: map[string]interface{}{},
			},
			{
				name: "single",
				input: map[string]interface{}{
					"foo": "bar",
				},
			},
			{
				name: "multiple",
				input: map[string]interface{}{
					"foo": "bar",
					"bar": 1,
					"baz": nil,
				},
			},
		}

		log = logger{}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					store *MemoryTTL
					m     map[string]interface{}
					err   error
				)

				store, err = New(
					Config{
						TTL:        0 * time.Second,
						Resolution: 10 * time.Millisecond,
					},
					log,
				)
				if err != nil {
					t.Error(err)
					return
				}
				defer store.Close()

				for k, v := range sample.input {
					err = store.Set(k, v)
					assert.Equal(t, nil, err)
				}

				time.Sleep(100 * time.Millisecond)

				m, err = store.Map()
				assert.Equal(t, nil, err)
				assert.EqualValues(t, map[string]interface{}{}, m)
			},
		)
	}
}
