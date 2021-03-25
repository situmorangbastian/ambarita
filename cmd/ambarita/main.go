package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	log "github.com/sirupsen/logrus"
	"github.com/situmorangbastian/gower"
	"github.com/spf13/viper"

	articleHandler "github.com/situmorangbastian/ambarita/article/http"
	articleRepository "github.com/situmorangbastian/ambarita/article/repository"
	articleUsecase "github.com/situmorangbastian/ambarita/article/usecase"
)

func init() {
	viper.SetConfigFile(`configs/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Warn("service is running in DEBUG Mode")
		return
	}
	log.SetLevel(log.InfoLevel)
	log.Info("service is running in PRODUCTION Mode")
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
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

	// Server
	app := fiber.New(fiber.Config{
		Prefork:      viper.GetBool("gofiber.prefork"),
		ErrorHandler: gower.ErrMiddleware,
	})
	app.Use(recover.New())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Info("Gracefully shutting down...")
		if err := app.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

	// Domain
	ar := articleRepository.NewMysqlRepository(dbConn)
	au := articleUsecase.NewArticleUsecase(ar)
	articleHandler.NewHandler(app, au)

	// Start server
	if err := app.Listen(viper.GetString("server.address")); err != nil {
		log.Fatal(err)
	}
}
