package lua

import (
	"github.com/corpix/loggers"
	"github.com/corpix/lua/mapper"
	"github.com/corpix/lua/pool"
	lua "github.com/yuin/gopher-lua"
)

const (
	Name = "lua"
)

func newVM(c Config, l loggers.Logger) func() *lua.LState {
	return func() *lua.LState {
		l.Debug("Creating new Lua VM")

		var (
			l   = lua.NewState()
			err error
		)

		err = l.DoString(c.Code)
		if err != nil {
			panic(err)
		}

		return l
	}
}

type Lua struct {
	config Config
	pool   *pool.Pool
	log    loggers.Logger
}

func (l *Lua) Preprocess(v interface{}) (interface{}, error) {
	var (
		vm  = l.pool.Get()
		lv  lua.LValue
		err error
	)
	defer l.pool.Put(vm)

	// FIXME: Dangerous indirection
	lv, err = mapper.ToValue(**(v.(**interface{})))
	if err != nil {
		return nil, err
	}

	err = vm.CallByParam(
		lua.P{
			Fn:      vm.GetGlobal(l.config.FunctionName),
			NRet:    1,
			Protect: true,
		},
		lv,
	)
	if err != nil {
		return nil, err
	}

	lv = vm.Get(-1)
	vm.Pop(1)

	return mapper.FromValue(lv)
}

func (l *Lua) Close() error {
	l.log.Debug("Closing")
	l.pool.Close()

	return nil
}

func New(c Config, l loggers.Logger) (*Lua, error) {
	var (
		p = pool.New(newVM(c, l))
	)

	return &Lua{
		config: c,
		pool:   p,
		log:    l,
	}, nil
}
