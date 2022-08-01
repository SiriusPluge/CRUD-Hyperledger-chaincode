package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// AssetExists returns true when asset with given ID exists in the ledger.
func (s *SmartContract) serviceExists(ctx contractapi.TransactionContextInterface, serviceID string) (bool, error) {
	projectBytes, err := ctx.GetStub().GetState(serviceID)
	if err != nil {
		return false, fmt.Errorf("failed to read asset %s from world state. %v", serviceID, err)
	}

	return projectBytes != nil, nil
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Service, error) {
	var services []*Service
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var service Service
		err = json.Unmarshal(queryResult.Value, &service)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}

	return services, nil
}

func contains(sli []string, str string) bool {
	for _, a := range sli {
		if a == str {
			return true
		}
	}
	return false
}
