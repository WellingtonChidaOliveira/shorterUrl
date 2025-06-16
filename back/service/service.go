package service

import (
	"context"
	"fmt"

	"github.com/wellingtonchida/shortner/models"
	"github.com/wellingtonchida/shortner/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShorterService interface {
	GetAll(ctx context.Context) ([]models.Shorter, error)
	GetById(ctx context.Context, id string) (models.Shorter, error)
	Create(ctx context.Context, shorter models.Shorter) (string, error)
	Update(ctx context.Context, id string, shorter models.Shorter) error
	Delete(ctx context.Context, id string) error
	Inactivate(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) error 
}

type shorterService struct {
	repo repository.DataBase
}

func Service(repo repository.DataBase) ShorterService {
	return &shorterService{
		repo: repo,
	}
}

func (s *shorterService) GetAll(ctx context.Context) ([]models.Shorter, error) {
	shorters, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return shorters, nil
}

func (s *shorterService) GetById(ctx context.Context, id string) (models.Shorter, error) {
	shorter, err := s.repo.GetById(ctx, id)
	if err != nil {
		return models.Shorter{}, err
	}
	return shorter, nil
}

func (s *shorterService) Create(ctx context.Context, shorter models.Shorter) (string, error) {
	hasShorter, err := s.GetById(ctx, shorter.Shorter)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Println(err)
		return "", err
	}

	fmt.Printf("has data %v\n", hasShorter)

	if !hasShorter.ID.IsZero(){
		return "", fmt.Errorf("already exist this shorter: %s", shorter.Shorter)
	}

	fmt.Println("\nno has other register")

	oid, err := s.repo.Create(ctx, shorter)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Printf("add oid %v\n", oid)
	return oid, err
}

func (s *shorterService) Update(ctx context.Context, id string, shorter models.Shorter) error {
	fmt.Printf("id to update %v\n",id)
	err := s.repo.Update(ctx, id, shorter)
	if err != nil {
		return err
	}
	return nil
}
func (s *shorterService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (s *shorterService) Inactivate(ctx context.Context, id string) error {
	err := s.repo.Inactivate(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (s *shorterService) Activate(ctx context.Context, id string) error {
	err := s.repo.Activate(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

