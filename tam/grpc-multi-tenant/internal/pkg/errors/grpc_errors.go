package errors

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Common error types that can be used across services
const (
	ErrInvalidArgument   = "invalid_argument"
	ErrNotFound          = "not_found"
	ErrAlreadyExists     = "already_exists"
	ErrPermissionDenied  = "permission_denied"
	ErrUnauthenticated   = "unauthenticated"
	ErrInternal          = "internal"
	ErrResourceExhausted = "resource_exhausted"
	ErrFailedPrecond     = "failed_precondition"
)

// GRPCError creates a gRPC error with the appropriate status code based on the error type
func GRPCError(errType string, msg string) error {
	var code codes.Code

	switch errType {
	case ErrInvalidArgument:
		code = codes.InvalidArgument
	case ErrNotFound:
		code = codes.NotFound
	case ErrAlreadyExists:
		code = codes.AlreadyExists
	case ErrPermissionDenied:
		code = codes.PermissionDenied
	case ErrUnauthenticated:
		code = codes.Unauthenticated
	case ErrResourceExhausted:
		code = codes.ResourceExhausted
	case ErrFailedPrecond:
		code = codes.FailedPrecondition
	default:
		code = codes.Internal
	}

	return status.Error(code, msg)
}

// InvalidArgument creates a gRPC error with InvalidArgument status code
func InvalidArgument(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

// NotFound creates a gRPC error with NotFound status code
func NotFound(entity string, id string) error {
	return status.Error(codes.NotFound, fmt.Sprintf("%s with ID %s not found", entity, id))
}

// AlreadyExists creates a gRPC error with AlreadyExists status code
func AlreadyExists(msg string) error {
	return status.Error(codes.AlreadyExists, msg)
}

// InternalError creates a gRPC error with Internal status code
func InternalError(msg string) error {
	return status.Error(codes.Internal, fmt.Sprintf("internal error: %s", msg))
}

// FailedPrecondition creates a gRPC error with FailedPrecondition status code
func FailedPrecondition(msg string) error {
	return status.Error(codes.FailedPrecondition, msg)
}

// FromError converts a standard error into a gRPC error with an appropriate status code
func FromError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Map common error messages to appropriate gRPC error codes
	switch {
	case strings.Contains(errMsg, "not found"):
		return status.Error(codes.NotFound, errMsg)
	case strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "already joined") ||
		strings.Contains(errMsg, "already accepted"):
		return status.Error(codes.AlreadyExists, errMsg)
	case strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "required"):
		return status.Error(codes.InvalidArgument, errMsg)
	default:
		return status.Error(codes.Internal, errMsg)
	}
}
