package pools

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

// Pool represents a connection pool for a database.
type PGXPool struct {
	mu        sync.Mutex
	conns     chan *pgx.Conn
	connCount int
	maxConns  int
	config    pgx.ConnConfig
	ctx       context.Context
}

// NewPool creates a new connection pool for a database.
func NewPool(maxConns int, config pgx.ConnConfig) *PGXPool {
	return &PGXPool{
		conns:     make(chan *pgx.Conn, maxConns),
		maxConns:  maxConns,
		config:    config,
		ctx:       context.Background(),
		connCount: 0,
	}
}

// Get retrieves a connection from the pool.
func (p *PGXPool) Get() (*pgx.Conn, error) {
	select {
	case conn := <-p.conns:
		return conn, nil
	default:
		if p.connCount < p.maxConns {
			conn, err := pgx.Connect(p.ctx, p.config.ConnString())
			if err != nil {
				return nil, err
			}
			p.connCount++
			return conn, nil
		}
		return nil, errors.New("pool is full")
	}
}

// Put returns a connection to the pool.
func (p *PGXPool) Put(conn *pgx.Conn) {
	select {
	case p.conns <- conn:
	default:
		conn.Close(context.Background())
		p.connCount--
	}
}

// Close closes all connections in the pool.
func (p *PGXPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		select {
		case conn := <-p.conns:
			conn.Close(context.Background())
			p.connCount--
		default:
			return
		}
	}
}

// Set the underlaying pgx connection timeout
func (p *PGXPool) SetConnTimeout(d time.Duration) {
	p.config.ConnectTimeout = d
}

// IdleConns returns the number of idle connections in the pool.
func (p *PGXPool) IdleConns() int {
	return len(p.conns)
}

// SetIdleTimeout sets the idle timeout for connections in the pool.
func (p *PGXPool) SetIdleTimeout(d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Cancel the previous context to prevent a context leak
	if p.ctx != nil {
		p.ctx.Done()
	}

	p.ctx, _ = context.WithTimeout(p.ctx, d)
	for {
		select {
		case conn := <-p.conns:
			if p.ctx.Err() != nil {
				conn.Close(context.Background())
				p.connCount--
			} else {
				p.conns <- conn
			}
		default:
			return
		}
	}
}
