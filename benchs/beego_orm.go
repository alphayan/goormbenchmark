package benchs

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

var bo orm.Ormer

func init() {
	st := NewSuite("beego_orm")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, BeegoOrmInsert)
		st.AddBenchmark("BulkInsert 100 row", 2000*ORM_MULTI, 0, BeegoOrmInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, BeegoOrmUpdate)
		st.AddBenchmark("Read", 2000*ORM_MULTI, 0, BeegoOrmRead)
		st.AddBenchmark("MultiRead limit 2000", 2000*ORM_MULTI, 2000, BeegoOrmReadSlice)

		orm.RegisterDataBase("default", "postgres", ORM_SOURCE, ORM_MAX_IDLE, ORM_MAX_CONN)
		orm.RegisterModel(new(Model))

		bo = orm.NewOrm()
	}
}

func BeegoOrmInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 0; i < b.N; i++ {
		m.Id = 0
		if _, err := bo.Insert(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func BeegoOrmInsertMulti(b *B) {
	var ms []*Model
	wrapExecute(b, func() {
		initDB()
		ms = make([]*Model, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, NewModel())
		}
	})

	for i := 0; i < b.N; i++ {
		if _, err := bo.InsertMulti(100, ms); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func BeegoOrmUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := bo.Insert(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		if _, err := bo.Update(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func BeegoOrmRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := bo.Insert(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		if err := bo.Read(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func BeegoOrmReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			m.Id = 0
			if _, err := bo.Insert(m); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})

	for i := 0; i < b.N; i++ {
		var models []*Model
		if _, err := bo.QueryTable("models").Filter("id__gt", 0).Limit(b.L).All(&models); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
