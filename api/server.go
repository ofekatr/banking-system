package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/ofekatr/simple-bank/db/sqlc"
	"github.com/pkg/errors"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	router.POST("/accounts", server.createAccount)

	router.GET("/accounts/:id", server.getAccount)

	router.GET("/accounts", server.listAccounts)

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	if err := server.router.Run(address); err != nil {
		return errors.Wrap(err, "cannot start server")
	}

	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
