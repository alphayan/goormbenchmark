package benchs

import (
	"fmt"

	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/mysql"
)

var godbcon *godb.DB

func init() {
	st := NewSuite("godb")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, GoDBInsert)
		st.AddBenchmark("BulkInsert 100 row", 500*ORM_MULTI, 0, GoDBInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, GoDBUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, 0, GoDBRead)
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, GoDBReadSlice)
		var err error
		godbcon, err = godb.Open(mysql.Adapter, ORM_SOURCE)
		if err != nil {
			fmt.Printf("conn err: %v\n", err)
		}
		//db.SetLogger(log.New(os.Stderr, "", 0))
	}
}

func GoDBInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 0; i < b.N; i++ {
		m.Id = 0
		if err := godbcon.Insert(m).Do(); err != nil {
			fmt.Printf("insert err: %v\n", err)
			b.FailNow()
		}
	}
}

func GoDBInsertMulti(b *B) {
	var ms []*Model
	wrapExecute(b, func() {
		initDB()
	})

	for i := 0; i < b.N; i++ {
		ms = make([]*Model, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, NewModel())
		}
		if err := godbcon.BulkInsert(&ms).Do(); err != nil {
			fmt.Printf("bulkinsert err: %v\n", err)
			b.FailNow()
		}
	}
}

func GoDBUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if err := godbcon.Insert(m).Do(); err != nil {
			fmt.Printf("insert before update err: %v\n", err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		if err := godbcon.Update(m).Do(); err != nil {
			fmt.Printf("update err: %v\n", err)
			b.FailNow()
		}
	}
}

func GoDBRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if err := godbcon.Insert(m).Do(); err != nil {
			fmt.Printf("insert before read err: %v\n", err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		if err := godbcon.Select(m).Do(); err != nil {
			fmt.Printf("read err: %v\n", err)
			b.FailNow()
		}
	}
}

func GoDBReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			m.Id = 0
			if err := godbcon.Insert(m).Do(); err != nil {
				fmt.Printf("insert before readslice err: %v\n", err)
				b.FailNow()
			}
		}
	})
	/*
		f, err := os.Create("/tmp/godb.cprof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	*/
	for i := 0; i < b.N; i++ {
		var models []*Model
		if err := godbcon.Select(&models).Where("id > ?", 0).Limit(b.L).Do(); err != nil {
			fmt.Printf("slice err: %v\n", err)
			b.FailNow()
		}
	}
	/*runtime.GC()
	memProfile, err := os.Create("/tmp/godb.mprof")
	if err != nil {
		log.Fatal(err)
	}
	defer memProfile.Close()
	//go tool pprof --alloc_space mem /tmp/godb.mprof
	if err := pprof.WriteHeapProfile(memProfile); err != nil {
		log.Fatal(err)
	}*/
}
