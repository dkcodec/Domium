package service

import (
	"context"
	"property-service/internal/repository"
	"time"

	"github.com/nats-io/nats.go"
)

type PropertyService struct {
	repo repository.Repository
	nats *nats.Conn
}

func New(repo repository.Repository, nc *nats.Conn) *PropertyService {
	return &PropertyService{
		repo: repo,
		nats: nc,
	}
}

func (s *PropertyService) Create(ctx context.Context, p *repository.Property) (string, error) {
	id, err := s.repo.Create(ctx, p)

	if err == nil {
		_ = s.nats.Publish("property.created", []byte(id))
	}

	return id, err
}

func (s *PropertyService) Get(ctx context.Context, id string) (*repository.Property, error) {
	return s.repo.Get(ctx, id)
}

func (s *PropertyService) Update(ctx context.Context, id string, p *repository.Property) error {
	p.UpdatedAt = time.Now()
	err := s.repo.Update(ctx, id, p)
	if err == nil {
		_ = s.nats.Publish("property.updated", []byte(id))
	}

	return err
}

func (s *PropertyService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		_ = s.nats.Publish("property.deleted", []byte(id))
	}

	return err
}