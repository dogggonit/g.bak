package server

import (
	"flag"
	"fmt"
	"g.bak/internal/pkg/application"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func Server() {
	fs := flag.NewFlagSet("server", flag.ExitOnError)
	dbName := fs.String("db", "email.db", "filename for database")
	host := fs.String("h", "localhost", "host")
	port := fs.Uint("p", 8080, "port")
	_ = fs.Parse(os.Args)

	db, err := application.NewSqlite3(*dbName)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := gin.Default()
	r.Use(cors.Default())
	group := r.Group("api").Group("v1")

	group.GET("labels", application.GetLabels(db)...)
	group.GET("token", application.GetTokenLink(db)...)
	group.POST("token", application.PostToken(db)...)

	email := group.Group("emails")

	email.GET("/:email/label/:label", application.GetListEmails(db)...)
	email.GET("/:email/message/:id", application.GetEmail(db)...)

	log.Println("starting http server")
	log.Fatal(r.Run(fmt.Sprintf("%s:%d", *host, *port)).Error())
}