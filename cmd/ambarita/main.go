package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	articleHandler "github.com/situmorangbastian/ambarita/article/http"
	articleRepository "github.com/situmorangbastian/ambarita/article/repository"
	articleUsecase "github.com/situmorangbastian/ambarita/article/usecase"
	"github.com/situmorangbastian/ambarita/cmd/logger"
	"github.com/situmorangbastian/ambarita/models"
	"github.com/situmorangbastian/eclipse"
)

func init() {
	viper.SetConfigFile(`configs/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	logger.SetupLogs()
}

func main() {
	database := viper.GetString("database")

	var ar models.ArticleRepository

	switch database {
	case "mysql":
		dbHost := viper.GetString(`mysql.host`)
		dbPort := viper.GetString(`mysql.port`)
		dbUser := viper.GetString(`mysql.user`)
		dbPass := viper.GetString(`mysql.pass`)
		dbName := viper.GetString(`mysql.name`)
		connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
		val := url.Values{}
		val.Add("parseTime", "1")
		val.Add("loc", "Asia/Jakarta")
		dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
		dbConn, err := sql.Open(`mysql`, dsn)

		if err != nil {
			log.Fatal(err)
		}
		err = dbConn.Ping()
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			err := dbConn.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		// Migration
		driver, err := mysql.WithInstance(dbConn, &mysql.Config{})
		if err != nil {
			log.Fatal(err)
		}
		m, err := migrate.NewWithDatabaseInstance("file://migrations/", "mysql", driver)
		if err != nil {
			log.Fatal(err)
		}
		_ = m.Up()

		ar = articleRepository.NewMysqlRepository(dbConn)
	case "mongo":
		mongoURI := viper.GetString("mongo.uri")

		mongoClient, err := mongo.Connect(
			context.Background(),
			options.Client().ApplyURI(mongoURI).
				SetConnectTimeout(2*time.Second).
				SetServerSelectionTimeout(3*time.Second),
		)
		if err != nil {
			log.Fatal("Mongo connection failed: ", err.Error())
		}

		mongoDBName := viper.GetString("mongo.dbname")

		ar = articleRepository.NewMongoRepository(mongoClient.Database(mongoDBName))
	default:
		panic("please select your database: mysql or mongo")
	}

	e := echo.New()
	e.Use(eclipse.Error())

	// Domain
	au := articleUsecase.NewArticleUsecase(ar)
	articleHandler.NewHandler(e, au)

	// Start server
	go func() {
		if err := e.Start(viper.GetString("server.address")); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
