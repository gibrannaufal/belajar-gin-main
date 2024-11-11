package productController

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
	var Product []models.Product
	models.DB.Find(&Product)
	controllers.HandleOK(c, "Product Found", Product)

}
func Show(c *gin.Context) {
	id := c.Param("id")

	var Product models.Product
	result := models.DB.First(&Product, id)
	if result.Error != nil {
		controllers.HandleError(c, "Product not found", result.Error, http.StatusNotFound)

		return
	}

	controllers.HandleOK(c, "Product Found", Product)
}
func Create(c *gin.Context) {

	params := map[string]string{
		"name":         c.Query("name"),
		"description":  c.Query("description"),
		"price":        c.Query("price"),
		"availability": c.Query("availability"),
	}

	parsedParams := make(map[string]int)
	for key, value := range params {
		if value == "" {
			err := errors.New("Required")
			controllers.HandleError(c, fmt.Sprintf("%s is required", key), err, http.StatusBadRequest)
			return
		}

		if key == "price" || key == "availability" {
			intValue, err := controllers.StringToInt(value)
			if err != nil {
				controllers.HandleError(c, fmt.Sprintf("%s must be an integer", key), err, http.StatusBadRequest)
				return
			}
			parsedParams[key] = intValue
		}
	}

	product := models.Product{
		Name:         params["name"],
		Description:  params["description"],
		Price:        parsedParams["price"],
		Availability: parsedParams["availability"],
		IsDeleted:    0,
		CreatedAt:    int(time.Now().Unix()),
	}

	if err := models.DB.Create(&product).Error; err != nil {
		controllers.HandleError(c, "Product not found", err, http.StatusInternalServerError)

		return
	}

	controllers.HandleOK(c, "Product created successfully", product)

}
func Update(c *gin.Context) {
	productid := c.Param("id")

	params := map[string]string{
		"name":         c.Query("name"),
		"description":  c.Query("description"),
		"price":        c.Query("price"),
		"availability": c.Query("availability"),
	}

	parsedParams := make(map[string]int)
	for key, value := range params {
		if key == "price" || key == "availability" {
			intValue, err := controllers.StringToInt(value)
			if err != nil {
				controllers.HandleError(c, fmt.Sprintf("%s must be an integer", key), err, http.StatusBadRequest)
				return
			}
			parsedParams[key] = intValue
		}
	}

	var product models.Product
	if err := models.DB.Where("id = ?", productid).First(&product).Error; err != nil {
		controllers.HandleError(c, "Product not found", err, http.StatusNotFound)
		return
	}

	productData := models.Product{
		Name: func() string {
			if params["name"] == "" {
				return product.Name
			}
			return params["name"]
		}(),
		Description: func() string {
			if params["description"] == "" {
				return product.Description
			}
			return params["description"]
		}(),
		Price: func() int {
			if parsedParams["price"] == 0 {
				return product.Price
			}
			return parsedParams["price"]
		}(),
		Availability: func() int {
			if parsedParams["availability"] == 0 {
				return product.Availability
			}
			return parsedParams["availability"]
		}(),
		IsDeleted: 0,
		UpdatedAt: int(time.Now().Unix()),
	}

	if productid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Add the Id"})
		return
	}

	if err := models.DB.Model(&models.Product{}).Where("id = ?", productid).Updates(productData).Error; err != nil {
		controllers.HandleError(c, "Failed to update product", err, http.StatusInternalServerError)

		return
	}

	controllers.HandleOK(c, "Product updated successfully", productid)

}

func Delete(c *gin.Context) {
	id := c.Param("id")

	isDeleted := 1

	if err := models.DB.Model(&models.Product{}).Where("id = ?", id).Updates(map[string]interface{}{"is_deleted": isDeleted}).Error; err != nil {
		controllers.HandleError(c, "Failed to update product", err, http.StatusInternalServerError)
		return
	}

	controllers.HandleOK(c, "Product deleted successfully", id)
}
