package main

import "github.com/Twister915/acme-sidecar/store"

func main() {
	kubeStorage := store.GetProvider("kubernetes")

}
