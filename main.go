package main

import (
	"log"
)

func main() {

	gracefulShutdown()

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
