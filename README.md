# search-service
Микросервис поиска товаров с использованием Elasticsearch

🚀 Возможности

Поиск по названию - полнотекстовый поиск с поддержкой нечеткого поиска
1) Фильтрация по цене - поиск в диапазоне цен
2) Пагинация - постраничный вывод результатов
3) Сортировка - по релевантности и названию
4) Валидация - проверка входных параметров

API Endpoints

1. Поиск товаров
GET /search-service/products/search

Параметры:
- name (string) - название товара для поиска
- min_price (int) - минимальная цена
- max_price (int) - максимальная цена  
- page (int) - номер страницы (по умолчанию: 1)
- limit (int) - количество элементов на странице (по умолчанию: 20, максимум: 100)

Примеры запросов:

Поиск по названию
GET /search-service/products/search?name=футболка

Поиск с фильтром по цене
GET /search-service/products/search?name=джинсы&min_price=1000&max_price=5000

Пагинация
GET /search-service/products/search?page=2&limit=10

Получить все товары
GET /search-service/products/search

2. Получить все товары
GET /search-service/products/

Параметры:
- page (int) - номер страницы (по умолчанию: 1)
- limit (int) - количество элементов на странице (по умолчанию: 20)

📋 Формат ответа
{
  "products": [
    {
      "id": 1,
      "name": "Футболка хлопковая",
      "description": "Мягкая хлопковая футболка",
      "category_id": 1,
      "brand_id": 1,
      "variants": [
        {
          "variant_id": "var_1",
          "sku": "T-SHIRT-001",
          "price": 1500.0,
          "sizes": [42, 44, 46],
          "colors": ["белый", "черный"],
          "material": "хлопок",
          "stock": 100,
          "rating": 4.5
        }
      ]
    }
  ],
  "total": 150,
  "page": 1,
  "limit": 20,
  "pages": 8
}

🛠️ Запуск
docker-compose up --build