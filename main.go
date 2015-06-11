package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	// import server packages
	serverMiddlewares "github.com/sogko/slumber/middlewares"
	"github.com/sogko/slumber/server"
	"github.com/sogko/slumber/sessions"
	"github.com/sogko/slumber/users"

	// import our project packages
	"github.com/sogko/slumber-books-api-example/books"
	"github.com/sogko/slumber-books-api-example/hooks"
	"time"
)

func main() {

	// try to load signing keys for token authority
	// NOTE: DO NOT USE THESE KEYS FOR PRODUCTION! FOR DEMO ONLY
	privateSigningKey, err := ioutil.ReadFile("keys/demo.rsa")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading private signing key: %v", err.Error())))
	}
	publicSigningKey, err := ioutil.ReadFile("keys/demo.rsa.pub")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading public signing key: %v", err.Error())))
	}

	// load and merge routes
	routes := users.Routes.Append(&sessions.Routes, &books.Routes)

	// load and merge ACL maps
	aclMap := users.ACL.Append(&sessions.ACL, &books.ACL)

	// set server configuration
	config := server.Config{
		Database: &serverMiddlewares.MongoDBOptions{
			ServerName:   "localhost",
			DatabaseName: "slumber-books-example",
		},
		Renderer: &serverMiddlewares.RendererOptions{
			IndentJSON: true,
		},
		TokenAuthority: &serverMiddlewares.TokenAuthorityOptions{
			PrivateSigningKey: privateSigningKey,
			PublicSigningKey:  publicSigningKey,
		},
		Routes:          &routes,
		ACLMap:          &aclMap,
		ControllerHooks: &hooks.HooksMap,
	}

	// init server and run
	s := server.NewServer(&config)

	// you can add your own middlewares here
	// and have it set up before routes are added to the router

	// setup route after middlewares
	s.SetupRoutes()

	// bam!
	gracefulTimeout := 5 * time.Second
	s.Run(":3001", gracefulTimeout)
}
