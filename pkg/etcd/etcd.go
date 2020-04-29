package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	log "github.com/sirupsen/logrus"
	"go-websocket/define"
	"sync"
	"time"
)

var etcdKvClient *clientv3.Client
var mu sync.Mutex

func GetInstance() *clientv3.Client {
	if etcdKvClient == nil {
		if client, err := clientv3.New(clientv3.Config{
			Endpoints:   define.EtcdEndpoints,
			DialTimeout: 5 * time.Second,
		}); err != nil {
			log.Error(err)
			return nil
		} else {
			//创建时才加锁
			mu.Lock()
			defer mu.Unlock()
			etcdKvClient = client
			return etcdKvClient
		}

	}
	return etcdKvClient
}

func Put(key, value string) error {
	_, err := GetInstance().Put(context.Background(), key, value)
	return err
}

func Get(key string) (resp *clientv3.GetResponse, err error) {
	resp, err = GetInstance().Get(context.Background(), key)
	return resp, err
}
