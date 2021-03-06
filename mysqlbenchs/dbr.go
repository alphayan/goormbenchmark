package benchs

import (
	"fmt"

	"github.com/gocraft/dbr"
)

var dbrsession *dbr.Session

func init() {
	st := NewSuite("dbr")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, DbrInsert)
		st.AddBenchmark("BulkInsert 100 row", 2000*ORM_MULTI, 0, DbrInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, DbrUpdate)
		st.AddBenchmark("Read", 2000*ORM_MULTI, 0, DbrRead)
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, DbrReadSlice)
		conn, _ := dbr.Open("mysql", ORM_SOURCE, nil)
		conn.SetMaxIdleConns(ORM_MAX_IDLE)
		conn.SetMaxOpenConns(ORM_MAX_CONN)
		sess := conn.NewSession(nil)
		dbrsession = sess
	}
}

func DbrInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Id = 0
		if _, err := dbrsession.InsertInto("models").Columns("name", "title", "fax", "web", "age", "counter").Record(m).Exec(); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func DbrInsertMulti(b *B) {
	panic(fmt.Errorf("Don't support bulk insert"))
}

func DbrUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := dbrsession.InsertInto("models").Columns("name", "title", "fax", "web", "age", "counter").Record(m).Exec(); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := dbrsession.Update("models").
			Set("name", m.Name).
			Set("title", m.Title).
			Set("fax", m.Fax).
			Set("web", m.Web).
			Set("age", m.Age).
			Set("counter", m.Counter).Exec(); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func DbrRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := dbrsession.InsertInto("models").Columns("name", "title", "fax", "web", "age", "counter").Record(m).Exec(); err != nil {
			fmt.Printf("insert before read err: %v\n", err)
			b.FailNow()
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := dbrsession.Select("*").From("models").Where("id = ?", m.Id).Load(&m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func DbrReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			if _, err := dbrsession.InsertInto("models").Columns("name", "title", "fax", "web", "age", "counter").Record(m).Exec(); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var m []Model
		if _, err := dbrsession.Select("*").From("models").Where("id > ?", 0).Limit(uint64(b.L)).Load(&m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
