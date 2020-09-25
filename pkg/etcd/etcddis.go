package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	log "github.com/sirupsen/logrus"
	"github.com/woodylan/go-websocket/pkg/setting"
	"time"
)

type ClientDis struct {
	client *clientv3.Client
}

func NewClientDis(addr []string) (*ClientDis, error) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}
	if client, err := clientv3.New(conf); err == nil {
		return &ClientDis{
			client: client,
		}, nil
	} else {
		return nil, err
	}
}

func (this *ClientDis) GetService(prefix string) ([]string, error) {
	resp, err := this.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	addrs := this.extractAddrs(resp)

	go this.watcher(prefix)
	return addrs, nil
}

func (this *ClientDis) watcher(prefix string) {
	rch := this.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				this.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				this.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

func (this *ClientDis) extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			this.SetServiceList(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
			addrs = append(addrs, string(v))
		}
	}
	return addrs
}

func (this *ClientDis) SetServiceList(key, val string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()
	setting.GlobalSetting.ServerList[key] = val
	log.Info("发现服务：", key, " 地址:", val)
}

func (this *ClientDis) DelServiceList(key string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()
	delete(setting.GlobalSetting.ServerList, key)
	log.Println("服务下线:", key)
}
