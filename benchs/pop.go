package benchs

import (
	"fmt"

	"github.com/gobuffalo/pop"
)

var popdb *pop.Connection

func PopConnect(name string) (*pop.Connection, error) {
	deet := &pop.ConnectionDetails{
		URL:  "postgres://postgres:root123456@192.168.199.248:5432/test?sslmode=disable",
		Pool: 4,
	}
	if c, err := pop.NewConnection(deet); err != nil {
		return nil, err
	} else {
		pop.Connections[name] = c
		return pop.Connections[name], nil
	}
}

func init() {
	st := NewSuite("pop")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, PopInsert)
		st.AddBenchmark("BulkInsert 100 row", 500*ORM_MULTI, 0, PopInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, PopUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, 0, PopRead)
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, PopReadSlice)
		var err error
		popdb, err = PopConnect("bechdb")
		if err != nil {
			fmt.Printf("Can not connect to db err: %v\n", err)
		}
		err = popdb.Open()
		if err != nil {
			fmt.Printf("Can not connect to db err: %v\n", err)
		}
		//pop.Debug = true
	}
}

func PopInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})
	for i := 0; i < b.N; i++ {
		m.Id = 0
		if err := popdb.Create(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PopInsertMulti(b *B) {
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
		if err := popdb.Create(&ms); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PopUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if err := popdb.Create(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	for i := 0; i < b.N; i++ {
		if err := popdb.Update(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PopRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if err := popdb.Create(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		if err := popdb.First(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PopReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			m.Id = 0
			if err := popdb.Create(m); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})

	for i := 0; i < b.N; i++ {
		var models []*Model
		if err := popdb.Where("id > ?", 0).Limit(b.L).All(&models); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
