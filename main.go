package main

import (
	"context"
	"log"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

func main() {
	// key
	prefixKey := "/globalmap/"
	prefixLock := "/lock/"

	// Create a etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// create a sessions to aqcuire a lock for reset
	sReset, err := concurrency.NewSession(cli, concurrency.WithTTL(10))
	if err != nil {
		log.Fatal(err)
	}
	defer sReset.Close()
	ctxReset := context.Background()
	resetMutex := concurrency.NewMutex(sReset, prefixLock)

	//update
	err = UpdateGlobalMap(cli, prefixKey+"zecrey", "is amazing", resetMutex)
	if err != nil {
		log.Fatal(err)
	}

	//update
	err = UpdateGlobalMap(cli, prefixKey+"zecrey1", "is amazing", resetMutex)
	if err != nil {
		log.Fatal(err)
	}

	//get
	_, err = GetGlobalMap(cli, prefixKey+"zecrey", resetMutex)
	if err != nil {
		log.Fatal(err)
	}

	// //delete
	// err = DeleteGlobalMap(cli, prefixKey+"zecrey", resetMutex)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// reset
	err = ResetGlobalMap(cli, prefixLock, prefixKey, ctxReset, resetMutex)
	if err != nil {
		log.Fatal(err)
	}

}
