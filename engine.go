package tinybreaker

import "sync"

type Engine struct {
	pool     sync.Pool
	breakers sync.Map
}

type Command struct {
	breaker *Breaker
}

func NewEngine() *Engine {
	engine := new(Engine)
	engine.pool.New = func() interface{} {
		return engine.allocateCommand()
	}

	return engine
}

func (engine *Engine) allocateCommand() *Command {
	return new(Command)
}

func (engine *Engine) GetCommand(name string) (cmd *Command) {
	cmd = engine.pool.Get().(*Command)

	cmd.reset()
	cmd.breaker = engine.getBreaker(name)

	return cmd
}

func (engine *Engine) getBreaker(name string) *Breaker {
	value, ok := engine.breakers.Load(name)
	if ok {
		return value.(*Breaker)
	}

	engine.RegisterBreaker(name)
	value, _ = engine.breakers.Load(name)
	return value.(*Breaker)
}

func (engine *Engine) RegisterBreaker(name string) {
	engine.breakers.Store(name, new(Breaker))
}

func (cmd *Command) reset() {
	cmd.breaker = nil
}
