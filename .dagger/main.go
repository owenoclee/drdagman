package main

import (
	"context"
	"dagger/drdagman/internal/dagger"
	"fmt"
)

type Drdagman struct{}

func (d *Drdagman) BuildSimpleNode(ctx context.Context, root *dagger.Directory, strToAppend string) *dagger.Container {
	return dag.Container().
		From("golang:1.23-alpine").
		WithDirectory("/drdagman", root).
		WithWorkdir("/drdagman/simplenode").
		WithExec([]string{"go", "build", "-ldflags", fmt.Sprintf("-X main.strToAppend=%s", strToAppend)}).
		WithEntrypoint([]string{"/drdagman/simplenode/simplenode"})
}
