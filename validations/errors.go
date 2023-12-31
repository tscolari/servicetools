package validations

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Error(err error) error {
	if err == nil {
		return nil
	}

	return status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
}
