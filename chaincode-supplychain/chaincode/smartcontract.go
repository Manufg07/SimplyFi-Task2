package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract implements the supply chain tracking logic
type SmartContract struct {
	contractapi.Contract
}

// StatusEntry records each status change with timestamp
type StatusEntry struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// Product represents a tracked item in the supply chain
type Product struct {
	ProductID     string        `json:"productID"`
	CurrentStatus string        `json:"currentStatus"`
	StatusHistory []StatusEntry `json:"statusHistory"`
}

// Valid status transitions for business logic enforcement
var validStatuses = map[string]struct{}{
	"Manufactured": {},
	"Shipped":      {},
	"In-Transit":   {},
	"Delivered":    {},
}

// Added validation for status transitions
// MODIFICATION: Added status transition rules to enforce valid state machine logic
var validTransitions = map[string][]string{
	"Manufactured": {"Shipped"},
	"Shipped":      {"In-Transit"},
	"In-Transit":   {"Delivered"},
	"Delivered":    {}, // No further transitions allowed
}

// RegisterProduct initializes a new product with Manufactured status
func (s *SmartContract) RegisterProduct(ctx contractapi.TransactionContextInterface, productID string) error {
	if productID == "" {
		return fmt.Errorf("product ID cannot be empty")
	}

	exists, err := s.productExists(ctx, productID)
	if err != nil {
		return fmt.Errorf("existence check failed: %v", err)
	}
	if exists {
		return fmt.Errorf("product %s already exists", productID)
	}

	timestamp, err := s.getTxTimestamp(ctx)
	if err != nil {
		return fmt.Errorf("timestamp error: %v", err)
	}

	product := Product{
		ProductID:     productID,
		CurrentStatus: "Manufactured",
		StatusHistory: []StatusEntry{
			{
				Status:    "Manufactured",
				Timestamp: timestamp,
			},
		},
	}

	productBytes, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("marshaling error: %v", err)
	}

	return ctx.GetStub().PutState(productID, productBytes)
}

// UpdateStatus transitions product to new state with validation
func (s *SmartContract) UpdateStatus(ctx contractapi.TransactionContextInterface, productID string, newStatus string) error {
	if productID == "" {
		return fmt.Errorf("product ID cannot be empty")
	}

	if _, valid := validStatuses[newStatus]; !valid {
		return fmt.Errorf("invalid status: %s", newStatus)
	}

	product, err := s.getProduct(ctx, productID)
	if err != nil {
		return err
	}

	if product.CurrentStatus == newStatus {
		return fmt.Errorf("product already in status: %s", newStatus)
	}

	// MODIFICATION: Added transition validation to enforce proper status flow
	if !isValidTransition(product.CurrentStatus, newStatus) {
		return fmt.Errorf("invalid transition from %s to %s", product.CurrentStatus, newStatus)
	}

	timestamp, err := s.getTxTimestamp(ctx)
	if err != nil {
		return fmt.Errorf("timestamp error: %v", err)
	}

	product.CurrentStatus = newStatus
	product.StatusHistory = append(product.StatusHistory, StatusEntry{
		Status:    newStatus,
		Timestamp: timestamp,
	})

	productBytes, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("marshaling error: %v", err)
	}

	return ctx.GetStub().PutState(productID, productBytes)
}

// GetProduct retrieves complete product history
func (s *SmartContract) GetProduct(ctx contractapi.TransactionContextInterface, productID string) (*Product, error) {
	if productID == "" {
		return nil, fmt.Errorf("product ID cannot be empty")
	}

	productBytes, err := ctx.GetStub().GetState(productID)
	if err != nil {
		return nil, fmt.Errorf("ledger read error: %v", err)
	}
	if productBytes == nil {
		return nil, fmt.Errorf("product %s not found", productID)
	}

	var product Product
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("unmarshaling error: %v", err)
	}

	return &product, nil
}

// MODIFICATION: Added dedicated function to get only status history for clear separation of concerns
// GetProductHistory retrieves only the status timeline of a product
func (s *SmartContract) GetProductHistory(ctx contractapi.TransactionContextInterface, productID string) ([]StatusEntry, error) {
	product, err := s.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}
	return product.StatusHistory, nil
}

// GetAllProducts returns all registered products with their statuses
func (s *SmartContract) GetAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("ledger scan failed: %v", err)
	}
	defer resultsIterator.Close()

	var products []*Product
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("iterator error: %v", err)
		}

		var product Product
		if err := json.Unmarshal(queryResult.Value, &product); err != nil {
			return nil, fmt.Errorf("unmarshaling error: %v", err)
		}
		products = append(products, &product)
	}

	return products, nil
}

// MODIFICATION: Added query function to filter products by status
// GetProductsByStatus retrieves all products with a specific status
func (s *SmartContract) GetProductsByStatus(ctx contractapi.TransactionContextInterface, status string) ([]*Product, error) {
	if _, valid := validStatuses[status]; !valid {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	allProducts, err := s.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	var filteredProducts []*Product
	for _, product := range allProducts {
		if product.CurrentStatus == status {
			filteredProducts = append(filteredProducts, product)
		}
	}

	return filteredProducts, nil
}

// productExists checks for product existence
func (s *SmartContract) productExists(ctx contractapi.TransactionContextInterface, productID string) (bool, error) {
	productBytes, err := ctx.GetStub().GetState(productID)
	if err != nil {
		return false, fmt.Errorf("ledger read error: %v", err)
	}
	return productBytes != nil, nil
}

// getProduct internal helper with basic validation
func (s *SmartContract) getProduct(ctx contractapi.TransactionContextInterface, productID string) (*Product, error) {
	product, err := s.GetProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product fetch failed: %v", err)
	}
	return product, nil
}

// getTxTimestamp converts transaction timestamp to RFC3339 format
func (s *SmartContract) getTxTimestamp(ctx contractapi.TransactionContextInterface) (string, error) {
	txTime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("timestamp unavailable: %v", err)
	}
	return time.Unix(txTime.Seconds, int64(txTime.Nanos)).UTC().Format(time.RFC3339), nil
}

// MODIFICATION: Added helper function to validate state transitions
// isValidTransition checks if a status transition is valid according to business rules
func isValidTransition(currentStatus, newStatus string) bool {
	allowedNextStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, status := range allowedNextStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}

