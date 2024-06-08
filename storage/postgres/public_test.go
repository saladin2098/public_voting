package postgres


import (
	// "context"
	// "database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/saladin2098/month3/lesson11/public_voting/genproto"
)
func TestCreatePublic(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	storage := NewPublicStorage(db)
	public := &pb.PublicCreate{
		Id:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Birthday:  "2000-01-01",
		Gender:    "Male",
		Nation:    "USA",
		Party:     "1",
	}

	query := `insert into public\(id,first_name,last_name,birthday,gender,nation,party_id\) values\(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)`
	mock.ExpectExec(query).
		WithArgs(public.Id, public.FirstName, public.LastName, public.Birthday, public.Gender, public.Nation, public.Party).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := storage.CreatePublic(public)
	if err != nil {
		t.Fatalf("failed to create public: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestDeletePublic(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	id := &pb.ById{Id: "1"}

	query := `delete from public where id=\$1`
	mock.ExpectExec(query).WithArgs(id.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	storage := NewPublicStorage(db)
	_, err := storage.DeletePublic(id)
	if err != nil {
		t.Fatalf("failed to delete public: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}


func TestUpdatePublic(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	public := &pb.PublicCreate{
		Id:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Birthday:  "2000-01-01",
		Gender:    "Male",
		Nation:    "USA",
		Party:     "1",
	}

	query := `update public set first_name = \$1, last_name = \$2, birthday = \$3, gender = \$4, nation = \$5, party_id = \$6 where id = \$7`
	mock.ExpectExec(query).
		WithArgs(public.FirstName, public.LastName, public.Birthday, public.Gender, public.Nation, public.Party, public.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	storage := NewPublicStorage(db)
	_, err := storage.UpdatePublic(public)
	if err != nil {
		t.Fatalf("failed to update public: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetByIdPublic(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	public := &pb.Public{
		Id:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Birthday:  "2000-01-01",
		Gender:    "Male",
		Nation:    "USA",
		Party:     &pb.Party{Id: "1", Name: "Party Name", Slogan: "Party Slogan", OpenDate: "2022-01-01", Description: "Party Description"},
	}

	query := `select id, first_name, last_name, birthday, gender, nation, party_id from public where id=\$1`
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "birthday", "gender", "nation", "party_id"}).
		AddRow(public.Id, public.FirstName, public.LastName, public.Birthday, public.Gender, public.Nation, public.Party.Id)

	mock.ExpectQuery(query).WithArgs(public.Id).WillReturnRows(rows)

	query2 := `select id, name, slogan, open_date, description from party where id=\$1`
	rows2 := sqlmock.NewRows([]string{"id", "name", "slogan", "open_date", "description"}).
		AddRow(public.Party.Id, public.Party.Name, public.Party.Slogan, public.Party.OpenDate, public.Party.Description)

	mock.ExpectQuery(query2).WithArgs(public.Party.Id).WillReturnRows(rows2)

	storage := NewPublicStorage(db)
	gotPublic, err := storage.GetByIdPublic(&pb.ById{Id: "1"})
	if err != nil {
		t.Fatalf("failed to get public by id: %v", err)
	}

	if gotPublic.Id != public.Id || gotPublic.FirstName != public.FirstName || gotPublic.LastName != public.LastName ||
		gotPublic.Birthday != public.Birthday || gotPublic.Gender != public.Gender || gotPublic.Nation != public.Nation ||
		gotPublic.Party.Id != public.Party.Id || gotPublic.Party.Name != public.Party.Name || gotPublic.Party.Slogan != public.Party.Slogan ||
		gotPublic.Party.OpenDate != public.Party.OpenDate || gotPublic.Party.Description != public.Party.Description {
		t.Errorf("got public %+v, want %+v", gotPublic, public)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetAllPublics(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	publics := []pb.Public{
		{Id: "1", FirstName: "John", LastName: "Doe", Birthday: "2000-01-01", Gender: "Male", Nation: "USA", Party: &pb.Party{Id: "1", Name: "Party Name", Slogan: "Party Slogan", OpenDate: "2022-01-01", Description: "Party Description"}},
		{Id: "2", FirstName: "Jane", LastName: "Doe", Birthday: "2000-02-01", Gender: "Female", Nation: "USA", Party: &pb.Party{Id: "1", Name: "Party Name", Slogan: "Party Slogan", OpenDate: "2022-01-01", Description: "Party Description"}},
	}

	query := `select p.id, p.first_name, p.last_name, p.birthday, p.gender, p.nation, pt.id, pt.name, pt.slogan, pt.open_date, pt.description from public p left join party pt on p.party_id = pt.id`
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "birthday", "gender", "nation", "party_id", "name", "slogan", "open_date", "description"})
	for _, public := range publics {
		rows.AddRow(public.Id, public.FirstName, public.LastName, public.Birthday, public.Gender, public.Nation, public.Party.Id, public.Party.Name, public.Party.Slogan, public.Party.OpenDate, public.Party.Description)
	}

	mock.ExpectQuery(query).WillReturnRows(rows)

	storage := NewPublicStorage(db)
	filter := &pb.Filter{}
	gotPublics, err := storage.GetAllPublics(filter)
	if err != nil {
		t.Fatalf("failed to get all publics: %v", err)
	}

	if len(gotPublics.Publics) != len(publics) {
		t.Errorf("got %d publics, want %d", len(gotPublics.Publics), len(publics))
	}

	for i, want := range publics {
		got := gotPublics.Publics[i]
		if got.Id != want.Id || got.FirstName != want.FirstName || got.LastName != want.LastName || got.Birthday != want.Birthday || got.Gender != want.Gender || got.Nation != want.Nation {
			t.Errorf("got public %+v, want %+v", got, want)
		}
		if got.Party == nil || got.Party.Id != want.Party.Id || got.Party.Name != want.Party.Name || got.Party.Slogan != want.Party.Slogan || got.Party.OpenDate != want.Party.OpenDate || got.Party.Description != want.Party.Description {
			t.Errorf("got public party %+v, want %+v", got.Party, want.Party)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
