// golden-field 组合根（VP-023 产线化实验下游 · 双方言参数化：R4）。
// 消费形态：apps/api/kernel + assembly（B+ 层）+ 模块包——Go 类型推断，无 internal 命名。
// 用法：golden-field [-dialect sqlite|postgres] [-dsn <conn>|<sqlite-path>]
//   SQLite 默认（内嵌默认）：golden-field ./data.db
//   PostgreSQL（生产权威方言）：golden-field -dialect postgres -dsn "postgres://…"
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/assembly"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/compiled"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/modules/users"
)

func main() {
	// v0.3.0 迁移（按 changelog-breaking-v0.3.0 第 2 步改写）：JoinKeys → JoinIdentifiers
	_ = kernel.JoinIdentifiers("a", "b")
	dialect := flag.String("dialect", "sqlite", "sqlite | postgres")
	dsn := flag.String("dsn", "", "连接串（sqlite = 文件路径；postgres = DSN）")
	flag.Parse()
	ctx := context.Background()

	catalog, err := compiled.PersistenceCatalog()
	if err != nil {
		log.Fatal("catalog:", err)
	}
	if *dsn == "" {
		log.Fatal("usage: golden-field [-dialect sqlite|postgres] -dsn <path|conn>")
	}
	var path, pgDSN string
	if *dialect == "postgres" {
		pgDSN = *dsn
	} else {
		path = *dsn
	}
	st, err := assembly.OpenStore(ctx, kernel.Dialect(*dialect), path, pgDSN, catalog)
	if err != nil {
		log.Fatal("store:", err)
	}
	defer st.Close()

	authn := assembly.NewAuthenticator([]byte("golden-field-secret"), 15*time.Minute, 24*time.Hour, st)
	repo := authsession.NewRepository(st)
	ops := operationlog.NewRepository(st)
	mailer := assembly.NewMailSender(st, 1000)
	usersProvider := users.New(authn, repo, ops, mailer, "http://127.0.0.1")

	plan, err := kernel.ResolveProfile("admin", nil)
	if err != nil {
		log.Fatal("resolve:", err)
	}
	builtin := kernel.BuiltinModules()
	if _, err := kernel.NewRegistry(builtin); err != nil {
		log.Fatal("registry:", err)
	}
	byID := make(map[string]kernel.Module, len(builtin))
	for _, m := range builtin {
		byID[m.ID] = m
	}
	planModules := make([]kernel.Module, 0, len(plan.Modules))
	for _, id := range plan.Modules {
		if m, ok := byID[id]; ok {
			planModules = append(planModules, m)
		}
	}
	contribs, err := kernel.RegisterContributions(ctx, kernel.Plan{
		Modules:      planModules,
		Capabilities: kernel.StandardAdminCapabilities(),
	}, []kernel.Provider{usersProvider})
	if err != nil {
		log.Fatal("contribs:", err)
	}

	fmt.Printf("golden-field kernel=%s profile=%s dialect=%s fresh=%v contrib{r=%d p=%d perm=%d nav=%d frag=%d}\n",
		kernel.KernelAPIVersion, plan.Name, st.Dialect(), st.WasFresh(),
		len(contribs.Routes), len(contribs.Pages), len(contribs.Permissions),
		len(contribs.Navigation), len(contribs.Fragments))
}