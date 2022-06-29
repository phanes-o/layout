package traefik

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	etcd = &etcdRegistry{
		srvsRW: &sync.RWMutex{},
		srvs:   make(map[string]*etcdSrv),
	}
)

type etcdSrv struct {
	Config  *Config
	LeaseID clientv3.LeaseID
}

type etcdRegistry struct {
	cli *clientv3.Client

	srvsRW *sync.RWMutex
	srvs   map[string]*etcdSrv
}

func (r *etcdRegistry) init(addr string) error {
	var err error
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: time.Minute,
	})
	if err == nil {
		go r.loop()
	}
	return err
}

func (r *etcdRegistry) loop() {
	for {
		for _, srv := range r.srvs {
			if _, err := r.cli.KeepAliveOnce(context.Background(), srv.LeaseID); err != nil {
				if err == rpctypes.ErrLeaseNotFound {
					if srv.LeaseID, err = r.registerCore(srv.Config); err != nil {
						fmt.Println("(traefik) etcd registerCore errors: ", err)
					}
				} else {
					fmt.Println("(traefik) etcd KeepAliveOnce errors: ", err)
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func (r *etcdRegistry) register(conf *Config) error {
	r.srvsRW.RLock()
	_, b := r.srvs[conf.SrvName]
	r.srvsRW.RUnlock()
	if b {
		return nil
	}

	r.srvsRW.Lock()
	defer r.srvsRW.Unlock()

	var (
		err error
		lid clientv3.LeaseID
	)
	if lid, err = r.registerCore(conf); err == nil {
		r.srvs[conf.SrvName] = &etcdSrv{Config: conf, LeaseID: lid}
	}

	return err
}

func (r *etcdRegistry) registerCore(conf *Config) (clientv3.LeaseID, error) {
	var (
		err error
		kvs map[string]string
		ttl = 10 * time.Second
		lgr *clientv3.LeaseGrantResponse
	)
	if conf.TTL < ttl {
		conf.TTL = ttl
	}
	if lgr, err = r.cli.Grant(context.Background(), int64(conf.TTL.Seconds())); err != nil {
		return 0, err
	}
	if kvs, err = r.build(conf); err != nil {
		return 0, err
	}
	for k, v := range kvs {
		_, err = r.cli.Put(context.Background(), k, v, clientv3.WithLease(lgr.ID))
		if err != nil {
			return 0, err
		}
	}
	return lgr.ID, nil
}

func (r *etcdRegistry) build(conf *Config) (map[string]string, error) {
	conf.SrvName = strings.Replace(conf.SrvName, ".", "-", -1)
	var (
		t   string
		srv = conf.SrvName + "-srv"
	)

	switch conf.Type {
	case ReverseTypeHttp:
		if !strings.HasPrefix(conf.SrvAddr, "httpc") {
			conf.SrvAddr = "http://" + conf.SrvAddr
		}
		t = "httpc"
	case ReverseTypeH2c:
		if !strings.HasPrefix(conf.SrvAddr, "h2c") {
			conf.SrvAddr = "h2c://" + conf.SrvAddr
		}
		t = "httpc"
	case ReverseTypeTcp:
		t = "tcp"
	case ReverseTypeUdp:
		t = "udp"
	}

	dict := map[string]string{
		"traefik/enable": "true",
		fmt.Sprintf("traefik/%s/routers/%v/rule", t, conf.SrvName):    conf.Rule,
		fmt.Sprintf("traefik/%s/routers/%v/service", t, conf.SrvName): srv,
	}

	for i, ep := range conf.EndPoints {
		switch conf.Type {
		case ReverseTypeUdp:
			fmt.Println(t, conf.SrvName, i, ep)
			dict[fmt.Sprintf("traefik/%s/routers/%v/entrypoints/%v", t, conf.SrvName, i)] = ep
		case ReverseTypeTcp:
			dict[fmt.Sprintf("traefik/%s/routers/%v/entrypoints/%v", t, conf.SrvName, i)] = ep
		case ReverseTypeHttp, ReverseTypeH2c:
			dict[fmt.Sprintf("traefik/%s/routers/%v/entrypoints/%v", t, conf.SrvName, i)] = ep
		}
	}

	if len(conf.Prefix) > 0 {
		prefix := fmt.Sprintf("%v-prefix", conf.SrvName)
		dict[fmt.Sprintf("traefik/%s/routers/%v/middlewares/0", t, conf.SrvName)] = prefix
		dict[fmt.Sprintf("traefik/%s/middlewares/%v/stripPrefix/prefixes/0", t, prefix)] = conf.Prefix
	}

	urlPrefix := fmt.Sprintf("traefik/%s/services/%v/loadbalancer/servers/", t, srv)
	urlResp, err := r.cli.Get(context.Background(), urlPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	index := urlResp.Count
	for _, kv := range urlResp.Kvs {
		arr := strings.Split(string(kv.Key), "/")
		if len(arr) < 2 {
			continue
		}
		i, e := strconv.ParseInt(arr[len(arr)-2], 10, 64)
		if e != nil {
			continue
		}
		if i >= index {
			index = i + 1
		}
	}
	switch conf.Type {
	case ReverseTypeUdp:
		dict[fmt.Sprintf("traefik/udp/services/%v/loadbalancer/servers/%v/address", srv, index)] = conf.SrvAddr
		fmt.Println(dict)
	case ReverseTypeTcp:
		dict[fmt.Sprintf("traefik/tcp/services/%v/loadbalancer/servers/%v/address", srv, index)] = conf.SrvAddr
	case ReverseTypeHttp:
		dict[fmt.Sprintf("traefik/%s/services/%v/loadbalancer/servers/%v/url", t, srv, index)] = conf.SrvAddr
	case ReverseTypeH2c:
		dict[fmt.Sprintf("traefik/%s/services/%v/loadbalancer/servers/%v/url", t, srv, index)] = conf.SrvAddr
		dict[fmt.Sprintf("traefik/%s/services/%v/loadbalancer/healthcheck/scheme", t, srv)] = "h2c"
	}

	return dict, nil
}
