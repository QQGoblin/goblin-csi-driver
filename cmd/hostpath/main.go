package main

import (
	"flag"
	"goblin-csi-driver/internal"
	"goblin-csi-driver/internal/hostpath"
	utilconfig "goblin-csi-driver/internal/utils/config"
	"k8s.io/klog/v2"
)

var (
	config     string
	nodeid     string
	drivername string
	endpoint   string
	ptype      string
)

func init() {

	// common flags
	flag.StringVar(&config, "config", "csi-driver.json", "csi driver config file.")
	flag.StringVar(&nodeid, "nodeid", "", "node id for csi plugin.")
	flag.StringVar(&drivername, "drivername", "hostpath.csi.goblin.io", "csi driver name for this csi.")
	flag.StringVar(&endpoint, "endpoint", internal.DefaultEndPoint, "default endpoint.")
	klog.InitFlags(nil)

	if err := flag.Set("logtostderr", "true"); err != nil {
		klog.Exitf("failed to set logtostderr flag: %v", err)
	}

	flag.Parse()
}

func main() {

	internal.DriverVersion = "canary"
	option := utilconfig.NewOptions()
	driver := hostpath.NewDriver(drivername, nodeid, endpoint, option)
	driver.Run()
}
