package middleware

import (
	"go-wal/pkg/helper/response"
	"go-wal/pkg/logger"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stackTrace := debug.Stack()
				logger.Error(ctx).Str("stack_trace", string(stackTrace)).Msg("[Oops Panic]")
				response.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, 500)
			}
		}()
		ctx.Next()
	}
}
