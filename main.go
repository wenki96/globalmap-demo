package main

import (
	"log"
)

func main() {

	// Create a etcd client
	// cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer cli.Close()

	// create a sessions to aqcuire a lock for reset

	// sReset, err := concurrency.NewSession(cli, concurrency.WithTTL(10))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer sReset.Close()
	// ctxReset := context.Background()
	// resetMutex := concurrency.NewMutex(sReset, prefixLock)

	//update
	err := UpdateGlobalMap(PrefixKey+"zecrey", "is amazing")
	if err != nil {
		log.Fatal(err)
	}

	//update
	err = UpdateGlobalMap(PrefixKey+"zecrey1", "is amazing")
	if err != nil {
		log.Fatal(err)
	}

	//get
	_, err = GetGlobalMap(PrefixKey + "zecrey")
	if err != nil {
		log.Fatal(err)
	}

	//delete
	err = DeleteGlobalMap(PrefixKey + "zecrey")
	if err != nil {
		log.Fatal(err)
	}

	// reset
	err = ResetGlobalMap(PrefixLock, PrefixKey)
	if err != nil {
		log.Fatal(err)
	}

}
