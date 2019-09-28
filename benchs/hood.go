package benchs

import (
	"database/sql"
	"fmt"

	"github.com/eaigner/hood"
)

var hd *hood.Hood

func init() {
	st := NewSuite("hood")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, HdInsert)
		st.AddBenchmark("BulkInsert 100 row", 500*ORM_MULTI, 0, HdInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, HdUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, 0, HdRead)
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, HdReadSlice)
		db, err := sql.Open("postgres", ORM_SOURCE)
		if err != nil {
			fmt.Printf("conn err: %v\n", err)
		}
		hd = hood.New(db, hood.NewPostgres())
	}
}

func HdInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 0; i < b.N; i++ {
		m.Id = 0
		if _, err := hd.Save(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func HdInsertMulti(b *B) {
	panic(fmt.Errorf("Problematic bulk insert, too slow"))
	var ms []Model
	wrapExecute(b, func() {
		initDB()
	})

	for i := 0; i < b.N; i++ {
		ms = make([]Model, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, *NewModel())
		}
		if _, err := hd.SaveAll(&ms); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func HdUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := hd.Save(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	for i := 0; i < b.N; i++ {
		if _, err := hd.Save(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func HdRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := hd.Save(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	//github.com/eaigner/hood/base.go:50 should be patched to `fieldValue.SetString(string(driverValue.Elem().String()))`
	for i := 0; i < b.N; i++ {
		var models []Model
		if err := hd.Where("id", "=", m.Id).Find(&models); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}

}

func HdReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			m.Id = 0
			if _, err := hd.Save(m); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})

	for i := 0; i < b.N; i++ {
		var models []Model
		if err := hd.Where("id", ">", 0).Limit(b.L).Find(&models); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}

}
