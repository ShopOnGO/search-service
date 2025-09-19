package search

import (
	"encoding/json"
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
)

func HandleProductEvent(msg []byte, key string, searchSvc *SearchService) error {
	logger.Infof("Получено сообщение: %s", string(msg))

	var event ProductCreatedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
	}

	eventHandlers := map[string]func([]byte, *SearchService) error{
		"create": HandleCreateProductEvent,
		// "update": HandleUpdateProductEvent,
		// "delete": HandleDeleteProductEvent,
	}

	handler, exists := eventHandlers[event.Action]
	if !exists {
		return fmt.Errorf("неизвестное действие для продукта: %s", event.Action)
	}

	return handler(msg, searchSvc)
}

func HandleCreateProductEvent(msg []byte, searchSvc *SearchService) error {
	var event ProductCreatedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
	}

	logger.Infof("Получены данные для индексации продукта: %+v", event)

	if err := searchSvc.AddProductToElasticSearch(event); err != nil {
		logger.Errorf("Ошибка при индексации продукта: %v", err)
		return err
	}

	logger.Infof("Продукт успешно индексирован: %+v", event)
	return nil
}