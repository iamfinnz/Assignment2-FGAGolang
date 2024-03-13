package controllers

import (
	"errors"
	"log"
	"net/http"
	"slices"

	"app/database"
	"app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrder(ctx *gin.Context) {
	var (
		newOrder models.Order
		db = database.GetConnection()
	)

	if err := ctx.ShouldBindJSON(&newOrder); err != nil {
		log.Fatalln("error :", err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	trx := db.Begin()
	defer func() {
		if r := recover(); r != nil || ctx.IsAborted(){
			trx.Rollback()
		}
	}()

	if err := trx.Create(&newOrder).Error; err != nil {
		log.Fatalln("error :", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"message": "Server Error",
		})
		return
	}

	trx.Commit()
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"message": "Order berhasil dibuat",
	})
}

func GetAllOrder(ctx *gin.Context) {
	var (
		orders []models.Order
		db = database.GetConnection()
	)

	db.Preload("Items").Find(&orders)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Berhasil mengambil data order",
		"data": orders,
	})
}

func UpdateOrder(ctx *gin.Context) {
	var (
		updateOrder models.Order
		order models.Order
		items []models.Item
		updateIds = []uint{}
		itemIds = []uint{}
		orderId = ctx.Param("orderId")
		db = database.GetConnection()
	)

	if err := ctx.ShouldBindJSON(&updateOrder); err != nil {
		log.Fatalln("error :", err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	trx := db.Begin()
	defer func() {
		if r := recover(); r != nil || ctx.IsAborted(){
			trx.Rollback()
		}
	}()

	// get order data
	if err := trx.First(&order, orderId).Error;
		err != nil && errors.Is(err, gorm.ErrRecordNotFound){

		log.Fatalln("error :", err.Error())
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"message": "Data order tidak ditemukan",
		})
		return
	}

	// update order data
	if err := trx.Model(&order).Updates(updateOrder).Error; err != nil {
		log.Fatalln("error :", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"message": "Server Error",
		})
		return
	}

	for _, item := range updateOrder.Items {
		updateIds = append(updateIds, item.ItemId)
	}

	trx.Where("order_id = ? AND item_id IN ?", orderId, updateIds).Find(&items)
	for _, item := range items {
		itemIds = append(itemIds, item.ItemId)
	}

	for _, item := range updateOrder.Items {
		if slices.Contains(itemIds, item.ItemId) {
			if err := trx.Model(models.Item{}).Where("item_id = ?", item.ItemId).Updates(item).Error;
				err != nil {

				log.Fatalln("error :", err.Error())
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status": "fail",
					"message": "Server Error",
				})
				return
			}
		}
	}

	trx.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Berhasil update data order",
	})
}

func DeleteOrder(ctx *gin.Context) {
	var (
		order models.Order
		orderId = ctx.Param("orderId")
		db = database.GetConnection()
	)

	trx := db.Begin()
	defer func() {
		if r := recover(); r != nil || ctx.IsAborted() {
			trx.Rollback()
		}
	}()

	// get order data
	if err := trx.Preload("Items").First(&order, orderId).Error;
		err != nil && errors.Is(err, gorm.ErrRecordNotFound) {

		log.Fatalln("error :", err.Error())
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"message": "Data order tidak ditemukan",
		})
		return
	}

	// delete items of order
	for _, item := range order.Items {
		if err := trx.Delete(&item).Error; err != nil {
			log.Fatalln("error :", err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": "fail",
				"message": "Server Error",
			})
			return
		}
	}

	// delete order
	if err := trx.Delete(&order).Error; err != nil {
		log.Fatalln("error :", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"message": "Server Error",
		})
		return
	}

	trx.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Berhasil hapus data order",
	})
}