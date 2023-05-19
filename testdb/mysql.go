package testdb

import (
	"fmt"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net"
	"os"
	"strings"
	"time"
)

const (
	dbHost            = "localhost"
	dbPort            = 3306
	dbDefaultUser     = "root"
	dbDefaultPassword = "123456"
	dbDefaultCharset  = "utf8mb4"
)

// InitServerAndClients start a mysql server and init mysql connection clients according to configs
func InitServerAndClients(configs map[string]MysqlConf) (clients map[string]*gorm.DB) {
	if len(configs) == 0 {
		panic(errors.New("empty configs"))
	}

	initMySQLServerAndDatabases(configs)
	return initMySQLClients(configs)
}

type TableOption func(db *memory.Database)

type MysqlConf struct {
	Name            string
	Addr            string        `yaml:"addr"`
	Database        string        `yaml:"database"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Charset         string        `yaml:"charset"`
	MaxIdleConns    int           `yaml:"maxidleconns"`
	MaxOpenConns    int           `yaml:"maxopenconns"`
	ConnMaxIdlTime  time.Duration `yaml:"maxIdleTime"`
	ConnMaxLifeTime time.Duration `yaml:"connMaxLifeTime"`
	ConnTimeOut     time.Duration `yaml:"connTimeOut"`
	WriteTimeOut    time.Duration `yaml:"writeTimeOut"`
	ReadTimeOut     time.Duration `yaml:"readTimeOut"`
	Tables          []TableOption
}

func StartMySQLServer(databases ...sql.Database) {
	_ = sql.NewEmptyContext()
	engine := sqle.NewDefault(
		//memory.NewDBProvider(
		//	databases...,
		//),
		memory.NewMemoryDBProvider(
			databases...,
		),
	)

	//if err := createDefaultDBUser(engine); err != nil {
	//	panic(err)
	//}

	config := server.Config{
		Auth:     auth.NewNativeSingle("root", dbDefaultPassword, auth.AllPermissions),
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", dbHost, dbPort),
		Version:  "5.7.24-log",
	}
	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}

	if err = s.Start(); err != nil {
		panic(err)
	}
}

//func createDefaultDBUser(engine *sqle.Engine) error {
//	mysqlDb := engine.Analyzer.Catalog.MySQLDb
//	mysqlDb.AddSuperUser(dbDefaultUser, dbHost, dbDefaultPassword)
//	return nil
//}

func initMySQLClients(configs map[string]MysqlConf) (clients map[string]*gorm.DB) {
	var err error
	clients = make(map[string]*gorm.DB, len(configs))

	for name, dbConf := range configs {

		c := &gorm.Config{
			SkipDefaultTransaction:                   true,
			NamingStrategy:                           nil,
			FullSaveAssociations:                     false,
			NowFunc:                                  nil,
			DryRun:                                   false,
			PrepareStmt:                              false,
			DisableAutomaticPing:                     false,
			DisableForeignKeyConstraintWhenMigrating: false,
			AllowGlobalUpdate:                        false,
			ClauseBuilders:                           nil,
			ConnPool:                                 nil,
			Dialector:                                nil,
			Plugins:                                  nil,
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=True&loc=Asia%%2FShanghai",
			dbConf.User,
			dbConf.Password,
			dbConf.Addr,
			dbConf.Database,
			dbConf.ConnTimeOut,
			dbConf.ReadTimeOut,
			dbConf.WriteTimeOut,
		)

		clients[name], err = gorm.Open(mysql.Open(dsn), c)
		if err != nil {
			panic(err)
		}
	}

	return
}

func initMySQLServerAndDatabases(configs map[string]MysqlConf) {
	memoryDBs := make([]sql.Database, 0, len(configs))
	for _, v := range configs {
		memoryDBs = append(memoryDBs, registerDB(v.Database, v.Tables...))
	}

	go StartMySQLServer(memoryDBs...)

	conn := make(chan struct{}, 1)
	ready := time.Now()
	tick := time.NewTicker(time.Millisecond * 100)
	defer tick.Stop()

	for {
		isPortOn(dbHost, dbPort, conn)

		select {
		case now := <-tick.C:
			if sub := now.Sub(ready); sub > time.Second*120 {
				panic(fmt.Sprintf("failed to start test db server, cost%fs", sub.Seconds()))
			}
		case <-conn:
			return
		}
	}
}

func InitData(sqlFile string, db *gorm.DB) error {
	b, err := os.ReadFile(sqlFile)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, _sql := range strings.Split(strings.Trim(string(b), "\n"), ";") {
		if err = db.Exec(_sql).Error; err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func isPortOn(host string, port int, ch chan struct{}) {
	hostPort := fmt.Sprintf("%s:%d", host, port)
	// 检测端口
	conn, err := net.DialTimeout("tcp", hostPort, 3*time.Second)
	if err != nil || conn == nil {
		return
	}

	_ = conn.Close()
	ch <- struct{}{}
}

func registerDB(dbName string, opts ...TableOption) *memory.Database {
	db := memory.NewDatabase(dbName)
	db.EnablePrimaryKeyIndexes()

	for _, opt := range opts {
		opt(db)
	}

	return db
}
