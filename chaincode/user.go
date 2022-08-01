package chaincode

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
)

func (s *SmartContract) Whoami(ctx contractapi.TransactionContextInterface) ([]string, error) {
	// Получения MSPID организации
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, errors.Wrap(err, "get user MSPID error")
	}
	// Получение ID пользователя
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, errors.Wrap(err, "get user ID error")
	}
	return []string{mspID, clientID}, nil
}
