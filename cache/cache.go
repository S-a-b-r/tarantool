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
	// Subscriber(poolSize int, channel Channel, handler func(m string))
	// Publisher(channel Channel) chan<- string
	HGet(ctx context.Context, hash, field string) *StringCmd // работает
	Get(ctx context.Context, hash string) *StringCmd         // работает
	MGet(ctx context.Context, keys ...string) *SliceCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd // работает
	Keys(ctx context.Context, pattern string) *StringSliceCmd
	Del(ctx context.Context, keys ...string) *IntCmd
}

type Session struct {
	conn *tarantool.Connection
	ctx  context.Context
	l    *zerolog.Logger
}

func Init(ctx context.Context, logger *zerolog.Logger, url string) (Cache, error) {
	l := logger.With().Str("address", url).Logger()

	dialer := getDialer(url)
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		l.Error().Err(err).Msg("Connection refused")
		return &Session{}, err
	}

	_, err = conn.Do(tarantool.NewPingRequest()).Get()
	if err != nil {
		l.Error().Err(err).Msg("ping error")
		return &Session{}, err
	}

	return &Session{
		conn: conn,
		ctx:  ctx,
		l:    &l,
	}, nil
}

// Subscriber handler (m), m - префиксы Хеша, например(cmd, unitInfo и тд)
// func (s *Session) Subscriber(poolSize int, channel Channel, handler func(m string)) {
// 	l := s.l.With().Str("channel", string(channel)).Logger()
//
// 	stm, err := s.conn.NewStream()
// 	if err != nil {
// 		l.Error().Err(err).Msg("")
// 	}
// 	stm.Conn.Do(tarantool.NewPingRequest())
//
// 	// callback := func(event tarantool.WatchEvent) {
// 	// 	fmt.Printf("event connection: %s\n", event.Conn.Addr())
// 	// 	fmt.Printf("event key: %s\n", event.Key)
// 	// 	fmt.Printf("event value: %v\n", event.Value)
// 	// }
//
// 	// watcher, err := s.conn.NewWatcher(string(channel), callback)
// 	// if err != nil {
// 	// 	fmt.Printf("Failed to connect watcher: %s\n", err)
// 	// 	return
// 	// }
// 	// defer watcher.Unregister()
// }
//
// func (s *Session) subscriber(l *zerolog.Logger, poolSize int, ch chan *queue.Task, h func(m string), q queue.Queue) {
// 	l.Println("start new subscriber")
// 	defer l.Println("subscriber done")
//
// 	pool := make(chan struct{}, poolSize)
//
// 	go func() {
// 		for {
// 			msg, err := q.Take()
// 			if err != nil {
// 				l.Fatal().Err(err).Msg("error taking message from queue")
// 				return
// 			}
// 			ch <- msg
// 		}
// 	}()
//
// 	for {
// 		select {
// 		case m, ok := <-ch:
// 			if !ok {
// 				l.Fatal().Msg("empty message received")
// 			}
//
// 			pool <- struct{}{}
//
// 			switch mm := m.Data().(type) {
// 			case string:
// 				go s.handle(pool, mm, h)
// 			}
//
// 		case <-s.ctx.Done():
// 			return
// 		}
// 	}
// }
//
// func (s *Session) handle(pool chan struct{}, prefixHash string, h func(m string)) {
// 	h(prefixHash)
// 	<-pool
// }
//
// func (s *Session) Publisher(channel Channel) chan<- string {
// 	// l := s.l.With().Interface("channel", channel).Logger()
//
// 	pubCh := make(chan string, 1)
//
// 	// go s.publisher(&l, channel, pubCh)
//
// 	return pubCh
// }

func (s *Session) HGet(ctx context.Context, hash, field string) *StringCmd {
	req := crud.MakeGetRequest("cache").Context(ctx).Key(hash).Opts(crud.GetOpts{
		Fields: crud.MakeOptTuple([]interface{}{"hash_table"}),
	})
	ret := crud.Result{}

	if err := s.conn.Do(req).GetTyped(&ret); err != nil {
		return NewStringCmd("", fmt.Errorf("failed to execute request: %w", err))
	}

	return NewStringCmd(getFieldOnHashMap(ret, field))
}

func (s *Session) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd {
	var object = crud.MapObject{
		"hash":       key,
		"hash_table": value,
	}

	opts := crud.SimpleOperationObjectOpts{
		Timeout: crud.MakeOptFloat64(expiration.Seconds()),
	}

	req := crud.MakeInsertObjectRequest("cache").Context(ctx).Object(object).Opts(opts)

	res, err := s.conn.Do(req).Get()
	if err != nil {
		return NewStatusCmd("", fmt.Errorf("failed to execute request: %w", err))
	}

	fmt.Println(res)
	return NewStatusCmd("success", nil)
}

func (s *Session) Get(ctx context.Context, hash string) *StringCmd {
	req := crud.MakeGetRequest("cache").Context(ctx).Key(hash)
	ret := crud.Result{}

	if err := s.conn.Do(req).GetTyped(&ret); err != nil {
		return NewStringCmd("", fmt.Errorf("failed to execute request: %w", err))
	}

	return NewStringCmd(getValue(ret))
}

// MGet Лучше переписать на lua функции
func (s *Session) MGet(ctx context.Context, keys ...string) *SliceCmd {
	res := make([]interface{}, 0)
	for _, key := range keys {
		val, err := s.Get(ctx, key).Result()
		if err != nil {
			s.l.Warn().Err(err).Str("key", key).Msg("failed get data with key")
			continue
		}

		res = append(res, val)
	}

	return NewSliceCmd(res, nil)
}

func (s *Session) Keys(ctx context.Context, pattern string) *StringSliceCmd {
	req := tarantool.NewCallRequest("get_all_hashes_by_pattern").Args([]interface{}{[]string{pattern}}).Context(ctx)

	data, err := s.conn.Do(req).Get()
	if err != nil {
		return NewStringSliceCmd([]string{}, fmt.Errorf("failed to execute request: %w", err))
	}

	return NewStringSliceCmd(getStringSliceRes(data))
}

func (s *Session) Del(ctx context.Context, keys ...string) *IntCmd {
	countDeleted := 0

	for _, key := range keys {
		req := crud.MakeDeleteRequest("cache").Context(ctx).Key(key)

		if _, err := s.conn.Do(req).Get(); err != nil {
			return NewIntCmd(countDeleted, fmt.Errorf("failed to execute request: %w", err))
		}
		countDeleted++
	}

	return NewIntCmd(countDeleted, nil)
}

func (s *Session) Close() {
	if s == nil || s.conn == nil {
		return
	}

	if err := s.conn.Close(); err != nil {
		s.l.Error().Err(err).Msg("failed to close connection")
	}
}
