//go:generate mockgen -source=kv.go -destination=mock/mock_kv.go -package=mock
package storage

import (
	"context"

	"github.com/Terence1105/Tinder/pkg/storage/redis/tinder"
	"github.com/Terence1105/Tinder/pkg/storage/redis/tinder/dto"
)

var _ TinderStorage = (*tinder.TinderKV)(nil)

type TinderStorage interface {
	AddPerson(ctx context.Context, person *dto.Person) error
	RemovePerson(ctx context.Context, name string, gender int) error
	GetPeople(ctx context.Context, min, max float64, count, gender int) ([]dto.Person, error)
	DecrementDateCount(ctx context.Context, name string) (int, error)
	GetDateCount(ctx context.Context, name string) (string, error)
}
