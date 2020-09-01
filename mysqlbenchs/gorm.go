package benchs

import (
	"fmt"

	"gorm.io/driver/mysql"
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
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, GormReadSlice)

		conn, err := gorm.Open(mysql.Open(ORM_SOURCE), &gorm.Config{SkipDefaultTransaction: true,
			PrepareStmt: true})
		if err != nil {
			fmt.Println(err)
			return
		}
		db, err := conn.DB()
		if err != nil {
			fmt.Println(err)
		}
		db.SetMaxIdleConns(ORM_MAX_IDLE)
		db.SetMaxOpenConns(ORM_MAX_CONN)
		gormdb = conn

	}
}

func GormInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Id = 0
		d := gormdb.Create(m)
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
		d := gormdb.Create(m)
		if d.Error != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := gormdb.Model(m).Updates(m)
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
		d := gormdb.Create(m)
		if d.Error != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := gormdb.Find(m)
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
			d := gormdb.Create(m)
			if d.Error != nil {
				fmt.Println(d.Error)
				b.FailNow()
			}
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var models []*Model
		d := gormdb.Where("id > ?", 0).Order("id asc").Limit(b.L).Find(&models)
		if d.Error != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	}

}
