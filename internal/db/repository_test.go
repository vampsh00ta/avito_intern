package repository

//func TestShouldUpdateStats(t *testing.T) {
//	mock, err := pgxmock.NewConn()
//	if err != nil {
//		t.Fatal(err)
//	}
//	db := New(mock)
//	//rows := mock.NewRows([]string{"id", "title", "body"}).
//	//	AddRow(1, "post 1", "hello").
//	//	AddRow(2, "post 2", "world")
//
//	//mock.ExpectQuery("^SELECT (.+) FROM posts$").WillReturnRows(rows)
//
//	ctx := context.TODO()
//	defer mock.Close(ctx)
//	tx, err := mock.ExpectBeginTx
//
//	if err != nil {
//		t.Fatal(err)
//	}
//	mock.ExpectBegin()
//
//	testData := Segment{Slug: "test"}
//	mock.ExpectExec("INSERT INTO users").
//		WithArgs(testData).WillReturnResult(pgxmock.NewResult("INSERT", 1))
//
//	_, err = db.CreateSegment(ctx, testData)
//	if err != nil {
//		tx.Rollback(ctx)
//		t.Errorf("error was not expected while updating: %s", err)
//	}
//
//	// we make sure that all expectations were met
//	if err := mock.ExpectationsWereMet(); err != nil {
//		tx.Rollback(ctx)
//
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//	tx.Commit(ctx)
//}
