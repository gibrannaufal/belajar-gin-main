package transactionController

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gibrannaufal/belajar-main-gin/controllers"
	"github.com/gibrannaufal/belajar-main-gin/models"
	"github.com/gin-gonic/gin"
)

func BuyProduct(c *gin.Context) {

	params := map[string]string{
		"user_id":        c.Query("user_id"),
		"wallet_id":      c.Query("wallet_id"),
		"product_id":     c.Query("product_id"),
		"amount_product": c.Query("amount_product"),
	}

	parsedParams := make(map[string]int)
	for key, value := range params {
		if value == "" {
			err := errors.New("Required")

			controllers.HandleError(c, fmt.Sprintf("%s is required", key), err, http.StatusBadRequest)

			return
		}
		intValue, err := controllers.StringToInt(value)
		if err != nil {
			controllers.HandleError(c, fmt.Sprintf("%s must be an integer", key), err, http.StatusBadRequest)

			return
		}

		parsedParams[key] = intValue
	}

	paramsAuthorize := map[string]string{
		"username": c.Query("username"),
		"password": c.Query("password"),
	}

	if !controllers.CheckUserCredentials(c, parsedParams["user_id"], paramsAuthorize["username"], paramsAuthorize["password"]) {
		controllers.HandleError(c, "Invalid credentials", nil, http.StatusUnauthorized)
		return
	}

	var wallet models.Wallet
	if err := models.DB.Where("id = ? AND user_id = ?", parsedParams["wallet_id"], parsedParams["user_id"]).First(&wallet).Error; err != nil {
		controllers.HandleError(c, "Wallet ID and user ID do not match", err, http.StatusBadRequest)
		return
	}

	// Retrieve product
	var product models.Product
	if err := models.DB.Where("id = ?", parsedParams["product_id"]).First(&product).Error; err != nil {
		controllers.HandleError(c, "Product not found", err, http.StatusNotFound)

		return
	}

	if product.IsDeleted == 1 {
		err := errors.New("product was deleted")
		controllers.HandleError(c, "product was deleted", err, http.StatusNotFound)

		return
	}

	if product.Availability < parsedParams["amount_product"] || product.Availability == 0 {
		err := errors.New("Required")
		controllers.HandleError(c, "insufficient product stoc", err, http.StatusNotFound)

		return
	}

	totalPayment := product.Price * parsedParams["amount_product"]

	if wallet.Balance < totalPayment {
		err := errors.New("insufficient wallet balance")
		controllers.HandleError(c, "insufficient wallet balance", err, http.StatusBadRequest)

		return
	}

	// Update product availability
	product.UpdatedAt = int(time.Now().Unix())
	product.Availability -= parsedParams["amount_product"]
	if err := models.DB.Model(&models.Product{}).
		Where("id = ?", parsedParams["product_id"]).
		Updates(map[string]interface{}{
			"availability": product.Availability,
			"updated_at":   wallet.UpdatedAt,
		}).Error; err != nil {
		controllers.HandleError(c, "Failed to update product availability", err, http.StatusInternalServerError)
		return
	}

	transaction := models.Transaction{
		WalletId:        parsedParams["wallet_id"],
		Amount:          totalPayment,
		TransactionType: "withdraw",
		CreatedAt:       int(time.Now().Unix()),
		IsDeleted:       0,
	}
	if err := models.DB.Create(&transaction).Error; err != nil {
		controllers.HandleError(c, "Failed to create transaction", err, http.StatusInternalServerError)

		return
	}

	// Update wallet balance
	wallet.UpdatedAt = int(time.Now().Unix())
	wallet.Balance -= totalPayment

	if err := models.DB.Model(&wallet).Where("id = ? AND user_id = ?", parsedParams["wallet_id"], parsedParams["user_id"]).Update("balance", wallet.Balance).Update("updated_at", wallet.UpdatedAt).Error; err != nil {
		controllers.HandleError(c, "Failed to update wallet balance", err, http.StatusInternalServerError)
		return
	}

	controllers.HandleOK(c, "Product sold successfully", map[string]interface{}{
		"wallet":  wallet,
		"product": product,
	})

}

func SellProduct(c *gin.Context) {

	params := map[string]string{
		"user_id":        c.Query("user_id"),
		"wallet_id":      c.Query("wallet_id"),
		"product_id":     c.Query("product_id"),
		"amount_product": c.Query("amount_product"),
	}

	parsedParams := make(map[string]int)
	for key, value := range params {
		if value == "" {
			err := errors.New("Required")
			controllers.HandleError(c, fmt.Sprintf("%s is required", key), err, http.StatusBadRequest)
			return
		}
		intValue, err := controllers.StringToInt(value)
		if err != nil {
			controllers.HandleError(c, fmt.Sprintf("%s must be an integer", key), err, http.StatusBadRequest)
			return
		}
		parsedParams[key] = intValue
	}

	paramsAuthorize := map[string]string{
		"username": c.Query("username"),
		"password": c.Query("password"),
	}

	if !controllers.CheckUserCredentials(c, parsedParams["user_id"], paramsAuthorize["username"], paramsAuthorize["password"]) {
		controllers.HandleError(c, "Invalid credentials", nil, http.StatusUnauthorized)
		return
	}

	var wallet models.Wallet
	if err := models.DB.Where("id = ? AND user_id = ?", parsedParams["wallet_id"], parsedParams["user_id"]).First(&wallet).Error; err != nil {
		controllers.HandleError(c, "Wallet ID and user ID do not match", err, http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := models.DB.Where("id = ?", parsedParams["product_id"]).First(&product).Error; err != nil {
		controllers.HandleError(c, "Product not found", err, http.StatusNotFound)
		return
	}

	if product.IsDeleted == 1 {
		err := errors.New("product was deleted")
		controllers.HandleError(c, "product was deleted", err, http.StatusNotFound)

		return
	}

	if product.Availability < parsedParams["amount_product"] || product.Availability == 0 {
		err := errors.New("Required")
		controllers.HandleError(c, "insufficient product stock", err, http.StatusBadRequest)
		return
	}

	totalPayment := product.Price * parsedParams["amount_product"]

	product.UpdatedAt = int(time.Now().Unix())
	product.Availability += parsedParams["amount_product"]
	if err := models.DB.Model(&models.Product{}).
		Where("id = ?", parsedParams["product_id"]).
		Updates(map[string]interface{}{
			"availability": product.Availability,
			"updated_at":   wallet.UpdatedAt,
		}).Error; err != nil {
		controllers.HandleError(c, "Failed to update product availability", err, http.StatusInternalServerError)
		return
	}

	transaction := models.Transaction{
		WalletId:        parsedParams["wallet_id"],
		Amount:          totalPayment,
		TransactionType: "deposit",
		CreatedAt:       int(time.Now().Unix()),
		IsDeleted:       0,
	}

	if err := models.DB.Create(&transaction).Error; err != nil {
		controllers.HandleError(c, "Failed to create transaction", err, http.StatusInternalServerError)
		return
	}

	wallet.UpdatedAt = int(time.Now().Unix())
	wallet.Balance += totalPayment
	if err := models.DB.Model(&wallet).Where("id = ? AND user_id = ?", parsedParams["wallet_id"], parsedParams["user_id"]).Update("balance", wallet.Balance).Update("updated_at", wallet.UpdatedAt).Error; err != nil {
		controllers.HandleError(c, "Failed to update wallet balance", err, http.StatusInternalServerError)
		return
	}

	controllers.HandleOK(c, "Product sold successfully", map[string]interface{}{
		"wallet":  wallet,
		"product": product,
	})
}
