package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	pb "github.com/saladin2098/month3/lesson11/public_voting/genproto"
)

type PublicStorage struct {
	db *sql.DB
}

func NewPublicStorage(db *sql.DB) *PublicStorage {
	return &PublicStorage{db}
}

func (s *PublicStorage) CreatePublic(public *pb.PublicCreate) (*pb.Void, error) {
	query := `
		insert into public(
			id,
			first_name,
			last_name,
			birthday,
			gender,
            nation,
            party_id
		)
		values($1,$2,$3,$4,$5,$6,$7)
	`
	_, err := s.db.Exec(query,
		public.Id,
		public.FirstName,
		public.LastName,
		public.Birthday,
		public.Gender,
		public.Nation,
		public.Party,
	)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (s *PublicStorage) DeletePublic(id *pb.ById) (*pb.Void, error) {
	query := `
        delete from public
        where id = $1
    `
	_, err := s.db.Exec(query,
		id.Id,
	)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
func (s *PublicStorage) UpdatePublic(public *pb.PublicCreate) (*pb.Void, error) {
	query := `update public set `
	var conditions []string
	var args []interface{}
	if public.FirstName != "" {
		conditions = append(conditions, fmt.Sprintf("first_name = $%d", len(args)+1))
		args = append(args, public.FirstName)
	}
	if public.LastName != "" {
		conditions = append(conditions, fmt.Sprintf("last_name = $%d", len(args)+1))
		args = append(args, public.LastName)
	}
	if public.Birthday != "" {
		conditions = append(conditions, fmt.Sprintf("birthday = $%d", len(args)+1))
		args = append(args, public.Birthday)
	}
	if public.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", len(args)+1))
		args = append(args, public.Gender)
	}
	if public.Nation != "" {
		conditions = append(conditions, fmt.Sprintf("nation = $%d", len(args)+1))
		args = append(args, public.Nation)
	}
	if public.Party != "" {
		conditions = append(conditions, fmt.Sprintf("party_id = $%d", len(args)+1))
		args = append(args, public.Party)
	}
	if len(conditions) > 0 {
		query += strings.Join(conditions, ", ")
	}
	query += fmt.Sprintf(" where id = $%d",len(args)+1)
	args = append(args, public.Id)
	_, err := s.db.Exec(query, args...)
	if err!= nil {
        return nil, err
    }
	return &pb.Void{}, nil
}
func (s *PublicStorage) GetByIdPublic(id *pb.ById) (*pb.Public, error) {
	query := `
        select
            id,
            first_name,
            last_name,
            birthday,
            gender,
            nation,
            party_id
        from public
        where id = $1
    `
    row := s.db.QueryRow(query,
        id.Id,
    )
    public := &pb.Public{}
	var party_id string
    err := row.Scan(
        &public.Id,
        &public.FirstName,
        &public.LastName,
        &public.Birthday,
        &public.Gender,
        &public.Nation,
        &party_id,
    )
    if err!= nil {
        return nil, err
    }
	var party pb.Party

	query2 := `select * from party where id = $1`
	row2 := s.db.QueryRow(query2,party_id)
	err = row2.Scan(
		&party.Id,
		&party.Name,
		&party.Slogan,
		&party.OpenDate,
		&party.Description,
	)
	if err!= nil {
		return nil, err
	}
	public.Party = &party
    return public, nil
}
func (s *PublicStorage) GetAllPublics(filter *pb.Filter) (*pb.GetAllPublic, error) {
	query := `select * from public `
	var conditions []string
	var args []interface{}
	if filter.Party != "" {
        conditions = append(conditions, fmt.Sprintf("party = $%d", len(args)+1))
        args = append(args, filter.Party)
    }
	if filter.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", len(args)+1))
        args = append(args, filter.Gender)
	}
	if filter.Nation!= "" {
        conditions = append(conditions, fmt.Sprintf("nation = $%d", len(args)+1))
        args = append(args, filter.Nation)
    }
	if filter.Age!= 0 {
        conditions = append(conditions, fmt.Sprintf("extract(year(age(birthday))) >= $%d", len(args)+1))
        args = append(args, filter.Age)
    }
	if len(conditions) > 0 {
        query += "where " + strings.Join(conditions, " and ")
    }
	rows, err := s.db.Query(query, args...)
	if err!= nil {
        return nil, err
    }
	defer rows.Close()
	publics := pb.GetAllPublic{}
	for rows.Next() {
		public := pb.Public{}
		var party_id string
        err := rows.Scan(
            &public.Id,
            &public.FirstName,
            &public.LastName,
            &public.Birthday,
            &public.Gender,
            &public.Nation,
            &party_id,
        )
        if err!= nil {
            return nil, err
        }
		var party pb.Party

		query := `select * from party where id = $1`
		row := s.db.QueryRow(query,party_id)
		err = row.Scan(
			&party.Id,
			&party.Name,
			&party.Slogan,
			&party.OpenDate,
			&party.Description,
		)
		if err!= nil {
            return nil, err
        }
		public.Party = &party
        publics.Publics = append(publics.Publics, &public)
	}
	return &publics,nil
}
