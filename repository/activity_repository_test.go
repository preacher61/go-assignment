package repository

import (
	"context"
	"preacher61/go-assignment/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
)

func TestActivityRepositoryInsertActivitiesSucess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	data := []*model.Activity{
		{
			Activity: "test activity",
			Key:      "789789",
		},
	}

	//sqlmock.AnyArg().Match(data)
	mock.ExpectBegin()
	mock.ExpectPrepare("^COPY *").
		ExpectExec().
		WithArgs("789789", "test activity").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^COPY *").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ar := &ActivityRepository{
		db: db,
	}

	err = ar.InsertActivities(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestActivityRepositoryInsertActivitiesExecFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	data := []*model.Activity{
		{
			Activity: "test activity",
			Key:      "789789",
		},
	}

	//sqlmock.AnyArg().Match(data)
	mock.ExpectBegin()
	mock.ExpectPrepare("^COPY *").
		ExpectExec().
		WithArgs("789789", "test activity").
		WillReturnError(errors.New("i/o error"))

	ar := &ActivityRepository{
		db: db,
	}

	err = ar.InsertActivities(context.Background(), data)
	if err == nil {
		t.Fatal("error expected")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestActivityRepositoryInsertActivitiesPrepareFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	data := []*model.Activity{
		{
			Activity: "test activity",
			Key:      "789789",
		},
	}

	//sqlmock.AnyArg().Match(data)
	mock.ExpectBegin()
	mock.ExpectPrepare("^COPY *").
		WillReturnError(errors.New("i/o error"))

	ar := &ActivityRepository{
		db: db,
	}

	err = ar.InsertActivities(context.Background(), data)
	if err == nil {
		t.Fatal("error expected")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestActivityRepositoryInsertActivitiesCommitErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	data := []*model.Activity{
		{
			Activity: "test activity",
			Key:      "789789",
		},
	}

	//sqlmock.AnyArg().Match(data)
	mock.ExpectBegin()
	mock.ExpectPrepare("^COPY *").
		ExpectExec().
		WithArgs("789789", "test activity").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^COPY *").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(errors.New("i/o error"))

	ar := &ActivityRepository{
		db: db,
	}

	err = ar.InsertActivities(context.Background(), data)
	if err == nil {
		t.Fatal("error expected")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
