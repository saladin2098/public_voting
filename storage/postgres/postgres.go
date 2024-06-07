package postgres

import (
	"database/sql"
	"fmt"

	"github.com/saladin2098/month3/lesson11/public_voting/config"
	"github.com/saladin2098/month3/lesson11/public_voting/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	db      *sql.DB
	PublicS storage.PublicI
	PartyS  storage.PartyI
}

func ConnectDB() (*Storage, error) {
	cfg := config.Load()
	dbConn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase)
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	publicS := NewPublicStorage(db)
	partyS := NewPartyStorage(db)
	return &Storage{
        db:      db,
        PublicS: publicS,
        PartyS:  partyS,
    }, nil
}
func (s *Storage) Public() storage.PublicI {
	if s.PublicS == nil {
		s.PublicS = NewPublicStorage(s.db)
	}
	return s.PublicS
}
func (s *Storage) Party() storage.PartyI {
	if s.PartyS == nil {
		s.PartyS = NewPartyStorage(s.db)
	}
	return s.PartyS
}
