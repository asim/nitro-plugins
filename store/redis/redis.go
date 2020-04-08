package redis

import (
	"github.com/micro/go-micro/v2/store"
	redis "gopkg.in/redis.v3"
)

type rkv struct {
	options store.Options
	Client  *redis.Client
}

func (r *rkv) Init(...store.Option) error {
	return nil
}

func (r *rkv) Close() error {
	return r.Client.Close()
}

func (r *rkv) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	records := make([]*store.Record, 0, 1)

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

	return records, nil
}

func (r *rkv) Delete(key string, opts ...store.DeleteOption) error {
	return r.Client.Del(key).Err()
}

func (r *rkv) Write(record *store.Record, opts ...store.WriteOption) error {
	return r.Client.Set(record.Key, record.Value, record.Expiry).Err()
}

func (r *rkv) List(opts ...store.ListOption) ([]string, error) {
	keys, err := r.Client.Keys("*").Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *rkv) Options() store.Options {
	return r.options
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
