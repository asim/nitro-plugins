package redis

import (
	"github.com/micro/go-micro/store"
	redis "gopkg.in/redis.v3"
)

type rkv struct {
	options store.Options
	Client  *redis.Client
}

func (r *rkv) Read(keys ...string) ([]*store.Record, error) {
	records := make([]*store.Record, 0, len(keys))

	for _, key := range keys {
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

		records = append(records, &store.Record{
			Key:    key,
			Value:  val,
			Expiry: d,
		})
	}

	return records, nil
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

func (r *rkv) Write(records ...*store.Record) error {
	var err error

	for _, record := range records {
		err = r.Client.Set(record.Key, record.Value, record.Expiry).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *rkv) List() ([]*store.Record, error) {
	keys, err := r.Client.Keys("*").Result()
	if err != nil {
		return nil, err
	}

	vals := make([]*store.Record, 0, len(keys))
	for _, k := range keys {
		i, err := r.Read(k)
		if err != nil {
			return nil, err
		}
		vals = append(vals, i...)
	}
	return vals, nil
}

func (r *rkv) String() string {
	return "redis"
}

func NewStore(opts ...store.Option) store.Store {
	var options store.Options
	for _, o := range opts {
		o(&options)
	}

	nodes := options.Nodes

	if len(nodes) == 0 {
		nodes = []string{"127.0.0.1:6379"}
	}

	return &rkv{
		options: options,
		Client: redis.NewClient(&redis.Options{
			Addr:     nodes[0],
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}
