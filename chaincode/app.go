package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Service struct {
	Key          string
	TypeServices string `json:"type_services"`       // Тип услуги
	Comment      string `json:"comment,omitempty"`   // Комментарий
	Email        string `json:"email"`               // Рабочий E-mail адрес
	FirstName    string `json:"first_name"`          // Фамилия
	LastName     string `json:"last_name,omitempty"` // Имя
	Phone        string `json:"phone"`               // Телефон
	Address      string `json:"address"`             // Адрес
	Status       int    `json:"status"`
	ClientID     string
	MspID        string
}

const (
	mspAdmin   = "Org1MSP"
	mspUser    = "Org2MSP"
	keyByRange = "service"
)

const (
	statusOpen = iota + 1
	statusConsideration
	statusDataRefinement
	statusSatisfied
	statusWithdrawn
	statusRefusal
	statusDelete
)

// getServiceByUUID Для просмотра конкретной заявки по её uuid
func (s *SmartContract) GetServiceByUUID(ctx contractapi.TransactionContextInterface, typeServices, phone string) (string, error) {
	// Создание композитного ключа
	keyService, err := ctx.GetStub().CreateCompositeKey("service", []string{typeServices, phone})
	if err != nil {
		return "", fmt.Errorf("failed to create composite key: %v", err)
	}

	// Получение данных по ключу
	serviceByte, err := ctx.GetStub().GetState(keyService)
	if err != nil {
		return "", fmt.Errorf("failed get service for composite key: %v", err)
	}

	return string(serviceByte), nil
}

// createService Для создания новых заявок
func (s *SmartContract) СreateService(ctx contractapi.TransactionContextInterface, data string) error {
	inputData := &Service{}
	err := json.Unmarshal([]byte(data), &inputData)
	if err != nil {
		fmt.Printf("error unmarshal: %v\n", err)
		return fmt.Errorf("unmarshall error: %v", err)
	}

	// Создание композитного ключа
	keyService, err := ctx.GetStub().CreateCompositeKey("service", []string{inputData.TypeServices, inputData.Phone})
	if err != nil {
		fmt.Printf("failed to create composite key: %v\n", err)
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	// Проверка существования ключа и данных
	exists, err := s.serviceExists(ctx, keyService)
	if err != nil {
		fmt.Printf("failed to get asset: %v\n", err)
		return fmt.Errorf("failed to get asset: %v", err)
	}
	if exists {
		fmt.Printf("asset already exists: %s\n", keyService)
		return fmt.Errorf("asset already exists: %s", keyService)
	}

	// получить ID отправившего клиента
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		fmt.Printf("failed to get client identity %v\n", err)
		return fmt.Errorf("failed to get client identity %v", err)
	}
	// получить ID организации
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		fmt.Printf("failed to get org identity %v\n", err)
		return fmt.Errorf("failed to get org identity %v", err)
	}

	// Заполнение данных
	inputData.Status = statusOpen
	inputData.ClientID = clientID
	inputData.MspID = mspID

	projectJSON, err := json.Marshal(inputData)
	if err != nil {
		fmt.Printf("marshal error %v\n", err)
		return fmt.Errorf("marshal error %v", err)
	}

	// Загружаем проект в БЧ
	err = ctx.GetStub().PutState(keyService, projectJSON)
	if err != nil {
		fmt.Printf("failed to put auction in public data: %v\n", err)
		return fmt.Errorf("failed to put auction in public data: %v", err)
	}

	return nil
}

