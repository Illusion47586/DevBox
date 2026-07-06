package service

import "github.com/dhruv/devbox/internal/devbox/state"

type Operator struct {
	Config RuntimeConfig
	Store  *state.Store
}

func NewOperator(config RuntimeConfig) *Operator {
	return &Operator{Config: config, Store: state.NewStore(config.StatePath)}
}
