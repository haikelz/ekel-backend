package app

import "ekel-backend/internal/env"

type Bootstrap struct {
	Port string
}

func InitBootstrap() Bootstrap {
	return Bootstrap{Port: env.GetEnv().PORT}
}
