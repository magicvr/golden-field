// golden-field 组合根（VP-024 R2 升级 · 2026-08-29）：thin wrapper → 公开 serve 面。
// 消费形态：apps/api v0.4.0（server pkg = config 装载 + 标准组合 + 中央面 + RT-D02 停机）。
// 用法：golden-field [-config <path>] [-dialect sqlite|postgres] [-dsn <path|conn>] [-addr <addr>]
package main

import (
	"flag"
	"log"

	"github.com/magicvr/schema-ui-core/apps/api/server"
)

func main() {
	configPath := flag.String("config", "", "serve 配置路径（缺省 = 内嵌默认 + env；参考 config.example.yaml）")
	dialect := flag.String("dialect", "", "覆盖方言（sqlite|postgres）")
	dsn := flag.String("dsn", "", "覆盖连接串（sqlite=文件路径；postgres=DSN）")
	addr := flag.String("addr", "", "覆盖监听地址")
	flag.Parse()

	cfg, err := server.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("config:", err)
	}
	if *dialect != "" {
		cfg.DBDialect = *dialect
	}
	if *dsn != "" {
		if cfg.DBDialect == "postgres" {
			cfg.DBDSN = *dsn
		} else {
			cfg.DBPath = *dsn
		}
	}
	if *addr != "" {
		cfg.HTTPAddr = *addr
	}
	if err := server.Serve(server.Options{Config: cfg}); err != nil {
		log.Fatal("serve:", err)
	}
}