package main

import "errors"

var RefusedByLockError = errors.New("refused by lock")

// key
var (
	PrefixKey  = "/globalmap/"
	PrefixLock = "/lock/"
)
