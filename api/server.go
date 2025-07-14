package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/volskyi-dmytro/st-bank/db/sqlc"
	"github.com/volskyi-dmytro/st-bank/token"
	"github.com/volskyi-dmytro/st-bank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	var tokenMaker token.Maker
	var err error
	
	switch config.TokenType {
	case "jwt":
		tokenMaker, err = token.NewJWTMaker(config.TokenSymmetricKey)
	case "paseto":
		tokenMaker, err = token.NewPasetoMaker(config.TokenSymmetricKey)
	default:
		tokenMaker, err = token.NewPasetoMaker(config.TokenSymmetricKey)
	}
	
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	
	// Register custom validators
	RegisterValidators()
	
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PUT("/accounts/:id", server.updateAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)
	
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error" : err.Error()}
}