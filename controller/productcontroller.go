package controller

import (
	"fmt"
	"strings"

	"github.com/PhanLuc1/tech-heim-backend/database"
	"github.com/PhanLuc1/tech-heim-backend/models"
	"github.com/gin-gonic/gin"
)

func GetImageProduct(idproduct int) (imageProduct []models.Image, err error) {
	result1, err := database.Client.Query("SELECT productimage.url,productimage.description FROM productimage WHERE idProduct = ?", idproduct)
	if err != nil {
		return nil, err
	}
	for result1.Next() {
		var imageTemp models.Image
		err := result1.Scan(&imageTemp.Url, &imageTemp.Description)
		if err != nil {
			return nil, err
		}
		imageProduct = append(imageProduct, imageTemp)
	}
	return imageProduct, nil
}
func GetTechnicalProduct(idproduct int) (productTechnical []models.Technical, err error) {
	result, err := database.Client.Query("SELECT producttechnical.title,producttechnical.description FROM producttechnical WHERE idProduct = ?", idproduct)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		var technicalTemp models.Technical
		err := result.Scan(&technicalTemp.Title, &technicalTemp.Description)
		if err != nil {
			return nil, err
		}
		productTechnical = append(productTechnical, technicalTemp)
	}
	return productTechnical, nil
}
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
			var idCategory int
			err := result.Scan(&product.ProductId, &product.ProductName, &product.Rate, &product.Sold, &product.Quantity, &product.CurrentPrice, &product.LastPrice, &idCategory)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			product.ProductImages, err = GetImageProduct(*product.ProductId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			products = append(products, product)
		}
		ctx.IndentedJSON(200, products)
	}
}
func GetProductDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		var product models.Product
		query := "SELECT product.id,product.name,product.rate,product.quantity,product.quantity,product.currentPrice,product.lastPrice FROM product WHERE id =?"
		err := database.Client.QueryRow(query, id).Scan(
			&product.ProductId,
			&product.ProductName,
			&product.Rate,
			&product.Sold,
			&product.Quantity,
			&product.CurrentPrice,
			&product.LastPrice,
		)
		if err != nil {
			ctx.JSON(500, gin.H{"Error": err.Error()})
		}
		product.ProductImages, _ = GetImageProduct(*product.ProductId)
		product.ProuctTechnical, _ = GetTechnicalProduct(*product.ProductId)
		// get Group and type product
		var productGroups []models.Group
		query = "SELECT DISTINCT productgroup.id,productgroup.title FROM productgroup	JOIN producttype ON producttype.idGroup = productgroup.id JOIN product_producttype ON product_producttype.idType = producttype.id JOIN product ON product_producttype.idProduct = product.id WHERE product.id = ? "
		query += " ORDER BY productgroup.id ASC "
		result2, err := database.Client.Query(query, *product.ProductId)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for result2.Next() {
			var productGroup models.Group
			var types []models.Type
			result2.Scan(&productGroup.Id, &productGroup.Title)
			result3, err := database.Client.Query("SELECT producttype.id,producttype.title,producttype.description FROM producttype JOIN product_producttype ON product_producttype.idType = producttype.id JOIN product ON product.id = product_producttype.idProduct WHERE idGroup = ? AND idProduct = ?", productGroup.Id, product.ProductId)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			for result3.Next() {
				var typeTemp models.Type
				err := result3.Scan(&typeTemp.Id, &typeTemp.Title, &typeTemp.Description)
				if err != nil {
					ctx.JSON(200, gin.H{"Erorr": err})
					return
				}
				types = append(types, typeTemp)
			}

			productGroup.Type = types
			productGroups = append(productGroups, productGroup)
		}
		product.ProductGroup = productGroups
		// get Commentproduct
		var productComments []models.Comment
		query = "SELECT user.firstName,user.lastName,productcomment.description FROM user JOIN productcomment ON productcomment.idUser = user.id JOIN product ON product.id = productcomment.idProduct WHERE product.id = ?"
		result4, err := database.Client.Query(query, *product.ProductId)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
		}
		for result4.Next() {
			var productComment models.Comment
			err := result4.Scan(&productComment.FirstName, &productComment.LastName, &productComment.Description)
			if err != nil {
				ctx.JSON(404, gin.H{"error": err.Error()})
			}
			productComments = append(productComments, productComment)
		}
		product.ProductComment = productComments

		ctx.JSON(200, product)
	}
}
func GetProductGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var productGroups []models.Group
		categoryId := ctx.Query("categoryId")
		query := "SELECT productgroup.id,productgroup.title FROM productgroup WHERE productgroup.categoryId = " + categoryId
		result, err := database.Client.Query(query)
		if err != nil {
			ctx.JSON(500, gin.H{"Error": err.Error()})
			return
		}
		for result.Next() {
			var productGroup models.Group
			var types []models.Type
			result.Scan(&productGroup.Id, &productGroup.Title)
			result1, err := database.Client.Query("SELECT producttype.id,producttype.title,producttype.description FROM producttype WHERE idGroup = ?", productGroup.Id)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			for result1.Next() {
				var typeTemp models.Type
				err := result1.Scan(&typeTemp.Id, &typeTemp.Title, &typeTemp.Description)
				if err != nil {
					ctx.JSON(200, gin.H{"Erorr": err.Error()})
					return
				}
				types = append(types, typeTemp)
			}

			productGroup.Type = types
			productGroups = append(productGroups, productGroup)
		}
		ctx.JSON(200, productGroups)
	}
}
