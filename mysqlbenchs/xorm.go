package benchs

import (
	"fmt"

	"xorm.io/xorm"
)

var xo *xorm.Session

func init() {
	st := NewSuite("xorm")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, XormInsert)
		st.AddBenchmark("BulkInsert 100 row", 2000*ORM_MULTI, 0, XormInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, XormUpdate)
		st.AddBenchmark("Read", 2000*ORM_MULTI, 0, XormRead)
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, XormReadSlice)

		engine, _ := xorm.NewEngine("mysql", ORM_SOURCE)

		engine.SetMaxIdleConns(ORM_MAX_IDLE)
		engine.SetMaxOpenConns(ORM_MAX_CONN)

		xo = engine.NewSession()
	}
}

func XormInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 0; i < b.N; i++ {
		m.Id = 0
		if _, err := xo.InsertOne(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func XormInsertMulti(b *B) {
	var ms []*Model
	wrapExecute(b, func() {
		initDB()
		ms = make([]*Model, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, NewModel())
		}
	})
	for i := 0; i < b.N; i++ {
		if _, err := xo.InsertMulti(&ms); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func XormUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := xo.Insert(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	for i := 0; i < b.N; i++ {
		if _, err := xo.Update(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func XormRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := xo.Insert(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		if _, err := xo.Get(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func XormReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			m.Id = 0
			if _, err := xo.Insert(m); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})

	for i := 0; i < b.N; i++ {
		var models []*Model
		if err := xo.Table("models").Where("id > ?", 0).Limit(b.L).Find(&models); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}

}
