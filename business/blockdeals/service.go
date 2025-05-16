package v1

import (
	"database/sql"
	"space/models"

	"github.com/gin-gonic/gin"
)

type blockDealService struct {
	DB *sql.DB
}

// NewBlockDealService returns a new BlockDealService instance
func NewBlockDealService(db *sql.DB) models.BlockDealService {
	return &blockDealService{DB: db}
}

// Create a new BlockDeal
func (s *blockDealService) CreateBlockDeal(c *gin.Context, blockDeal models.BlockDeal) error {
	_, err := s.DB.Exec(
		"INSERT INTO blockdeals (cocode, dealtype, scripcode, serial, date1, scripname, clientname, buysell, qtyshares, avgprice, unixtime) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		blockDeal.Cocode, blockDeal.DealType, blockDeal.ScripCode, blockDeal.Serial, blockDeal.Date1, blockDeal.ScripName, blockDeal.ClientName, blockDeal.BuySell, blockDeal.QtyShares, blockDeal.AvgPrice, blockDeal.UnixTime,
	)
	return err
}

// Get all BlockDeals
func (s *blockDealService) GetAllBlockDeals(c *gin.Context) ([]models.BlockDeal, error) {
	rows, err := s.DB.Query("SELECT cocode, dealtype, scripcode, serial, date1, scripname, clientname, buysell, qtyshares, avgprice, unixtime FROM blockdeals")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blockDeals []models.BlockDeal
	for rows.Next() {
		var blockDeal models.BlockDeal
		if err := rows.Scan(&blockDeal.Cocode, &blockDeal.DealType, &blockDeal.ScripCode, &blockDeal.Serial, &blockDeal.Date1, &blockDeal.ScripName, &blockDeal.ClientName, &blockDeal.BuySell, &blockDeal.QtyShares, &blockDeal.AvgPrice, &blockDeal.UnixTime); err != nil {
			return nil, err
		}
		blockDeals = append(blockDeals, blockDeal)
	}

	return blockDeals, nil
}

// Get a single BlockDeal by ID
func (s *blockDealService) GetBlockDealByID(c *gin.Context, id int) (models.BlockDeal, error) {
	var blockDeal models.BlockDeal
	err := s.DB.QueryRow(
		"SELECT cocode, dealtype, scripcode, serial, date1, scripname, clientname, buysell, qtyshares, avgprice, unixtime FROM blockdeals WHERE cocode = $1",
		id,
	).Scan(&blockDeal.Cocode, &blockDeal.DealType, &blockDeal.ScripCode, &blockDeal.Serial, &blockDeal.Date1, &blockDeal.ScripName, &blockDeal.ClientName, &blockDeal.BuySell, &blockDeal.QtyShares, &blockDeal.AvgPrice, &blockDeal.UnixTime)

	return blockDeal, err
}

// Update a BlockDeal
func (s *blockDealService) UpdateBlockDeal(cocode int, blockDeal models.BlockDeal) error {
	_, err := s.DB.Exec(
		"UPDATE blockdeals SET dealtype=$1, scripcode=$2, serial=$3, date1=$4, scripname=$5, clientname=$6, buysell=$7, qtyshares=$8, avgprice=$9, unixtime=$10 WHERE cocode=$11",
		blockDeal.DealType, blockDeal.ScripCode, blockDeal.Serial, blockDeal.Date1, blockDeal.ScripName, blockDeal.ClientName, blockDeal.BuySell, blockDeal.QtyShares, blockDeal.AvgPrice, blockDeal.UnixTime, cocode,
	)
	return err
}

// Delete a BlockDeal
func (s *blockDealService) DeleteBlockDeal(c *gin.Context, id int) error {
	_, err := s.DB.Exec("DELETE FROM blockdeals WHERE cocode = $1", id)
	return err
}
