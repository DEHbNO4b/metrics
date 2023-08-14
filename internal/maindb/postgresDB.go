package maindb

import (
	"database/sql"

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
	if dsn == "" {
		return nil
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log.Error("cannot open db", zap.Error(err))
		return nil
	}
	db.Exec(createMetricsTable)
	return &PostgresDB{DB: db}
}
func (pdb *PostgresDB) Ping() bool {
	if pdb.DB == nil {
		logger.Log.Error("database object is nil")
		return false
	}
	err := pdb.DB.Ping()
	if err != nil {
		logger.Log.Error(err.Error())
		return false
	}
	return true
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
