package payjs

import "github.com/yuyan2077/payjs/context"

// PayJS struct
type PayJS struct {
	Context *context.Context
}

// config for PayJS
type Config struct {
	Key   string
	Mchid string
}

func New(cfg *Config) *PayJS {
	context := new(context.Context)
	copyConfigToContext(cfg, context)
	return &Wechat{context}
}
