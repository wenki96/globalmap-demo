package main

import (
	"errors"

	"go.etcd.io/etcd/clientv3/concurrency"
)

var RefusedByLockError = errors.New("refused by lock")

// key
var (
	PrefixKey  = "/globalmap/"
	PrefixLock = "/lock/"
)

var resetLock = &concurrency.Mutex{}
