package repoiface

import (
	"context"
)

type SampleData interface {
	HealthCheck(ctx context.Context) error
}
