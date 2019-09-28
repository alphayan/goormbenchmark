package benchs

import (
	"fmt"

	"github.com/coocood/qbs"
)

var qo *qbs.Qbs

func init() {
	st := NewSuite("qbs")
	st.InitF = func() {
		st.AddBenchmark("Insert", 500*ORM_MULTI, 0, QbsInsert)
		st.AddBenchmark("BulkInsert 100 row", 500*ORM_MULTI, 0, QbsInsertMulti)
		st.AddBenchmark("Update", 500*ORM_MULTI, 0, QbsUpdate)
		st.AddBenchmark("Read", 500*ORM_MULTI, 0, QbsRead)
		st.AddBenchmark("MultiRead limit 500", 500*ORM_MULTI, 500, QbsReadSlice)

		qbs.Register("postgres", ORM_SOURCE, "q_model", qbs.NewPostgres())
		qbs.ChangePoolSize(ORM_MAX_IDLE)
		qbs.SetConnectionLimit(ORM_MAX_CONN, true)
		var err error
		qo, err = qbs.GetQbs()
		if err != nil {
			fmt.Printf("conn err: %v\n", err)
		}
	}
}

func QbsInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})
	defer qo.Close()
	for i := 0; i < b.N; i++ {
		m.Id = 0
		if _, err := qo.Save(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func QbsInsertMulti(b *B) {
	panic(fmt.Errorf("Don't support bulk insert, err driver: bad connection"))
	var ms []*Model
	wrapExecute(b, func() {
		initDB()

		ms = make([]*Model, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, NewModel())
		}
	})
	for i := 0; i < b.N; i++ {
		if err := qo.BulkInsert(ms); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func QbsUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		qo.Save(m)
	})

	for i := 0; i < b.N; i++ {
		if _, err := qo.Save(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func QbsRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		qo.Save(m)
	})

	for i := 0; i < b.N; i++ {
		if err := qo.Find(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func QbsReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			m.Id = 0
			if _, err := qo.Save(m); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})
	for i := 0; i < b.N; i++ {
		var models []*Model
		if err := qo.Where("id > ?", 0).Limit(b.L).FindAll(&models); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
