package filter

import (
	"context"
	"fmt"

	"github.com/argoproj-labs/argo-dataflow/sdks/golang"

	"github.com/antonmedv/expr"

	"github.com/argoproj-labs/argo-dataflow/runner/util"
)

func Exec(ctx context.Context, x string) error {
	prog, err := expr.Compile(x)
	if err != nil {
		return fmt.Errorf("failed to compile %q: %w", x, err)
	}
	return golang.StartWithContext(ctx, func(ctx context.Context, msg []byte) ([]byte, error) {
		env, err := util.ExprEnv(ctx, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to create expr env: %w", err)
		}
		res, err := expr.Run(prog, env)
		if err != nil {
			return nil, fmt.Errorf("failed to run program %x: %w", x, err)
		}
		accept, ok := res.(bool)
		if !ok {
			return nil, fmt.Errorf("%q must return bool", x)
		}
		if accept {
			return msg, nil
		} else {
			return nil, nil
		}
	})
}
