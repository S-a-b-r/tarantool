package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/crud"
)

type Cache interface {
	Close()
	Subscriber(poolSize int, channel Channel, handler func(ch, p, m string))
	Publisher(channel Channel) chan<- string
}

type Session struct {
	conn *tarantool.Connection
	ctx  context.Context
	l    *zerolog.Logger
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

func Init(ctx context.Context, logger *zerolog.Logger, url string) Cache {
	l := logger.With().Str("address", url).Logger()

	dealer := GetDealer(url)
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		fmt.Println("Connection refused:", err)
		return nil
	}

	return &Session{}
}

func (s *Session) Subscriber(poolSize int, channel Channel, handler func(ch, p, m string)) {}

func (s *Session) Publisher(channel Channel) chan<- string {
	l := s.l.With().Interface("channel", channel).Logger()

	pubCh := make(chan string, 1)

	go s.publisher(&l, channel, pubCh)

	return pubCh
}

func (s *Session) subscriber(l *log.Logger, poolSize int, ch <-chan *tarantool.Message, h func(ch string, p string, m string)) {
	l.Println("start new subscriber")
	defer l.Println("subscriber done")

	s.handle = h

	for {
		select {
		case m, ok := <-ch:
			if !ok {
				l.Fatal("empty message received")
			}

			s.pool <- struct{}{}

			go s.handleMessage(s.pool, m)

		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Session) handleMessage(pool chan struct{}, m *tarantool.Message) {
	defer func() {
		<-pool
	}()

	fiber, err := fiber.New(m.Payload())
	if err != nil {
		log.Printf("error creating fiber: %v", err)
		return
	}

	s.handle(fiber.Channel(), fiber.Peer(), fiber.Message())
}

func example(conn tarantool.Connector) {

	// req := crud.MakeSelectRequest("bands").
	// 	Opts(crud.SelectOpts{
	// 		First: crud.MakeOptInt(2),
	// 	})

	req := crud.MakeGetRequest("bands").Key(4) // getReq

	ret := crud.Result{}
	if err := conn.Do(req).GetTyped(&ret); err != nil {
		fmt.Printf("Failed to execute request: %s", err)
		return
	}

	fmt.Println("Tuple selected by the primary key value:", ret.Rows)
}
