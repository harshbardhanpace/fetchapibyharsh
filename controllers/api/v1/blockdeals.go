package v1

import (
	"net/http"
	"space/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var blockdealsobj models.BlockDealService

func NewBlockDealController(blockDealService models.BlockDealService) {
	blockdealsobj = blockDealService
}

func GetAllBlockDeals(c *gin.Context) {
	blockDeals, err := blockdealsobj.GetAllBlockDeals(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, blockDeals)
}

func CreateBlockDeal(c *gin.Context) {
	var blockDeal models.BlockDeal
	if err := c.ShouldBindJSON(&blockDeal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := blockdealsobj.CreateBlockDeal(c, blockDeal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "BlockDeal created successfully"})
}

// GetBlockDealByCocode retrieves a block deal by its cocode
func GetBlockDealByCocode(c *gin.Context) {
	cocodeStr := c.Param("cocode")
	cocode, err := strconv.Atoi(cocodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cocode"})
		return
	}

	blockDeal, err := blockdealsobj.GetBlockDealByID(c, cocode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Block deal not found"})
		return
	}

	c.JSON(http.StatusOK, blockDeal)
}

// UpdateBlockDeal updates a block deal by its cocode
func UpdateBlockDeal(c *gin.Context) {
	cocodeStr := c.Param("cocode")
	cocode, err := strconv.Atoi(cocodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cocode"})
		return
	}

	var updatedBlockDeal models.BlockDeal
	if err := c.ShouldBindJSON(&updatedBlockDeal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedBlockDeal.Cocode = cocode
	err = blockdealsobj.UpdateBlockDeal(cocode, updatedBlockDeal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update block deal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Block deal updated successfully"})
}

// DeleteBlockDeal deletes a block deal by its cocode
func DeleteBlockDeal(c *gin.Context) {
	cocodeStr := c.Param("cocode")
	cocode, err := strconv.Atoi(cocodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cocode"})
		return
	}

	err = blockdealsobj.DeleteBlockDeal(c, cocode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete block deal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Block deal deleted successfully"})
}
