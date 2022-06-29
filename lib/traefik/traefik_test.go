package traefik

import (
	"testing"
	"time"
)

func init() {
	if e := Init("127.0.0.1:2379"); e != nil {
		panic(e)
	}
}

func TestRegister(t *testing.T) {
	conf := &Config{
		SrvName:   "myTest",
		SrvAddr:   "127.0.0.1:6789",
		Rule:      "Host(`abc.test.com`) || PathPrefix(`/abc`)",
		Prefix:    "/test",
		EndPoints: []string{"httpc"},
	}
	e := Register(conf)
	if e != nil {
		t.Fatal(e)
	}
	t.Log("success")

	<-time.After(10 * time.Minute)
}
