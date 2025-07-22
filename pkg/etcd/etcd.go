package etcd

import (
	"context"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var (
	client *clientv3.Client
	once   sync.Once
)
type PutOption func(*clientv3.Op)
type GetOption func(*clientv3.Op)	
type EtcdOptions struct {
	Endpoints   []string
	DialTimeout time.Duration
	Username    string
	Password    string
}
type Lock interface {
	Key() string
	Release(context.Context) error
}
type EtcdService interface {
	Put(ctx context.Context, key, val string, opts ...PutOption) error
	Get(ctx context.Context, key string, opts ...GetOption) (string, error)
	Delete(ctx context.Context, key string) error
	Watch(ctx context.Context, key string, onChange func(string))
	Lock(ctx context.Context, key string, ttl int) (Lock, error)
	Unlock(lock Lock) error
	Close() error
}
type EtcdServiceImpl struct {
	client *clientv3.Client
}

// 初始化 etcd 客户端
func InitEtcd(cfg EtcdOptions) error {
	var err error
	once.Do(func() {
		client, err = clientv3.New(clientv3.Config{
			Endpoints:   cfg.Endpoints,
			DialTimeout: cfg.DialTimeout,
			Username:    cfg.Username,
			Password:    cfg.Password,
		})
	})
	return err
}

func GetEtcdInstance() *EtcdServiceImpl {
	return &EtcdServiceImpl{client: client}
}

func (e *EtcdServiceImpl) Put(ctx context.Context, key, val string, opts ...PutOption) error {
	_, err := e.client.Put(ctx, key, val)
	return err
}

func (e *EtcdServiceImpl) Get(ctx context.Context, key string, opts ...GetOption) (string, error) {
	resp, err := e.client.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

func (e *EtcdServiceImpl) Delete(ctx context.Context, key string) error {
	_, err := e.client.Delete(ctx, key)
	return err
}

func (e *EtcdServiceImpl) Close() error {
	return e.client.Close()
}

func (e *EtcdServiceImpl) Watch(ctx context.Context, key string, onChange func(string)) {
	rch := e.client.Watch(ctx, key)
	go func() {
		for wresp := range rch {
			for _, ev := range wresp.Events {
				onChange(string(ev.Kv.Value))
			}
		}
	}()
}

func (e *EtcdServiceImpl) Lock(ctx context.Context, key string, ttl int) (*concurrency.Mutex, *concurrency.Session, error) {
	session, err := concurrency.NewSession(e.client, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, nil, err
	}
	mutex := concurrency.NewMutex(session, key)
	if err := mutex.Lock(ctx); err != nil {
		return nil, nil, err
	}
	return mutex, session, nil
}

func (e *EtcdServiceImpl) Unlock(mutex *concurrency.Mutex, session *concurrency.Session) error {
	defer session.Close()
	return mutex.Unlock(context.TODO())
}