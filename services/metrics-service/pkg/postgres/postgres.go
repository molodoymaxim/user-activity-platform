package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type PostgreTx interface {
	Transact(ctxParent context.Context, txFunc func(context.Context, TxPG) error, commitCallback func()) (err error)
}

type Postgre interface {
	NewPoolConfig(maxConn int, connIdleTime, connLifeTime time.Duration) error                                       // Создание конфигурации пула
	ConnectionPool(ctx context.Context) error                                                                        // Подключаемся с помощью пула к Postgres
	GetSQL(sqlFunc func(db *sql.DB) error) error                                                                     // Выполнение функции от имени драйвера sql.DB
	Ping(ctx context.Context) error                                                                                  // Проверяем соединение
	Transact(ctxParent context.Context, txFunc func(context.Context, TxPG) error, commitCallback func()) (err error) // Обработчик транзакций
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)                               // Exec запрос
	Close()                                                                                                          // Закрытие соединения
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)                                       // Query запрос
	QueryRow(ctxParent context.Context, sql string, arguments ...any) pgx.Row                                        // QueryRow запрос
	GetPostgreTx() PostgreTx
}

type postgres struct {
	conn         *pgxpool.Pool
	connStr      string
	poolConfig   *pgxpool.Config
	queryTimeout time.Duration
}

func New(cfg *Config) Postgre {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.PostgresQueryTimeout)

	return &postgres{
		connStr:      connStr,
		queryTimeout: time.Duration(cfg.PostgresQueryTimeout) * time.Second,
	}
}

func (p *postgres) NewPoolConfig(maxConn int, connIdleTime, connLifeTime time.Duration) error {
	// Создание конфигурации пула
	poolConfig, err := pgxpool.ParseConfig(p.connStr)
	if err != nil {
		return err
	}

	// Проверка
	cpu := runtime.NumCPU()
	if maxConn > cpu {
		maxConn = cpu
	}

	poolConfig.MaxConns = int32(maxConn)
	poolConfig.MaxConnIdleTime = connIdleTime
	poolConfig.MaxConnLifetime = connLifeTime
	p.poolConfig = poolConfig
	return nil
}

func (p *postgres) ConnectionPool(ctx context.Context) error {
	conn, err := pgxpool.NewWithConfig(ctx, p.poolConfig)
	if err != nil {
		return err
	}
	p.conn = conn
	err = p.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgres) Ping(ctx context.Context) error {
	return p.conn.Ping(ctx)
}

func (p *postgres) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return p.conn.Exec(ctx, sql, arguments...)
}

func (p *postgres) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	return p.conn.Query(ctx, sql, arguments...)
}

func (p *postgres) QueryRow(ctxParent context.Context, sql string, arguments ...any) pgx.Row {
	ctx, cancel := context.WithTimeout(ctxParent, p.queryTimeout)
	defer cancel()
	return p.conn.QueryRow(ctx, sql, arguments...)
}

func (p *postgres) Close() {
	p.conn.Close()
}

func (p *postgres) Transact(ctxParent context.Context, txFunc func(context.Context, TxPG) error, commitCallback func()) (errResp error) {
	ctx, cancel := context.WithTimeout(ctxParent, p.queryTimeout)
	defer cancel()

	tx, err := p.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			//log.Printf("Panic Transact: %+v / Stacktrace: %+v\n", p, string(debug.Stack()))
			if errTx := tx.Rollback(ctx); errTx != nil {
				//log.Printf("Tx rollback panic: %v\n", errTx)
				errResp = fmt.Errorf("tx rollback panic: %v", errTx)
			} else {
				errResp = fmt.Errorf("panic transact: %v", p)
				//log.Printf("Rollback Success: %+v\n", errResp)
			}
		} else if errResp != nil {
			//log.Printf("Error Transact: %+v\n", errResp)
			if errTx := tx.Rollback(ctx); errTx != nil {
				//log.Printf("Tx rollback: %v\n", errTx)
				errResp = fmt.Errorf("tx rollback: %v", errTx)
			}
		} else {
			if errTx := tx.Commit(ctx); errTx != nil {
				//log.Printf("Tx commit: %v\n", errTx)
				errResp = fmt.Errorf("tx commit: %v", errTx)
				if commitCallback != nil {
					commitCallback()
				}
			}
		}
	}()

	errResp = txFunc(ctx, tx)

	return errResp
}

func (p *postgres) GetSQL(sqlFunc func(db *sql.DB) error) error {
	return sqlFunc(stdlib.OpenDBFromPool(p.conn))
}

func (p *postgres) GetPostgreTx() PostgreTx {
	return p
}
