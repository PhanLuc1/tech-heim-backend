package controller

import (
	"fmt"
	"strings"

	"github.com/PhanLuc1/tech-heim-backend/database"
	"github.com/PhanLuc1/tech-heim-backend/models"
	"github.com/gin-gonic/gin"
)

func GetProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var products []models.Product
		categoryID := ctx.Query("categoryId")
		productTypeId := ctx.Query("productTypeId")
		sort := ctx.Query("sort")
		var params []interface{}
		query := "SELECT p.* FROM product AS p JOIN product_producttype AS ppt ON p.id = ppt.idproduct JOIN producttype AS pt ON ppt.idtype = pt.id JOIN productgroup AS pg ON pt.idGroup = pg.id WHERE "
		if categoryID != "" {
			query += fmt.Sprintf("idCategory = %s ", categoryID)
		}
		if productTypeId != "" {
			typeIds := strings.Split(productTypeId, "-")
			if categoryID != "" {
				query += "AND pt.id IN (?"
			} else {
				query += "pt.id IN (?"
			}
			for i := 1; i < len(typeIds); i++ {
				query += ", ?"
			}
			query += ") "
			for _, typeId := range typeIds {
				params = append(params, typeId)
			}
		}
		if sort == "ascending" {
			query += "ORDER BY p.currentPrice ASC"
		} else if sort == "descending" {
			query += "ORDER BY p.currentPrice DESC"
		}

		if categoryID == "" && productTypeId == "" {
			query = "SELECT * FROM product"
		}
		result, err := database.Client.Query(query, params...)
		if err != nil {
			ctx.JSON(200, gin.H{"message": "Không tìm thấy sản phẩm phù hợp"})
			return
		}
		for result.Next() {
			var product models.Product
			var imageProduct []models.Image
			var idCategory int
			err := result.Scan(&product.ProductId, &product.ProductName, &product.Rate, &product.Sold, &product.CurrentPrice, &product.LastPrice, &idCategory)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err})
				return
			}
			result1, err := database.Client.Query("SELECT productimage.url,productimage.description FROM productimage WHERE idProduct = ?", *product.ProductId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err})
				return
			}
			for result1.Next() {
				var imageTemp models.Image
				err := result1.Scan(&imageTemp.Url, &imageTemp.Description)
				if err != nil {
					ctx.JSON(500, gin.H{"error": err.Error()})
					return
				}
				imageProduct = append(imageProduct, imageTemp)
			}
			product.ProductImages = imageProduct
			products = append(products, product)
		}
		ctx.IndentedJSON(200, products)
	}
}
