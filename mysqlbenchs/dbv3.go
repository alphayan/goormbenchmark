package benchs

import (
	"fmt"

	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

var sess sqlbuilder.Database

func init() {
	st := NewSuite("db.v3")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, DBv3Insert)
		st.AddBenchmark("BulkInsert 100 row", 500*ORM_MULTI, 0, DBv3InsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, DBv3Update)
		st.AddBenchmark("Read", 4000*ORM_MULTI, 0, DBv3Read)
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, DBv3ReadSlice)
		settings, err := mysql.ParseURL(ORM_SOURCE)
		checkErr(err)
		db, err := mysql.Open(settings)
		checkErr(err)
		sess = db
	}
}

func DBv3Insert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})
	col := sess.Collection("models")
	var err error
	for i := 0; i < b.N; i++ {
		m.Id = 0
		err = col.InsertReturning(m)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func DBv3InsertMulti(b *B) {
	wrapExecute(b, func() {
		initDB()
	})
	var err error
	for i := 0; i < b.N; i++ {
		v := NewModel()
		batch := sess.InsertInto("models").Columns("name", "title", "fax", "web", "age", "counter").Batch(100)
		go func() {
			for i := 0; i < 100; i++ {
				batch.Values(v.Name, v.Title, v.Fax, v.Web, v.Age, v.Counter)
			}
			batch.Done()
		}()
		if batch.Wait() != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func DBv3Update(b *B) {
	var m *Model
	col := sess.Collection("models")
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if err := col.InsertReturning(m); err != nil {
			b.FailNow()
		}
	})
	for i := 0; i < b.N; i++ {
		if err := col.UpdateReturning(m); err != nil {
			fmt.Printf("update err: %v\n", err)
			b.FailNow()
		}
	}
}

func DBv3Read(b *B) {
	var m *Model
	col := sess.Collection("models")
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if err := col.InsertReturning(m); err != nil {
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		if err := sess.SelectFrom("models").Where("id = ?", m.Id).One(m); err != nil {
			b.FailNow()
		}
	}
}

func DBv3ReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		col := sess.Collection("models")
		for i := 0; i < b.L; i++ {
			m.Id = 0
			if _, err := col.Insert(m); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})
	for i := 0; i < b.N; i++ {
		var models []*Model
		if err := sess.SelectFrom("models").Where("id > ?", 0).Limit(b.L).All(&models); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
