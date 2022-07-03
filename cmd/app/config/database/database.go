package database

import (
	"fmt"
	"github.com/xrexonx/togo/cmd/app/config/environment"
	"github.com/xrexonx/togo/internal/todo"
	"github.com/xrexonx/togo/internal/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

const (
	_dbConnectionSuccess = "Database connection was successful!"
	_dbConnectionFail    = "Could not connect to database"
	_dbDNS               = "%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	_dbErrCreateTable    = "Could not create tables to database"
)

// Config to maintain DB configuration properties
type Config struct {
	ServerName string
	User       string
	Password   string
	DB         string
}

// Init initialise database connection from DB configs
func Init() *gorm.DB {
	dbEnv := environment.GetDBEnv()
	dbConfig := Config{
		ServerName: dbEnv.DBHost + ":" + dbEnv.DBPort,
		User:       dbEnv.DBUser,
		Password:   dbEnv.DBPass,
		DB:         dbEnv.DBName,
	}
	log.Println("DBConfig:", dbConfig)
	dbConn, err := connect(dbConfig)
	if err != nil {
		log.Fatal(_dbConnectionFail)
	}

	// Connection pooling
	//sqlDB, _ := dbConn.DB()
	//sqlDB.SetConnMaxLifetime(time.Minute * 30)
	//sqlDB.SetMaxIdleConns(10)
	//sqlDB.SetMaxOpenConns(10)
	//defer sqlDB.Close()

	// Create tables
	dbMigrate(dbConn)

	return dbConn
}

// connect connection to the database
func connect(dbConfig Config) (*gorm.DB, error) {
	var err error
	dsn := getConnectionString(dbConfig)
	DBConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf(err.Error())
	}
	log.Println(_dbConnectionSuccess)
	return DBConn, nil
}

// getConnectionString setup database connection string
func getConnectionString(config Config) string {
	connectionString := fmt.Sprintf(_dbDNS, config.User, config.Password, config.ServerName, config.DB)
	return connectionString
}

// dbMigrate Create database and tables then seed sample users
func dbMigrate(db *gorm.DB) {
	errCreate := db.AutoMigrate(&user.User{}, &todo.Todo{})
	if errCreate != nil {
		log.Fatal(_dbErrCreateTable, errCreate)
	}

	// Seeds sample users
	result := db.Find(&user.User{})
	if result.RowsAffected == 0 {
		u1 := user.User{Name: "Rex", MaxDailyLimit: 5, Email: "rex@gmail.com.ph"}
		u2 := user.User{Name: "Riz", MaxDailyLimit: 4, Email: "roux@gmail.com.ph"}
		u3 := user.User{Name: "Roux", MaxDailyLimit: 3, Email: "roux@gmail.com.ph"}
		sampleUsers := []user.User{u1, u2, u3}
		for _, u := range sampleUsers {
			db.Create(&u)
		}
	}

}
