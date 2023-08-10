package maindb

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/DEHbNO4b/metrics/internal/data"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

var CreateGauges string = `CREATE TABLE IF NOT EXISTS gauges(
									id varchar(150) UNIQUE,
									value double precision
									);`
var CreateCounters string = `CREATE TABLE IF NOT EXISTS gauges(
									id varchar(150) UNIQUE,
									delta integer
									);`

type PostgresDB struct {
	Config data.StoreConfig
	DB     *sql.DB
}

func NewPostgresDB(dsn string) *PostgresDB {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log.Panic("cannot open db", zap.Error(err))
		return nil
	}
	_, err = db.Exec(CreateGauges)
	if err != nil {
		logger.Log.Panic("cannot create table in db", zap.Error(err))
	}
	_, err = db.Exec(CreateCounters)
	if err != nil {
		logger.Log.Panic("cannot create table in  db", zap.Error(err))
	}
	return &PostgresDB{DB: db}
}
func (mdb *PostgresDB) Ping() error {
	return mdb.DB.Ping()
}
func (mdb *PostgresDB) SetMetric(metric data.Metrics) error {
	switch metric.MType {
	case "gauge":
		_, err := mdb.DB.Exec(`insert into gauges (id,value)	values($1,$2);`, metric.ID, *metric.Value)
		if err != nil {
			logger.Log.Error("cannot set gauge to db", zap.Error(err))
			return err
		}
	case "counter":
		_, err := mdb.DB.Exec(`insert into counters (id,delta)	values($1,$2);`, metric.ID, *metric.Delta)
		if err != nil {
			logger.Log.Error("cannot set counter metric to db", zap.Error(err))
			return err
		}
	default:
		return errors.New("wrong metric type")
	}
	return nil
}
func (mdb *PostgresDB) GetMetrics() []string {
	m := make([]string, 0, 40)
	var id string
	var value float64
	var delta int64
	rows, err := mdb.DB.Query(`SELECT id,value from gauges;`)
	if err != nil {
		logger.Log.Error("unable to get metrics from db", zap.Error(err))
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id, &value)
		m = append(m, id+":"+strconv.FormatFloat(value, 'f', -1, 64))
	}
	rows, err = mdb.DB.Query(`SELECT id,delta from counters;`)
	if err != nil {
		logger.Log.Error("unable to get metrics from db", zap.Error(err))
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id, &delta)
		m = append(m, id+":"+strconv.FormatInt(delta, 10))
	}

	return m
}
func (mdb *PostgresDB) SetGauge(g data.Gauge) error {
	// ms.Gauges[g.Name] = g.Val
	return nil
}
func (mdb *PostgresDB) SetCounter(c data.Counter) error {
	// ms.Counters[c.Name] = ms.Counters[c.Name] + c.Val
	return nil
}
func (mdb *PostgresDB) GetGauge(name string) (data.Gauge, error) {
	g := data.Gauge{}
	g.Name = name
	row := mdb.DB.QueryRow(`select value from gauges where id = $1;`, name)
	err := row.Scan(&g.Val)
	if err != nil {
		logger.Log.Error("cannot get gauge from db", zap.Error(err))
		return data.Gauge{}, err
	}
	return g, nil
}
func (mdb *PostgresDB) GetCounter(name string) (data.Counter, error) {
	c := data.Counter{}
	c.Name = name
	row := mdb.DB.QueryRow(`select delta from counters where id = $1;`, name)
	err := row.Scan(&c.Val)
	if err != nil {
		logger.Log.Error("cannot get gauge from db", zap.Error(err))
		return data.Counter{}, err
	}
	return c, nil
}
