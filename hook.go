package mainager

import "context"

type Hook func(ctx context.Context) (context.Context, error)

type Hooks map[string]Hook
