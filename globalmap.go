package main

import (
	"context"
	"fmt"
	"log"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

func GetGlobalMap(key string) (value string, err error) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// reseting
	if len(resetLock.IsOwner().Key) > 1 {
		return "", RefusedByLockError
	}

	// create a sessions to aqcuire a lock
	s, err := concurrency.NewSession(cli, concurrency.WithTTL(10))
	if err != nil {
		return "", err
	}
	defer s.Close()

	ctx := context.Background()

	keyLock := "/lock" + key
	l := concurrency.NewMutex(s, keyLock)

	// acquire lock (or wait to have it)
	if err := l.Lock(ctx); err != nil {
		return "", err
	}

	// reseting
	if len(resetLock.IsOwner().Key) > 1 {
		if err := l.Unlock(ctx); err != nil {
			return "", err
		}
		return "", RefusedByLockError
	}

	fmt.Println("[Get] acquired lock for get ", keyLock)

	var getResp *clientv3.GetResponse
	// 实例化一个用于操作ETCD的KV
	kv := clientv3.NewKV(cli)

	if getResp, err = kv.Get(context.TODO(), key); err != nil {
		fmt.Println(err)
		return
	}

	// 输出本次的Revision
	if getResp.Kvs != nil {
		fmt.Printf("[Get] Key : %s, Value : %s \n", getResp.Kvs[0].Key, getResp.Kvs[0].Value)
	}

	if err := l.Unlock(ctx); err != nil {
		return "", err
	}

	fmt.Println("[Get] released lock for get ", keyLock)
	fmt.Println()

	return string(getResp.Kvs[0].Value), nil
}

func UpdateGlobalMap(key, value string) (err error) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// reseting
	if len(resetLock.IsOwner().Key) > 1 {
		return RefusedByLockError
	}

	// create a sessions to aqcuire a lock
	s, err := concurrency.NewSession(cli, concurrency.WithTTL(10))
	if err != nil {
		return err
	}
	defer s.Close()

	ctx := context.Background()

	keyLock := "/lock" + key
	l := concurrency.NewMutex(s, keyLock)

	// acquire lock (or wait to have it)
	if err := l.Lock(ctx); err != nil {
		return err
	}

	// reseting
	if len(resetLock.IsOwner().Key) > 1 {
		if err := l.Unlock(ctx); err != nil {
			return err
		}
		return RefusedByLockError
	}

	fmt.Println("[Update] acquired lock for update ", keyLock)

	var putResp *clientv3.PutResponse
	// 实例化一个用于操作ETCD的KV
	kv := clientv3.NewKV(cli)

	if putResp, err = kv.Put(context.TODO(), key, value, clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(putResp.Header.Revision)
	if putResp.PrevKv != nil {
		fmt.Printf("[Update] preValue: %s CreateRevision : %d  ModRevision: %d  Version: %d \n",
			putResp.PrevKv.Value, putResp.PrevKv.CreateRevision, putResp.PrevKv.ModRevision, putResp.PrevKv.Version)
	}
	fmt.Println("[Update] curValue: ", value)

	if err := l.Unlock(ctx); err != nil {
		return err
	}

	fmt.Println("[Update] released lock for update ", keyLock)
	fmt.Println()

	return nil
}

func DeleteGlobalMap(key string) (err error) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// reseting
	if len(resetLock.IsOwner().Key) > 1 {
		return RefusedByLockError
	}

	// create a sessions to aqcuire a lock
	s, err := concurrency.NewSession(cli, concurrency.WithTTL(10))
	if err != nil {
		return err
	}
	defer s.Close()

	ctx := context.Background()

	keyLock := "/lock" + key
	l := concurrency.NewMutex(s, keyLock)

	// acquire lock (or wait to have it)
	if err := l.Lock(ctx); err != nil {
		return err
	}

	// reseting
	if len(resetLock.IsOwner().Key) > 1 {
		if err := l.Unlock(ctx); err != nil {
			return err
		}
		return RefusedByLockError
	}

	fmt.Println("[Delete] acquired lock for delete ", keyLock)

	// 实例化一个用于操作ETCD的KV
	kv := clientv3.NewKV(cli)

	res, err := kv.Delete(context.TODO(), key)
	if err != nil {
		return err
	} else {
		fmt.Printf("[Delete] delete %d key\n", res.Deleted)
		for _, preKv := range res.PrevKvs {
			fmt.Printf("[Delete] del key: %s, value: %s\n", preKv.Key, preKv.Value)
		}
	}

	if err := l.Unlock(ctx); err != nil {
		return err
	}

	fmt.Println("[Delete] released lock for delete ", keyLock)
	fmt.Println()

	return nil
}

func ResetGlobalMap(prefixKeyLock, prefixKey string) (err error) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	sReset, err := concurrency.NewSession(cli, concurrency.WithTTL(10))
	if err != nil {
		log.Fatal(err)
	}
	defer sReset.Close()
	ctxReset := context.Background()
	resetLock := concurrency.NewMutex(sReset, PrefixLock)

	// acquire lock (or wait to have it)
	if err := resetLock.Lock(ctxReset); err != nil {
		return err
	}

	fmt.Println("[Reset] acquired lock for reset ", prefixKeyLock)

	kv := clientv3.NewKV(cli)
	res, err := kv.Delete(context.TODO(), prefixKey, clientv3.WithPrevKV(), clientv3.WithPrefix())
	if err != nil {
		return err
	} else {
		fmt.Printf("[Reset] delete %d keys\n", res.Deleted)
		for _, preKv := range res.PrevKvs {
			fmt.Printf("[Reset] del key: %s, value: %s\n", preKv.Key, preKv.Value)
		}
	}

	if err := resetLock.Unlock(ctxReset); err != nil {
		return err
	}

	fmt.Println("[Reset] released lock for reset ", prefixKeyLock)
	fmt.Println()

	return nil
}
