package env

import (
	"context"
	"os"

	"github.com/will-lol/personalWebsiteAwesome/lib/pointerify"
)

func GetEnv(ctx *context.Context) string {
	if ctx != nil {
		if req, ok := pointerify.DePointer(ctx).Value("env").(string); ok {
			return req
		}
	}
	return os.Getenv("ENVIRONMENT")
}
