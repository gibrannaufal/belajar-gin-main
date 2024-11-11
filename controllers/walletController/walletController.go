package walletController

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gibrannaufal/belajar-main-gin/controllers"
	"github.com/gibrannaufal/belajar-main-gin/models"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	userID, err := controllers.StringToInt(c.Param("user_id"))
	if controllers.HandleError(c, "Invalid user_id", err, http.StatusBadRequest) {
		return
	}

	params := map[string]string{
		"username": c.Query("username"),
		"password": c.Query("password"),
	}

	if !controllers.CheckUserCredentials(c, userID, params["username"], params["password"]) {
		controllers.HandleError(c, "Invalid credentials", nil, http.StatusUnauthorized)
		return
	}

	var Wallet []models.Wallet
	if err := models.DB.Where("user_id = ?", userID).Find(&Wallet).Error; err != nil {
		controllers.HandleError(c, "Failed to retrieve wallets", err, http.StatusInternalServerError)
		return
	}

	controllers.HandleOK(c, "Wallets retrieved successfully", Wallet)
}

func Create(c *gin.Context) {

	userID, err := controllers.StringToInt(c.Query("user_id"))
	if controllers.HandleError(c, "Invalid user_id", err, http.StatusBadRequest) {
		return
	}

	balance, err := controllers.StringToInt(c.Query("balance"))
	if controllers.HandleError(c, "Invalid balance", err, http.StatusBadRequest) {
		return
	}

	paramsAuthorize := map[string]string{
		"username": c.Query("username"),
		"password": c.Query("password"),
	}

	if !controllers.CheckUserCredentials(c, userID, paramsAuthorize["username"], paramsAuthorize["password"]) {
		controllers.HandleError(c, "Invalid credentials", nil, http.StatusUnauthorized)
		return
	}

	paramsValue := map[string]interface{}{
		"user_id": userID,
		"balance": balance,
	}

	for key, value := range paramsValue {
		if value == "" {
			controllers.HandleError(c, fmt.Sprintf("%s is required", key), nil, http.StatusBadRequest)

			return
		}
	}

	var existingWallet models.Wallet
	if err := models.DB.Where("user_id = ?", paramsValue["user_id"]).First(&existingWallet).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet for this user already exists"})
		return
	}

	Wallet := models.Wallet{
		UserId:    paramsValue["user_id"].(int),
		Balance:   paramsValue["balance"].(int),
		IsDeleted: 0,
		CreatedAt: int(time.Now().Unix()),
	}

	if err := models.DB.Create(&Wallet).Error; err != nil {
		controllers.HandleError(c, "Failed to create Wallet", err, http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet created successfully", "Wallet": Wallet})
}

func Update(c *gin.Context) {

	userID, err := controllers.StringToInt(c.Query("user_id"))
	if controllers.HandleError(c, "Invalid user_id", err, http.StatusBadRequest) {
		return
	}

	paramsAuthorize := map[string]string{
		"username": c.Query("username"),
		"password": c.Query("password"),
	}

	if !controllers.CheckUserCredentials(c, userID, paramsAuthorize["username"], paramsAuthorize["password"]) {
		controllers.HandleError(c, "Invalid credentials", nil, http.StatusUnauthorized)
		return
	}

	walletId, err := controllers.StringToInt(c.Param("wallet_id"))
	if controllers.HandleError(c, "Invalid wallet", err, http.StatusBadRequest) {
		return
	}
	balance, err := controllers.StringToInt(c.Query("balance"))
	if controllers.HandleError(c, "Invalid balance", err, http.StatusBadRequest) {
		return
	}

	if walletId == 0 {
		controllers.HandleError(c, "Add the id wallet", nil, http.StatusBadRequest)

		return
	}

	paramsValue := map[string]interface{}{
		"user_id": userID,
		"balance": balance,
	}

	var existingWallet models.Wallet
	if err := models.DB.Where("id = ? AND user_id = ?", walletId, paramsValue["user_id"]).First(&existingWallet).Error; err != nil {
		controllers.HandleError(c, "Wallet ID and user ID do not match", err, http.StatusBadRequest)
		return
	}

	Wallet := models.Wallet{
		UserId:    paramsValue["user_id"].(int),
		Balance:   paramsValue["balance"].(int),
		IsDeleted: 0,
		UpdatedAt: int(time.Now().Unix()),
	}

	if err := models.DB.Model(&models.Wallet{}).Where("id = ?", walletId).Updates(Wallet).Error; err != nil {
		controllers.HandleError(c, "Failed to update Wallet", nil, http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet updated successfully"})
}

func Withdraw(c *gin.Context) {

	userID, err := controllers.StringToInt(c.Query("user_id"))
	if controllers.HandleError(c, "Invalid user_id", err, http.StatusBadRequest) {
		return
	}

	paramsAuthorize := map[string]string{
		"username": c.Query("username"),
		"password": c.Query("password"),
	}

	if !controllers.CheckUserCredentials(c, userID, paramsAuthorize["username"], paramsAuthorize["password"]) {
		controllers.HandleError(c, "Invalid credentials", nil, http.StatusUnauthorized)
		return
	}

	walletId, err := controllers.StringToInt(c.Query("wallet_id"))
	if controllers.HandleError(c, "Invalid Wallet Id", err, http.StatusBadRequest) {
		return
	}

	var Wallet models.Wallet
	if err := models.DB.Where("id = ? AND user_id = ?", walletId, userID).First(&Wallet).Error; err != nil {
		controllers.HandleError(c, "Wallet ID and user ID do not match", err, http.StatusBadRequest)
		return
	}

	amount, err := controllers.StringToInt(c.Query("amount"))
	if controllers.HandleError(c, "Invalid Wallet Id", err, http.StatusBadRequest) {
		return
	}
	if Wallet.Balance < amount {
		err := errors.New("insufficient balance")
		controllers.HandleError(c, "Insufficient balance", err, http.StatusBadRequest)
		return
	}

	transaction := models.Transaction{
		WalletId:        walletId,
		Amount:          amount,
		TransactionType: "withdraw",
		IsDeleted:       0,
		CreatedAt:       int(time.Now().Unix()),
	}

	if err := models.DB.Create(&transaction).Error; err != nil {
		controllers.HandleError(c, "Failed to create Transaction", err, http.StatusInternalServerError)

		return
	}

	Wallet.Balance -= amount
	if err := models.DB.Save(&Wallet).Error; err != nil {
		controllers.HandleError(c, "Failed to update wallet balance", err, http.StatusInternalServerError)

		return
	}

	controllers.HandleOK(c, "Wallet withdraw successfully", Wallet)

}

func Deposit(c *gin.Context) {

	userID, err := controllers.StringToInt(c.Query("user_id"))
	if controllers.HandleError(c, "Invalid user_id", err, http.StatusBadRequest) {
		return
	}

	paramsAuthorize := map[string]string{
		"username": c.Query("username"),
		"password": c.Query("password"),
	}

	if !controllers.CheckUserCredentials(c, userID, paramsAuthorize["username"], paramsAuthorize["password"]) {
		controllers.HandleError(c, "Invalid credentials", nil, http.StatusUnauthorized)
		return
	}

	walletId, err := controllers.StringToInt(c.Query("wallet_id"))
	if controllers.HandleError(c, "Invalid Wallet Id", err, http.StatusBadRequest) {
		return
	}

	var Wallet models.Wallet
	if err := models.DB.Where("id = ? AND user_id = ?", walletId, userID).First(&Wallet).Error; err != nil {
		controllers.HandleError(c, "Wallet ID and user ID do not match", err, http.StatusBadRequest)
		return
	}

	amount, err := controllers.StringToInt(c.Query("amount"))
	if controllers.HandleError(c, "Invalid amount", err, http.StatusBadRequest) {
		return
	}

	transaction := models.Transaction{
		WalletId:        walletId,
		Amount:          amount,
		TransactionType: "deposit",
		IsDeleted:       0,
		CreatedAt:       int(time.Now().Unix()),
	}

	if err := models.DB.Create(&transaction).Error; err != nil {
		controllers.HandleError(c, "Failed to create Transaction", err, http.StatusInternalServerError)
		return
	}

	Wallet.Balance += amount
	if err := models.DB.Save(&Wallet).Error; err != nil {
		controllers.HandleError(c, "Failed to update wallet balance", err, http.StatusInternalServerError)
		return
	}

	controllers.HandleOK(c, "Wallet deposit successful", Wallet)
}
