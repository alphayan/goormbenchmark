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
		st.AddBenchmark("BulkInsert 100 row", 500*ORM_MULTI, 0, DbrInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, DbrUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, 0, DbrRead)
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, DbrReadSlice)

		conn, _ := dbr.Open("msyql", ORM_SOURCE, nil)
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

	for i := 0; i < b.N; i++ {
		m.Id = 0
		if _, err := dbrsession.InsertInto("model").Columns("name", "title", "fax", "web", "age", "counter").Record(m).Exec(); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func DbrInsertMulti(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})
	q := dbrsession.InsertInto("model").Columns("name", "title", "fax", "web", "age", "counter")
	for i := 0; i < b.N; i++ {
		m.Id = 0
		q = q.Record(m)
	}

	if _, err := q.Exec(); err != nil {
		fmt.Println(err)
		b.FailNow()
	}

}

func DbrUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := dbrsession.InsertInto("model").Columns("name", "title", "fax", "web", "age", "counter").Record(m).Exec(); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

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
		if _, err := dbrsession.InsertInto("model").Columns("name", "title", "fax", "web", "age", "counter").Record(m).Exec(); err != nil {
			fmt.Printf("insert before read err: %v\n", err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		if _, err := dbrsession.Select("*").From("model").Where("id = ?", m.Id).Load(&m); err != nil {
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
			if _, err := dbrsession.InsertInto("model").Columns("name", "title", "fax", "web", "age", "counter").Record(m).Exec(); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})
	for i := 0; i < b.N; i++ {
		var m []Model
		if _, err := dbrsession.Select("*").From("model").Where("id > ?", 0).Limit(uint64(b.L)).Load(&m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
