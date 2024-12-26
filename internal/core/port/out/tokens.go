package out

import "context"

type PortStoreInterface interface {
	Add(ctx context.Context, username, token string) error
	Get(ctc context.Context, username string) (string, error)
}
