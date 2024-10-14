package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/queue"
)

type Cache interface {
	Close()
	Subscriber(poolSize int, channel Channel, handler func(m string))
	Publisher(channel Channel) chan<- string
}

type Session struct {
	conn *tarantool.Connection
	ctx  context.Context
	l    *zerolog.Logger
}

func Init(ctx context.Context, logger *zerolog.Logger, url string) Cache {
	l := logger.With().Str("address", url).Logger()

	dialer := GetDialer(url)
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		l.Error().Err(err).Msg("Connection refused")
		return nil
	}

	data, err := conn.Do(tarantool.NewPingRequest()).Get()
	if err != nil {
		l.Error().Err(err).Msg("ping error")
		return nil
	}
	fmt.Println(data)

	return &Session{
		conn: conn,
		ctx:  ctx,
		l:    &l,
	}
}

// Subscriber handler (m), m - префиксы Хеша, например(cmd, unitInfo и тд)
func (s *Session) Subscriber(poolSize int, channel Channel, handler func(m string)) {
	l := s.l.With().Str("channel", string(channel)).Logger()

	stm, err := s.conn.NewStream()
	if err != nil {
		l.Error().Err(err).Msg("")
	}
	stm.Conn.Do(tarantool.NewPingRequest())

	// callback := func(event tarantool.WatchEvent) {
	// 	fmt.Printf("event connection: %s\n", event.Conn.Addr())
	// 	fmt.Printf("event key: %s\n", event.Key)
	// 	fmt.Printf("event value: %v\n", event.Value)
	// }

	// watcher, err := s.conn.NewWatcher(string(channel), callback)
	// if err != nil {
	// 	fmt.Printf("Failed to connect watcher: %s\n", err)
	// 	return
	// }
	// defer watcher.Unregister()
}

func (s *Session) subscriber(l *zerolog.Logger, poolSize int, ch chan *queue.Task, h func(m string), q queue.Queue) {
	l.Println("start new subscriber")
	defer l.Println("subscriber done")

	pool := make(chan struct{}, poolSize)

	go func() {
		for {
			msg, err := q.Take()
			if err != nil {
				l.Fatal().Err(err).Msg("error taking message from queue")
				return
			}
			ch <- msg
		}
	}()

	for {
		select {
		case m, ok := <-ch:
			if !ok {
				l.Fatal().Msg("empty message received")
			}

			pool <- struct{}{}

			switch mm := m.Data().(type) {
			case string:
				go s.handle(pool, mm, h)
			}

		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Session) handle(pool chan struct{}, prefixHash string, h func(m string)) {
	h(prefixHash)
	<-pool
}

func (s *Session) Publisher(channel Channel) chan<- string {
	// l := s.l.With().Interface("channel", channel).Logger()

	pubCh := make(chan string, 1)

	// go s.publisher(&l, channel, pubCh)

	return pubCh
}

func (s *Session) HGet(ctx context.Context, hash, field string) {
	// get by hash
	// hash = unitInfo:hash
}

func (s *Session) Close() {
	if s == nil || s.conn == nil {
		return
	}

	if err := s.conn.Close(); err != nil {
		s.l.Error().Err(err).Msg("failed to close connection")
	}
}
