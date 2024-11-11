package main

import (
	"github.com/gibrannaufal/belajar-main-gin/controllers/authController"
	"github.com/gibrannaufal/belajar-main-gin/controllers/productController"
	"github.com/gibrannaufal/belajar-main-gin/controllers/transactionController"
	"github.com/gibrannaufal/belajar-main-gin/controllers/walletController"
	"github.com/gibrannaufal/belajar-main-gin/middleware"
	"github.com/gibrannaufal/belajar-main-gin/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	models.ConnectDatabase()

	// Middleware CORS
	r.Use(cors.Default())

	r.GET("/api/products", middleware.AuthRequired(), productController.Index)
	r.GET("/api/products/:id", middleware.AuthRequired(), productController.Show)
	r.POST("/api/products", middleware.AuthRequired(), productController.Create)
	r.PUT("/api/products/:id", middleware.AuthRequired(), productController.Update)
	r.DELETE("/api/products/:id", middleware.AuthRequired(), productController.Delete)

	// wallet
	r.GET("/api/wallets/:user_id", middleware.AuthRequired(), walletController.Index)
	r.POST("/api/wallets", middleware.AuthRequired(), walletController.Create)
	r.PUT("/api/wallets/:wallet_id", middleware.AuthRequired(), walletController.Update)

	r.POST("/api/wallets/withdraw", middleware.AuthRequired(), walletController.Withdraw)
	r.POST("/api/wallets/deposit", middleware.AuthRequired(), walletController.Deposit)

	// transaction
	r.POST("/api/buy-product", middleware.AuthRequired(), transactionController.BuyProduct)
	r.POST("/api/sell-product", middleware.AuthRequired(), transactionController.SellProduct)

	r.POST("/api/login", authController.Login)
	r.POST("/api/logout", authController.Logout)

	r.Run()
}
