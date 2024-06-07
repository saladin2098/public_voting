package service

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/saladin2098/month3/lesson11/public_voting/genproto"
	"github.com/saladin2098/month3/lesson11/public_voting/storage/postgres"
)

type PartyService struct {
	stg *postgres.Storage
	pb.UnimplementedPartyServiceServer
}

func NewPartyService(stg *postgres.Storage) *PartyService {
	return &PartyService{stg: stg}
}

func (s *PartyService) CreateParty(ctx context.Context, party *pb.Party) (*pb.Void, error) {
	id := uuid.NewString()
	party.Id = id
	_, err := s.stg.Party().CreateParty(party)
	if err!= nil {
        return nil, err
    }
	return &pb.Void{}, nil
}
func (s *PartyService) UpdateParty(ctx context.Context, party *pb.Party) (*pb.Void, error) {
	_, err := s.stg.PartyS.UpdateParty(party)
    if err!= nil {
        return nil, err
    }
    return &pb.Void{}, nil
}
func (s *PartyService) DeleteParty(ctx context.Context, id *pb.ById) (*pb.Void, error) {
	res,err := s.stg.PartyS.DeleteParty(id)
	if err!= nil {
        return nil, err
    }
	return res, nil
}
func (s *PartyService) GetByIdParty(ctx context.Context, id *pb.ById) (*pb.Party, error) {
	res,err := s.stg.PartyS.GetByIdParty(id)
    if err!= nil {
        return nil, err
    }
    return res, nil
}
func (s *PartyService) GetAllPartys(ctx context.Context, filter *pb.Filter) (*pb.GetAllParty, error) {
	res,err := s.stg.PartyS.GetAllPartys(filter)
    if err!= nil {
        return nil, err
    }
    return res, nil
}