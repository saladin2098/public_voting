package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/saladin2098/month3/lesson11/public_voting/genproto"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	return db, mock
}

func TestCreateParty(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	storage := NewPartyStorage(db)
	party := &pb.Party{
		Id:          "1",
		Name:        "Party Name",
		Slogan:      "Party Slogan",
		OpenDate:    "2022-01-01",
		Description: "Party Description",
	}

	query := `INSERT INTO party\(id,name,slogan,open_date,description\) values\(\$1,\$2,\$3,\$4,\$5\)`
	mock.ExpectExec(query).
		WithArgs(party.Id, party.Name, party.Slogan, party.OpenDate, party.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := storage.CreateParty(party)
	if err != nil {
		t.Fatalf("failed to create party: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestDeleteParty(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	id := &pb.ById{Id: "1"}

	query := `DELETE FROM party WHERE id=\$1`
	mock.ExpectExec(query).WithArgs(id.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	storage := NewPartyStorage(db)
	_, err := storage.DeleteParty(id)
	if err != nil {
		t.Fatalf("failed to delete party: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUpdateParty(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	party := &pb.Party{
		Id:          "1",
		Name:        "Updated Party Name",
		Slogan:      "Updated Party Slogan",
		OpenDate:    "2022-01-01",
		Description: "Updated Party Description",
	}

	query := `UPDATE party SET name = \$1, slogan = \$2, open_date = \$3, description = \$4 WHERE id = \$5`
	mock.ExpectExec(query).
		WithArgs(party.Name, party.Slogan, party.OpenDate, party.Description, party.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	storage := NewPartyStorage(db)
	_, err := storage.UpdateParty(party)
	if err != nil {
		t.Fatalf("failed to update party: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetByIdParty(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	party := &pb.Party{
		Id:          "1",
		Name:        "Party Name",
		Slogan:      "Party Slogan",
		OpenDate:    "2022-01-01",
		Description: "Party Description",
	}

	query := `SELECT id, name, slogan, open_date, description FROM party WHERE id=\$1`
	rows := sqlmock.NewRows([]string{"id", "name", "slogan", "open_date", "description"}).
		AddRow(party.Id, party.Name, party.Slogan, party.OpenDate, party.Description)

	mock.ExpectQuery(query).WithArgs(party.Id).WillReturnRows(rows)

	storage := NewPartyStorage(db)
	gotParty, err := storage.GetByIdParty(&pb.ById{Id: "1"})
	if err != nil {
		t.Fatalf("failed to get party by id: %v", err)
	}

	if gotParty.Id != party.Id || gotParty.Name != party.Name || gotParty.Slogan != party.Slogan ||
		gotParty.OpenDate != party.OpenDate || gotParty.Description != party.Description {
		t.Errorf("got party %+v, want %+v", gotParty, party)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
func TestGetAllPartys(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	parties := []pb.Party{
		{Id: "1", Name: "Party One", Slogan: "Slogan One", OpenDate: "2022-01-01", Description: "Description One"},
		{Id: "2", Name: "Party Two", Slogan: "Slogan Two", OpenDate: "2022-01-02", Description: "Description Two"},
	}

	query := `SELECT id, name, slogan, open_date, description FROM party`
	rows := sqlmock.NewRows([]string{"id", "name", "slogan", "open_date", "description"})
	for _, party := range parties {
		rows.AddRow(party.Id, party.Name, party.Slogan, party.OpenDate, party.Description)
	}

	mock.ExpectQuery(query).WillReturnRows(rows)

	storage := NewPartyStorage(db)
	filter := &pb.Filter{}
	gotParties, err := storage.GetAllPartys(filter)
	if err != nil {
		t.Fatalf("failed to get all parties: %v", err)
	}

	if len(gotParties.Partys) != len(parties) {
		t.Errorf("got %d parties, want %d", len(gotParties.Partys), len(parties))
	}

	for i, party := range gotParties.Partys {
		if party.Id != parties[i].Id || party.Name != parties[i].Name || party.Slogan != parties[i].Slogan ||
			party.OpenDate != parties[i].OpenDate || party.Description != parties[i].Description {
			t.Errorf("got party %+v, want %+v", party, parties[i])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
