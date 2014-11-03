package cache

import (
	"errors"
	"io"
	"log"
	"os"
	"reflect"
	"sync"

	"github.com/dustin/gomemcached"
	memcache "github.com/dustin/gomemcached/client"
	"github.com/ugorji/go/codec"
)

var (
	CACHE CacheInterface
)

type CacheInterface interface {
	Set(uint16, string, int, int, []byte) (*gomemcached.MCResponse, error)
	Get(uint16, string) (*gomemcached.MCResponse, error)
	Auth(string, string) (*gomemcached.MCResponse, error)
	IsHealthy() bool
}

func Set(key string, data interface{}) {
	var (
		value []byte
		mh    codec.MsgpackHandle
		w     io.Writer
	)

	mh.MapType = reflect.TypeOf(data)
	mh.SliceType = reflect.TypeOf(data)

	enc := codec.NewEncoder(w, &mh)
	enc = codec.NewEncoderBytes(&value, &mh)
	err := enc.Encode(data)

	if err == nil {
		_, err = CACHE.Set(0, key, 0, 0, value)
		log.Println("Caching", key, "--", mh.MapType, " -->", data)
		if err != nil {
			log.Println("CACHE ERROR:", err)
		}
	} else {
		log.Println("ENCODE ERROR:", err)
	}
}

func Get(key string, as interface{}) (interface{}, error) {
	var (
		data interface{}
		mh   codec.MsgpackHandle
		r    io.Reader
	)

	mh.MapType = reflect.TypeOf(as)
	mh.SliceType = reflect.TypeOf(as)

	cached, err := CACHE.Get(0, key)
	if err != nil {
		return nil, err
	}

	dec := codec.NewDecoder(r, &mh)
	dec = codec.NewDecoderBytes(cached.Body, &mh)
	err = dec.Decode(&data)

	log.Println("Getting", key, "--", mh.MapType, " -->", data)

	return data, err
}

func init() {
	memcachedUrl := os.Getenv("MEMCACHEDCLOUD_SERVERS")
	memcachedPassword := os.Getenv("MEMCACHEDCLOUD_PASSWORD")
	memcachedUsername := os.Getenv("MEMCACHEDCLOUD_USERNAME")

	var memcachedErr error

	CACHE, memcachedErr = memcache.Connect("tcp", memcachedUrl)

	if memcachedErr != nil {
		log.Println("in APP CACHE")
		CACHE = inAppCache{}
	}

	if memcachedPassword != "" {
		CACHE.Auth(memcachedUsername, memcachedPassword)
	}

	log.Println("CACHE Healthy", CACHE.IsHealthy())
	log.Println(CACHE.Set(0, "test", 0, 5, []byte("Testing...")))
	log.Println(CACHE.Get(0, "test"))
}

var DATA = make(map[string][]byte, 0)

type inAppCache struct {
	sync.Mutex
}

func (c inAppCache) Set(_ uint16, key string, _ int, _ int, value []byte) (*gomemcached.MCResponse, error) {
	c.Lock()
	if DATA == nil {
		DATA = make(map[string][]byte, 0)
	}
	DATA[key] = value
	defer c.Unlock()

	return &gomemcached.MCResponse{Body: value}, nil
}

func (c inAppCache) Get(_ uint16, key string) (*gomemcached.MCResponse, error) {
	d := DATA[key]
	if d == nil {
		return nil, errors.New("Not found")
	}

	return &gomemcached.MCResponse{Body: d}, nil
}

func (c inAppCache) Auth(string, string) (*gomemcached.MCResponse, error) {
	return nil, nil
}

func (c inAppCache) IsHealthy() bool {
	return true
}
