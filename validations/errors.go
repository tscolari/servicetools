package validations

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error is a simple wrapper on an error that will
// transform it into an invalid argument status error.
func Error(err error) error {
	if err == nil {
		return nil
	}

	return status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
}
