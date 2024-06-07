package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	pb "github.com/saladin2098/month3/lesson11/public_voting/genproto"
)

type PartyStorage struct {
	db *sql.DB
}
func NewPartyStorage(db *sql.DB) *PartyStorage {
	return &PartyStorage{db}
}

func (s *PartyStorage) CreateParty(party *pb.Party) (*pb.Void, error) {
	query := `
		INSERT INTO party(id,name,slogan,open_date,description) 
		values($1,$2,$3,$4,$5)
		`
	_, err := s.db.Exec(query,
		party.Id,
		party.Name,
		party.Slogan,
		party.OpenDate,
		party.Description)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (s *PartyStorage) DeleteParty(id *pb.ById) (*pb.Void, error) {
	query := `
        DELETE FROM party WHERE id=$1
        `
	_, err := s.db.Exec(query,
		id.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (s *PartyStorage) UpdateParty(party *pb.Party) (*pb.Void, error) {
	var conditions []string
	var args []interface{}
	query := `update party set`

	if party.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, party.Name)
	}
	if party.Slogan != "" {
		conditions = append(conditions, fmt.Sprintf("slogan = $%d", len(args)+1))
		args = append(args, party.Slogan)
	}
	if party.OpenDate != "" {
		conditions = append(conditions, fmt.Sprintf("open_date = $%d", len(args)+1))
		args = append(args, party.OpenDate)
	}
	if party.Description != "" {
		conditions = append(conditions, fmt.Sprintf("description = $%d", len(args)+1))
		args = append(args, party.Description)
	}

	if len(conditions) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(conditions, ", ")
	query += fmt.Sprintf(" where id = $%d", len(args)+1)
	args = append(args, party.Id)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}
func (s *PartyStorage) GetByIdParty(id *pb.ById) (*pb.Party, error) {
	query := `
        SELECT id,
			name,
			slogan,
			open_date,
			description 
		    FROM party 
			WHERE id=$1
        `
    row := s.db.QueryRow(query,
        id.Id)
    party := &pb.Party{}
    err := row.Scan(
        &party.Id,
        &party.Name,
        &party.Slogan,
        &party.OpenDate,
        &party.Description,
    )
    if err!= nil {
        return nil, err
    }
    return party, nil
}
func (s *PartyStorage) GetAllPartys(filter *pb.Filter) (*pb.GetAllParty, error) {
	var conditions []string
	var args []interface{}
	query := `
        SELECT id,
            name,
            slogan,
            open_date,
            description 
        FROM party
    `

    if filter.Date != "" {
		conditions = append(conditions, fmt.Sprintf("open_date = $%d", len(args)+1))
		args = append(args, filter.Date)
    }
	if filter.Slogan != "" {
		conditions = append(conditions, fmt.Sprintf("slogan = $%d", len(args)+1))
		args = append(args, filter.Slogan)
    }
	
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parties := &pb.GetAllParty{}
	for rows.Next() {
		party := &pb.Party{}
        err := rows.Scan(
            &party.Id,
            &party.Name,
            &party.Slogan,
            &party.OpenDate,
            &party.Description,
        )
        if err != nil {
            return nil, err
        }
        parties.Partys = append(parties.Partys, party)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return parties, nil
}