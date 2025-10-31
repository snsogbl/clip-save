package main

import "sync/atomic"

var forceQuitFlag atomic.Bool

func setForceQuit() {
	forceQuitFlag.Store(true)
}

func isForceQuit() bool {
	return forceQuitFlag.Load()
}

func clearForceQuit() {
	forceQuitFlag.Store(false)
}
