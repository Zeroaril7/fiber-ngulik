package databases

import (
	"fiber-ngulik/modules/user/models"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	DatabaseConnection struct {
		Connection *gorm.DB
	}

	Database struct {
		Name string
	}

	DBInterface interface {
		Connect(string) *DatabaseConnection
	}
)

var (
	DBConnect  *DatabaseConnection
	accessOnce sync.Once
	access     DBInterface
)

func (db *Database) Connect(dbname string) *DatabaseConnection {
	dbConnection := &DatabaseConnection{}
	master := db.Name

	if master != "" {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second * 5, // Slow SQL threshold
				LogLevel:                  logger.Info,     // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,            // Don't include params in the SQL log
				Colorful:                  true,            // Disable color
			},
		)

		db, err := gorm.Open(postgres.Open(master), &gorm.Config{Logger: newLogger})
		if err != nil {
			log.Fatal("postgres ", "can not connect PostgreSQL", "connect", err)
		}

		err = db.AutoMigrate(&models.User{})
		if err != nil {
			log.Fatal("Failed to run migrations:", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal("postgres ", "can not connect PostgreSQL", "connect", err)
		}

		sqlDB.SetConnMaxIdleTime(5)
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetConnMaxLifetime(time.Minute * 5)

		dbConnection.Connection = db
	}

	DBConnect = &DatabaseConnection{Connection: dbConnection.Connection}

	return DBConnect
}

func InitConnection(dns string, dbname string) DBInterface {
	if access != nil {
		return access
	}

	accessOnce.Do(func() {
		dbClient := NewDatabaseGorm(dns)
		dbClient.Connect(dbname)
		access = dbClient
	})

	return access
}

func NewDatabaseGorm(config interface{}) *Database {
	cfg := config.(string)

	return &Database{
		Name: cfg,
	}
}
