package mainager

import "context"

type Hook func(ctx context.Context, params ...interface{}) (context.Context, error)

type Hooks map[string]Hook
