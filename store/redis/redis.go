package redis

import (
	"github.com/micro/go-micro/config/options"
	"github.com/micro/go-micro/store"
	redis "gopkg.in/redis.v3"
)

type rkv struct {
	options.Options
	Client *redis.Client
}

func (r *rkv) Read(key string) (*store.Record, error) {
	val, err := r.Client.Get(key).Bytes()

	if err != nil && err == redis.Nil {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, store.ErrNotFound
	}

	d, err := r.Client.TTL(key).Result()
	if err != nil {
		return nil, err
	}

	return &store.Record{
		Key:    key,
		Value:  val,
		Expiry: d,
	}, nil
}

func (r *rkv) Delete(keys ...string) error {
	var err error
	for _, key := range keys {
		if err = r.Client.Del(key).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *rkv) Write(record *store.Record) error {
	return r.Client.Set(record.Key, record.Value, record.Expiry).Err()
}

func (r *rkv) Sync() ([]*store.Record, error) {
	keys, err := r.Client.Keys("*").Result()
	if err != nil {
		return nil, err
	}
	var vals []*store.Record
	for _, k := range keys {
		i, err := r.Read(k)
		if err != nil {
			return nil, err
		}
		vals = append(vals, i)
	}
	return vals, nil
}

func (r *rkv) String() string {
	return "redis"
}

func NewStore(opts ...options.Option) store.Store {
	options := options.NewOptions(opts...)

	var nodes []string

	if n, ok := options.Values().Get("store.nodes"); ok {
		nodes = n.([]string)
	}

	if len(nodes) == 0 {
		nodes = []string{"127.0.0.1:6379"}
	}

	return &rkv{
		Options: options,
		Client: redis.NewClient(&redis.Options{
			Addr:     nodes[0],
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}
