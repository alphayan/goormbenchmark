package benchs

import (
	"context"
	"fmt"

	"gitee.com/chunanyong/zorm"
)

var zormdb *zorm.BaseDao

func init() {
	st := NewSuite("zorm")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, ZormInsert)
		st.AddBenchmark("BulkInsert 100 row", 2000*ORM_MULTI, 0, ZormInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, ZormUpdate)
		st.AddBenchmark("Read", 2000*ORM_MULTI, 0, ZormRead)
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, ZormReadSlice)
		dataSourceConfig := zorm.DataSourceConfig{
			DSN:        ORM_SOURCE,
			DriverName: "mysql",
			DBType:     "mysql",
		}
		zorm.NewBaseDao(&dataSourceConfig)
	}
}

func ZormInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})
	for i := 0; i < b.N; i++ {
		m.Id = 0
		d := zorm.SaveStruct(context.Background(), m)
		if d != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	}
}

func ZormInsertMulti(b *B) {
	panic(fmt.Errorf("Don't support bulk insert"))
}

func ZormUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		d := zorm.SaveStruct(context.Background(), m)
		if d != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		d := zorm.UpdateStruct(context.Background(), m)
		if d != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	}
}

func ZormRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		d := zorm.SaveStruct(context.Background(), m)
		if d != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	})
	for i := 0; i < b.N; i++ {
		//查询Struct对象列表
		d := zorm.QueryStruct(context.Background(), zorm.NewSelectFinder(m.TableName()), m)
		if d != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	}
}

func ZormReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			m.Id = 0
			d := zorm.SaveStruct(context.Background(), m)
			if d != nil {
				fmt.Println(d.Error)
				b.FailNow()
			}
		}
	})
	for i := 0; i < b.N; i++ {
		var models []*Model
		d := zorm.QueryStructList(context.Background(), zorm.NewSelectFinder(m.TableName()), &models, zorm.NewPage())
		if d != nil {
			fmt.Println(d.Error)
			b.FailNow()
		}
	}
}