// setServiceStatus Для изменения статуса заявки с нужным uuid
func (s *SmartContract) SetServiceStatus(ctx contractapi.TransactionContextInterface, typeServices, phone, status string) error {
	// Создание композитного ключа
	keyService, err := ctx.GetStub().CreateCompositeKey("service", []string{typeServices, phone})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	// проверка существования ключа и данных по нему
	exists, err := s.serviceExists(ctx, keyService)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	if !exists {
		return fmt.Errorf("asset not exists: %s", keyService)
	}

	// Идентификация по ID организации
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get org identity %v", err)
	}
	if mspID != mspAdmin {
		return fmt.Errorf("premission denied %v", err)
	}

	// Получение данных по ключу
	serviceByte, err := ctx.GetStub().GetState(keyService)
	if err != nil {
		return fmt.Errorf("failed get service for composite key: %v", err)
	}

	service := &Service{}
	err = json.Unmarshal(serviceByte, service)
	if err != nil {
		return fmt.Errorf("unmarshall error: %v", err)
	}

	i, err := strconv.Atoi(status)
	if err != nil {
		return fmt.Errorf("error converting to int: %v", err)
	}

	// Заполнение данных
	service.Status = i

	projectJSON, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("marshal error %v", err)
	}

	// Загружаем проект в БЧ
	err = ctx.GetStub().PutState(keyService, projectJSON)
	if err != nil {
		return fmt.Errorf("failed to put auction in public data: %v", err)
	}

	return nil
}

// setServiceStatus Для изменения статуса заявки с нужным uuid
func (s *SmartContract) WithDrawService(ctx contractapi.TransactionContextInterface, typeServices, phone, status string) error {
	// Создание композитного ключа
	keyService, err := ctx.GetStub().CreateCompositeKey("service", []string{typeServices, phone})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	// проверка на существование ключа и данных по нему
	exists, err := s.serviceExists(ctx, keyService)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	if !exists {
		return fmt.Errorf("asset not exists: %s", keyService)
	}

	// получить ID отправившего клиента
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		fmt.Printf("failed to get client identity %v\n", err)
		return fmt.Errorf("failed to get client identity %v", err)
	}
	//получение ID организации
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get org identity %v", err)
	}

	// получение данных по ключу
	serviceByte, err := ctx.GetStub().GetState(keyService)
	if err != nil {
		return fmt.Errorf("failed get service for composite key: %v", err)
	}

	service := &Service{}
	err = json.Unmarshal(serviceByte, service)
	if err != nil {
		return fmt.Errorf("unmarshall error: %v", err)
	}

	// Идентификация пользователя и проверка не удалена ли запись
	if mspID != service.MspID && clientID != service.ClientID {
		return fmt.Errorf("premission denied. mspID: %s not equal owner", mspID)
	}
	if service.Status != statusDelete {
		return fmt.Errorf("service request has already been deleted: %s", "DELETED")
	}

	service.Status = statusWithdrawn

	projectJSON, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("marshal error %v", err)
	}

	// Загружаем проект в БЧ
	err = ctx.GetStub().PutState(keyService, projectJSON)
	if err != nil {
		return fmt.Errorf("failed to put auction in public data: %v", err)
	}

	return nil
}

// setServiceStatus Для изменения статуса заявки с нужным uuid
func (s *SmartContract) DeleteService(ctx contractapi.TransactionContextInterface, typeServices, phone, status string) error {
	// Создание композитного ключа
	keyService, err := ctx.GetStub().CreateCompositeKey("service", []string{typeServices, phone})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	// проверка существования данных по ключу
	exists, err := s.serviceExists(ctx, keyService)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	if !exists {
		return fmt.Errorf("asset not exists: %s", keyService)
	}

	// получение ID организации и аутентификации по нему
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get org identity %v", err)
	}
	if mspID != mspAdmin {
		return fmt.Errorf("premission denied %v", err)
	}

	// получение данных по ключу
	serviceByte, err := ctx.GetStub().GetState(keyService)
	if err != nil {
		return fmt.Errorf("failed get service for composite key: %v", err)
	}

	service := &Service{}
	err = json.Unmarshal(serviceByte, service)
	if err != nil {
		return fmt.Errorf("unmarshall error: %v", err)
	}

	// Устанавливаем статус удален
	service.Status = statusDelete

	projectJSON, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("marshal error %v", err)
	}

	// Загружаем проект в БЧ
	err = ctx.GetStub().PutState(keyService, projectJSON)
	if err != nil {
		return fmt.Errorf("failed to put auction in public data: %v", err)
	}

	return nil
}
