# search-service

Микросервис поиска товаров с использованием Elasticsearch через GraphQL.

## 🚀 Возможности

- Поиск товаров по названию с поддержкой нечеткого поиска
- Фильтрация по цене, размеру, цвету, материалу, категории, бренду и складу
- Пагинация и сортировка по релевантности
- Валидация входных данных

---

## 📡 GraphQL API

### Query: `searchProducts`

**Описание:** поиск товаров с фильтрацией и пагинацией.

**Входные параметры (`SearchInput`):**

| Поле        | Тип         | Описание |
|------------|------------|----------|
| `name`      | String     | Название товара для поиска |
| `productID` | Int        | ID конкретного товара |
| `variantID` | Int        | ID варианта товара |
| `sku`       | String     | SKU варианта товара |
| `material`  | String     | Материал варианта |
| `color`     | String     | Цвет варианта |
| `size`      | String     | Размер варианта |
| `stock`     | Int        | Минимальный запас на складе |
| `categoryID`| Int        | ID категории |
| `brandID`   | Int        | ID бренда |
| `minPrice`  | Float      | Минимальная цена |
| `maxPrice`  | Float      | Максимальная цена |
| `page`      | Int        | Номер страницы (по умолчанию 1) |
| `limit`     | Int        | Количество элементов на странице (по умолчанию 20, максимум 100) |

---

### Пример запроса

```graphql
query {
  searchProducts(input: {
    name: "футболка",
    minPrice: 500,
    maxPrice: 2000,
    page: 1,
    limit: 10
  }) {
    products {
      id
      name
      description
      categoryID
      brandID
      variants {
        variantID
        sku
        price
        sizes
        colors
        material
        stock
        rating
      }
    }
    total
    page
    limit
    pages
  }
}


{
  "data": {
    "searchProducts": {
      "products": [
        {
          "id": 1,
          "name": "Футболка хлопковая",
          "description": "Мягкая хлопковая футболка",
          "categoryID": 1,
          "brandID": 1,
          "variants": [
            {
              "variantID": "var_1",
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
  }
}


🛠️ Запуск
docker-compose up --build

Примечания:

Все фильтры в SearchInput опциональны. Если не передавать параметры, возвращаются все товары с пагинацией.

page и limit задают постраничный вывод. Максимальный limit — 100.

Результаты сортируются по релевантности (_score) в Elasticsearch.