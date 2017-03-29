package main

import (
	"flag"
	"fmt"

	"fantasy/internal"
)

var WangyangRedis *internal.FRedis = new(internal.FRedis)

type Student struct {
	internal.Logger
}

func main() {
	var cfg_path = flag.String("cfg", "./config.cfg", "配置文件路径")
	flag.Parse()

	dcm := internal.DefaultConfigManager

	dcm.RegisterRedis(WangyangRedis)
	dcm.OpenConfig(*cfg_path)

	wy, _ := WangyangRedis.Get("wy")
	fmt.Println(wy)
}
