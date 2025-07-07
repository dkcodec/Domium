package repository

import "context"

type PropertyRepository interface {
	Create(ctx context.Context, p *Property) (string, error)
	Get(ctx context.Context, id string) (*Property, error)
	Update(ctx context.Context, id string, p *Property) error
	Delete(ctx context.Context, id string) error
}
