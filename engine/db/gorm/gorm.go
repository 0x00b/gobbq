package gorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MysqlDB 数据库操作实例
type GormDB struct {
	db *gorm.DB
	*Config
}

type Config struct {
	// for log
	SlowThreshold time.Duration `json:"slow_threshold" yaml:"slow_threshold"`

	LogLevel string `json:"log_level" yaml:"log_level"`
	UserName string `json:"user_name" yaml:"user_name"`
	Password string `json:"password" yaml:"password"`
	Protocol string `json:"protocol" yaml:"protocol"`
	Host     string `json:"host" yaml:"host"`
	DBName   string `json:"db_name" yaml:"db_name"`

	MaxOpen int `json:"max_open" yaml:"max_open"`
	MaxIdle int `json:"max_idle" yaml:"max_idle"`
	MaxLife int `json:"max_life" yaml:"max_life"`

	Charset string `json:"charset" yaml:"charset"`

	Params map[string]string `json:"params" yaml:"params"`

	//DatabaseType default mysql
	DatabaseType string `json:"database_type" yaml:"database_type"`

	Tabler
}

func FormatGormDB(db *gorm.DB) *GormDB {
	return &GormDB{
		db: db,
	}
}

func NewGormDB(c Config) *GormDB {
	e := CheckMysqlDatabase(c)
	if e != nil {
		panic(e)
	}

	return NewOnlyReadGormDB(c)
}

func NewOnlyReadGormDB(c Config) *GormDB {
	if c.Charset == "" {
		c.Charset = "utf8mb4"
	}

	db, err := newMysqlGormDB(c.UserName, c.Password, c.Protocol, c.Host,
		c.DBName, c.MaxOpen, c.MaxIdle, c.MaxLife, c.Params, c.Charset)

	if err != nil {
		panic(err)
	}

	if c.Tabler == nil {
		c.Tabler = defaultTabler{}
	}

	if c.SlowThreshold <= 0 {
		c.SlowThreshold = 2 * time.Second
	}

	db.Logger = NewGormLog(LogConfig{
		SlowThreshold: (c.SlowThreshold),
		LogLevel:      ParseGormLevel(c.LogLevel),
	})

	gdb := &GormDB{
		Config: &c,
		db:     db,
	}

	return gdb
}

//CheckMysqlDatabase 检查数据库是否存在，不存在则创建, for mysql
func CheckMysqlDatabase(c Config) error {
	dbURL := c.UserName + ":" + c.Password + "@tcp(" + c.Host + ")/?charset=" + c.Charset
	conn, err := sql.Open("mysql", dbURL)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil || conn == nil {
		return fmt.Errorf("open mysql fail, %v", err)
	}
	_, err = conn.Exec("CREATE DATABASE IF NOT EXISTS " + c.DBName)
	if err != nil {
		return fmt.Errorf("exec sql fail,%v", err)
	}
	return nil
}

func (c *GormDB) C(ctx context.Context, forupdate ...int) *gorm.DB {
	db := c.db.WithContext(ctx)
	if len(forupdate) > 0 {
		// db = db.Set("gorm:query_option", "FOR UPDATE")
		if forupdate[0] == 0 {
			db = db.Clauses(clause.Locking{Strength: "SHARE"})
		} else {
			db = db.Clauses(clause.Locking{Strength: "UPDATE"})
		}
	}
	return db
}

func (c *GormDB) DB() *gorm.DB {
	return c.db
}

func (c *GormDB) Begin() *GormDB {
	return c.begin()
}

func (c *GormDB) begin() *GormDB {
	db := c.db.Begin()
	return &GormDB{
		db:     db,
		Config: c.Config,
	}
}

func (c *GormDB) Model(ctx context.Context, model interface{}, forupdate ...int) *gorm.DB {
	//tableName := c.DB.NewScope(model).TableName()
	db := c.db.WithContext(ctx).Scopes(c.Table(ctx, model)).Model(model)

	if len(forupdate) > 0 {
		// db = db.Set("gorm:query_option", "FOR UPDATE")
		if forupdate[0] == 0 {
			db = db.Clauses(clause.Locking{Strength: "SHARE"})
		} else {
			db = db.Clauses(clause.Locking{Strength: "UPDATE"})
		}
	}

	return db
}

// func InitReport(db *GormDB) {
// 	if db == nil {
// 		panic("InitReport: db is nil")
// 	}
// 	ago.InitGormReport(db.db)
// }

func newMysqlGormDB(
	username, password, protocol, host, dbname string,
	maxOpen, maxIdle, maxLife int, params map[string]string, charset string) (*gorm.DB, error) {
	// username:password@protocol(address)/dbname?param=value
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=%s&parseTime=True&loc=Local",
		username, password, protocol, host, dbname, charset)
	if len(params) > 0 {
		for key, value := range params {
			dsn += fmt.Sprintf("&%s=%s", key, value)
		}
	}
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// database, err := gorm.Open("mysql", dbName)
	if err != nil {
		return nil, err
	}

	underlyingDB, err := database.DB()
	if err != nil {
		return nil, err
	}
	if underlyingDB == nil {
		return nil, errors.New("underlying DB is invalid")
	}
	underlyingDB.SetMaxOpenConns(maxOpen)
	if maxIdle > 0 {
		underlyingDB.SetMaxIdleConns(maxIdle)
	}
	underlyingDB.SetConnMaxLifetime(time.Duration(int64(maxLife) * int64(time.Second)))

	return database, nil
}
