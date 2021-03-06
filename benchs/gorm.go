package benchs

import (
	"fmt"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var gormdb *gorm.DB

func init() {
	st := NewSuite("gorm")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, GormInsert)
		st.AddBenchmark("BulkInsert 100 row", 2000*ORM_MULTI, 0, GormInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, GormUpdate)
		st.AddBenchmark("Read", 2000*ORM_MULTI, 0, GormRead)
		st.AddBenchmark("MultiRead limit 2000", 2000*ORM_MULTI, 2000, GormReadSlice)

		conn, err := gorm.Open(postgres.Open(ORM_SOURCE), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})
		if err != nil {
			fmt.Println(err)
		}
		gormdb = conn
	}
}

func GormInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 0; i < b.N; i++ {
		m.Id = 0
		d := gormdb.Create(&m)
		if d.Error != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	}
}

func GormInsertMulti(b *B) {
	var ms []Model
	wrapExecute(b, func() {
		initDB()
		ms = make([]Model, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, *NewModel())
		}
	})
	for i := 0; i < b.N; i++ {
		ms1 := make([]Model, 100)
		copy(ms1, ms)
		if err := gormdb.Create(ms1).Error; err != nil {
			//fmt.Println(err)
			b.FailNow()
		}
	}
}

func GormUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		d := gormdb.Create(&m)
		if d.Error != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		d := gormdb.Save(&m)
		if d.Error != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	}
}

func GormRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		d := gormdb.Create(&m)
		if d.Error != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	})
	for i := 0; i < b.N; i++ {
		d := gormdb.Find(&m)
		if d.Error != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	}
}

func GormReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			m.Id = 0
			d := gormdb.Create(&m)
			if d.Error != nil {
				fmt.Println(d.Error)
				b.FailNow()
			}
		}
	})

	for i := 0; i < b.N; i++ {
		var models []*Model
		d := gormdb.Where("id > ?", 0).Limit(b.L).Find(&models)
		if d.Error != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	}
}
