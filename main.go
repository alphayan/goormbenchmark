package main

import (
	"flag"
	"fmt"
	"goormbenchorm/benchs"
	mybenchs "goormbenchorm/mysqlbenchs"
	"math/rand"
	"runtime"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/lib/pq"
)

type ListOpts []string

func (opts *ListOpts) String() string {
	return fmt.Sprint(*opts)
}

func (opts *ListOpts) Set(value string) error {
	if value == "all" || strings.Index(" "+strings.Join(benchs.BrandNames, " ")+" ", " "+value+" ") != -1 {
	} else {
		return fmt.Errorf("wrong run name %s", value)
	}
	*opts = append(*opts, value)
	return nil
}

// Shuffle shuffles benchmark order
func (opts ListOpts) Shuffle() {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(opts); i++ {
		a := rd.Intn(len(opts))
		b := rd.Intn(len(opts))
		opts[a], opts[b] = opts[b], opts[a]
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var orms ListOpts
	var dbtype string
	var max_idle int
	var max_conn int
	var mysql_source string
	var postgresq_source string
	var multi int
	flag.IntVar(&max_idle, "max_idle", 20, "max idle conns")
	flag.IntVar(&max_conn, "max_conn", 200, "max open conns")
	flag.StringVar(&dbtype, "db type", "mysql", "mysql or postgresql")
	flag.StringVar(&postgresq_source, "psource", "host=127.0.0.1 port=5432 user=postgres password=root123456 dbname=test sslmode=disable", "postgres dsn source")
	flag.StringVar(&mysql_source, "msource", "root:root123456@(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local", "mysql dsn source")
	flag.IntVar(&multi, "multi", 1, "base query nums x multi")
	flag.Var(&orms, "orm", "orm name: all, "+strings.Join(benchs.BrandNames, ", "))
	flag.Parse()
	var all bool
	switch dbtype {
	case "postgresql":
		benchs.ORM_MAX_CONN = max_conn
		benchs.ORM_MAX_IDLE = max_idle
		benchs.ORM_MULTI = multi
		benchs.ORM_SOURCE = postgresq_source
	default:
		mybenchs.ORM_MAX_CONN = max_conn
		mybenchs.ORM_MAX_IDLE = max_idle
		mybenchs.ORM_MULTI = multi
		mybenchs.ORM_SOURCE = mysql_source
	}
	if len(orms) == 0 {
		all = true
	} else {
		for _, n := range orms {
			if n == "all" {
				all = true
			}
		}
	}

	if all {
		orms = benchs.BrandNames
	}

	orms.Shuffle()

	for _, n := range orms {
		fmt.Println(n)
		benchs.RunBenchmark(n)
	}

	fmt.Println("\nReports: \n")
	switch dbtype {
	case "postgresql":
		fmt.Print(benchs.MakeReport())
	default:
		fmt.Print(mybenchs.MakeReport())
	}

}
