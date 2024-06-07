package service

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/saladin2098/month3/lesson11/public_voting/genproto"
	"github.com/saladin2098/month3/lesson11/public_voting/storage/postgres"
)

type PublicService struct {
	stg *postgres.Storage
	pb.UnimplementedPublicServiceServer
}

func NewPublicService(stg *postgres.Storage) *PublicService {
	return &PublicService{stg: stg}
}
func (s *PublicService) CreatePublic(ctx context.Context, public *pb.PublicCreate) (*pb.Void, error) {
	id := uuid.NewString()
	public.Id = id
	_, err := s.stg.Public().CreatePublic(public)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
func (s *PublicService) UpdatePublic(ctx context.Context, public *pb.PublicCreate) (*pb.Void, error) {
	_, err := s.stg.Public().UpdatePublic(public)
    if err!= nil {
        return nil, err
    }
    return &pb.Void{}, nil
}
func (s *PublicService) DeletePublic(ctx context.Context, id *pb.ById) (*pb.Void, error) {
	res,err := s.stg.Public().DeletePublic(id)
    if err!= nil {
        return nil, err
    }
    return res, nil
}
func (s *PublicService) GetByIdPublic(ctx context.Context, id *pb.ById) (*pb.Public, error) {
	res,err := s.stg.Public().GetByIdPublic(id)
    if err!= nil {
        return nil, err
    }
    return res, nil
}
func (s *PublicService) GetAllPublics(ctx context.Context, filter *pb.Filter) (*pb.GetAllPublic, error) {
	res,err := s.stg.Public().GetAllPublics(filter)
    if err!= nil {
        return nil, err
    }
    return res, nil
}
