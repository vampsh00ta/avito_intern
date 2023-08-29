package repository

//
//func TestShouldUpdateStats(t *testing.T) {
//	mock, err := pgxmock.NewPool()
//	if err != nil {
//		t.Fatal(err)
//	}
//	ctx := context.Background()
//	begin := mock.ExpectBegin()
//	mock.ExpectExec("INSERT INTO users").
//		WithArgs("test1").
//		WillReturnResult(pgxmock.NewResult("INSERT", 1))
//
//	mock.ExpectCommit()
//	rep := New(mock, nil)
//	assert.Equal(t, rep.CreateUser(ctx, "test1"), nil)
//
//	// we make sure that all expectations were met
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//}
