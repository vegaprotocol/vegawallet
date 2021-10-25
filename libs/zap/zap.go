package zap

import (
	"fmt"

	"go.uber.org/zap"
)

func Sync(logger *zap.Logger) func() {
	return func() {
		err := logger.Sync()
		if err != nil {
			// This is the ultimate warning, as we can't do anything else.
			fmt.Printf("couldn't flush logger: %v", err) //nolint:forbidigo
		}
	}
}
