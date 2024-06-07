package storage
import pb "github.com/saladin2098/month3/lesson11/public_voting/genproto"

type StorageI interface {
	Public() PublicI
	Party() PartyI
}
type PublicI interface {
	CreatePublic(public *pb.PublicCreate) (*pb.Void, error)
	DeletePublic(id *pb.ById) (*pb.Void, error) 
	UpdatePublic(public *pb.PublicCreate) (*pb.Void, error)
    GetByIdPublic(id *pb.ById) (*pb.Public, error) 
	GetAllPublics(filter *pb.Filter) (*pb.GetAllPublic, error)
}
type PartyI interface {
	CreateParty(party *pb.Party) (*pb.Void, error)
	DeleteParty(id *pb.ById) (*pb.Void, error) 
	UpdateParty(party *pb.Party) (*pb.Void, error)
	GetByIdParty(id *pb.ById) (*pb.Party, error) 
	GetAllPartys(filter *pb.Filter) (*pb.GetAllParty, error) 
}