package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset structure
type Asset struct {
	AssetID string `json:"assetID"`
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	Value   int    `json:"value"` 
}

// CreateAsset 
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetID, name, owner, value string) error {
	if assetID == "" || name == "" || owner == "" || value == "" {
		return fmt.Errorf("all fields must be non-empty")
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("value must be a valid integer")
	}

	exists, err := s.AssetExists(ctx, assetID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("asset %s already exists", assetID)
	}

	asset := Asset{
		AssetID: assetID,
		Name:    name,
		Owner:   owner,
		Value:   valueInt,
	}

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset: %v", err)
	}

	err = ctx.GetStub().PutState(assetID, assetBytes)
	if err != nil {
		return fmt.Errorf("failed to put asset in world state: %v", err)
	}

	// Store composite key for querying by owner
	compositeKey, err := ctx.GetStub().CreateCompositeKey("owner~assetID", []string{owner, assetID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}
	err = ctx.GetStub().PutState(compositeKey, []byte{0}) 
	if err != nil {
		return fmt.Errorf("failed to store composite key: %v", err)
	}

	// Emit event 
	err = ctx.GetStub().SetEvent("AssetCreated", assetBytes)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// ReadAsset retrieves an asset by ID
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, assetID string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset %s: %v", assetID, err)
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", assetID)
	}

	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset 
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, assetID, name, owner, value string) error {
    exists, err := s.AssetExists(ctx, assetID)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the asset %s does not exist", assetID)
    }

    valueInt, err := strconv.Atoi(value)
    if err != nil {
        return fmt.Errorf("value must be a valid integer")
    }

    asset := Asset{
        AssetID: assetID,
        Name:    name,
        Owner:   owner,
        Value:   valueInt,
    }

    assetBytes, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState(assetID, assetBytes)
    if err != nil {
        return fmt.Errorf("failed to update asset: %v", err)
    }

    err = ctx.GetStub().SetEvent("AssetUpdated", assetBytes)
    if err != nil {
        return fmt.Errorf("failed to emit event: %v", err)
    }

    return nil
}

// DeleteAsset 
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
    exists, err := s.AssetExists(ctx, assetID)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the asset %s does not exist", assetID)
    }

    //  read the asset to get the owner for composite key deletion
    existingAsset, err := s.ReadAsset(ctx, assetID)
    if err != nil {
        return err
    }

    err = ctx.GetStub().DelState(assetID)
    if err != nil {
        return fmt.Errorf("failed to delete asset: %v", err)
    }

    // Remove composite key
    compositeKey, err := ctx.GetStub().CreateCompositeKey("owner~assetID", []string{existingAsset.Owner, assetID})
    if err != nil {
        return fmt.Errorf("failed to create composite key: %v", err)
    }
    err = ctx.GetStub().DelState(compositeKey)
    if err != nil {
        return fmt.Errorf("failed to delete composite key: %v", err)
    }

    err = ctx.GetStub().SetEvent("AssetDeleted", []byte(assetID))
    if err != nil {
        return fmt.Errorf("failed to emit event: %v", err)
    }

    return nil
}

// AssetExists returns true when an asset with given ID exists
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {
	assetBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetBytes != nil, nil
}

// QueryAssetsByOwner retrieves all assets owned by a specific owner
func (s *SmartContract) QueryAssetsByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("owner~assetID", []string{owner})
	if err != nil {
		return nil, fmt.Errorf("failed to get assets by owner: %v", err)
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// Extract assetID from composite key
		_, keyParts, err := ctx.GetStub().SplitCompositeKey(queryResult.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key: %v", err)
		}
		assetID := keyParts[1] 

		asset, err := s.ReadAsset(ctx, assetID)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}
