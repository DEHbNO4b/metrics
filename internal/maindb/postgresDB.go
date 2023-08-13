package maindb

import (
	"database/sql"
	"errors"

	"github.com/DEHbNO4b/metrics/internal/data"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

var createMetricsTable string = `CREATE TABLE IF NOT EXISTS metrics(
									id varchar(150) UNIQUE,
									type varchar(150),
									delta integer,
									value double precision
									);`

// var clearMetricsTable string = `delete from metrics;`

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB(dsn string) *PostgresDB {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log.Error("cannot open db", zap.Error(err))
		return &PostgresDB{}
	}
	db.Exec(createMetricsTable)
	return &PostgresDB{DB: db}
}
func (pdb *PostgresDB) Ping() error {
	if pdb.DB == nil {
		return errors.New("not connected to DB")
	}
	return pdb.DB.Ping()
}
func (pdb *PostgresDB) WriteMetrics(metrics []data.Metrics) error {
	if err := pdb.DB.Ping(); err != nil {
		return err
	}
	_, err := pdb.DB.Exec("DELETE FROM metrics")
	if err != nil {
		logger.Log.Error("unable to clear metric table ", zap.Error(err))
	}
	for _, metric := range metrics {
		pdb.Add(metric)
	}
	return nil
}
func (pdb *PostgresDB) ReadMetrics() ([]data.Metrics, error) {

	metrics := make([]data.Metrics, 0, 40)
	var (
		id    string
		mtype string
		value float64
		delta int64
	)
	if err := pdb.DB.Ping(); err != nil {
		return metrics, err
	}
	rows, err := pdb.DB.Query(`SELECT id,type,delta,value from metrics;`)
	if err != nil {
		logger.Log.Error("unable to get metrics from db", zap.Error(err))
		return metrics, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &mtype, &delta, &value)
		if err != nil {
			logger.Log.Error("unable to get metrics from db", zap.Error(err))
			return metrics, err
		}
		val := value
		del := delta
		metrics = append(metrics, data.Metrics{ID: id, MType: mtype, Delta: &del, Value: &val})
	}
	if err := rows.Err(); err != nil {
		logger.Log.Error("unable to get metrics from db", zap.Error(err))
		return metrics, err
	}
	return metrics, nil
}
func (pdb *PostgresDB) Add(metric data.Metrics) error {
	if err := pdb.DB.Ping(); err != nil {
		return err
	}
	_, err := pdb.DB.Exec(`insert into metrics (id,type,delta,value)	values($1,$2,$3,$4);`, metric.ID, metric.MType, *metric.Delta, *metric.Value)
	if err != nil {
		logger.Log.Error("cannot set gauge to db", zap.Error(err))
		return err
	}
	return nil
}

// func (pdb *PostgresDB) SetMetric(metric data.Metrics) error {
// 	switch metric.MType {
// 	case "gauge":
// 		_, err := pdb.DB.Exec(`insert into gauges (id,value)	values($1,$2);`, metric.ID, *metric.Value)
// 		if err != nil {
// 			logger.Log.Error("cannot set gauge to db", zap.Error(err))
// 			return err
// 		}
// 	case "counter":
// 		_, err := pdb.DB.Exec(`insert into counters (id,delta)	values($1,$2);`, metric.ID, *metric.Delta)
// 		if err != nil {
// 			logger.Log.Error("cannot set counter metric to db", zap.Error(err))
// 			return err
// 		}
// 	default:
// 		return errors.New("wrong metric type")
// 	}
// 	return nil
// }
// func (pdb *PostgresDB) GetMetrics() []string {
// 	m := make([]string, 0, 40)
// 	var id string
// 	var value float64
// 	var delta int64
// 	rows, err := pdb.DB.Query(`SELECT id,value from gauges;`)
// 	if err != nil {
// 		logger.Log.Error("unable to get metrics from db", zap.Error(err))
// 		return nil
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		err := rows.Scan(&id, &value)
// 		if err != nil {
// 			logger.Log.Error("unable to get metrics from db", zap.Error(err))
// 			return nil
// 		}
// 		m = append(m, id+":"+strconv.FormatFloat(value, 'f', -1, 64))
// 	}
// 	if err := rows.Err(); err != nil {
// 		logger.Log.Error("unable to get metrics from db", zap.Error(err))
// 		return nil
// 	}
// 	rowsC, err := pdb.DB.Query(`SELECT id,delta from counters;`)
// 	if err != nil {
// 		logger.Log.Error("unable to get metrics from db", zap.Error(err))
// 		return nil
// 	}
// 	defer rowsC.Close()
// 	for rowsC.Next() {
// 		err := rowsC.Scan(&id, &delta)
// 		if err != nil {
// 			logger.Log.Error("unable to get metrics from db", zap.Error(err))
// 			return nil
// 		}
// 		m = append(m, id+":"+strconv.FormatInt(delta, 10))
// 	}
// 	if err := rowsC.Err(); err != nil {
// 		logger.Log.Error("unable to get metrics from db", zap.Error(err))
// 		return nil
// 	}

// 	return m
// }

// func (pdb *PostgresDB) GetGauge(name string) (Gauge, error) {
// 	g := Gauge{}
// 	g.Name = name
// 	row := pdb.DB.QueryRow(`select value from gauges where id = $1;`, name)
// 	err := row.Scan(&g.Val)
// 	if err != nil {
// 		logger.Log.Error("cannot get gauge from db", zap.Error(err))
// 		return Gauge{}, err
// 	}
// 	return g, nil
// }
// func (pdb *PostgresDB) GetCounter(name string) (Counter, error) {
// 	c := Counter{}
// 	c.Name = name
// 	row := pdb.DB.QueryRow(`select delta from counters where id = $1;`, name)
// 	err := row.Scan(&c.Val)
// 	if err != nil {
// 		logger.Log.Error("cannot get gauge from db", zap.Error(err))
// 		return Counter{}, err
// 	}
// 	return c, nil
// }
