# search-service

Микросервис поиска товаров с использованием Elasticsearch через GraphQL.

## 🚀 Возможности

- Поиск товаров по названию с поддержкой нечеткого поиска
- Фильтрация по цене, размеру, цвету, материалу, категории, бренду и складу
- Пагинация и сортировка по релевантности
- Валидация входных данных

---

## 📡 GraphQL API

Query: `searchProducts`

Описание: поиск товаров с фильтрацией и пагинацией.

для генерации файлов надо зайти в internal и сделать go run github.com/99designs/gqlgen generate 

для синхронизации Elastic Search с базой данных команда

```
docker compose run --rm logstash -f /usr/share/logstash/logstash_internal/postgres_to_es.conf
```

Входные параметры (`SearchInput`)

```
query Search($in: SearchInput!) {
  searchProducts(input: $in) {
    total
    page
    limit
    pages
    products {
      id
      name
      description
      material
      category_id
      brand_id
      is_active
      image_urls
      video_urls
      variants {
        variant_id
        sku
        price
        sizes
        colors
        stock
        rating
        image_urls
      }
    }
  }
}
```


---

```
{
  "in": {
    "limit": 10,
    "page": 1
  }
}
```


это просто описание структуры
```
{
  "data": {
    "searchProducts": {
      "products": [
        {
          "id": 14,
          "name": "Жилет «Utility Vest»",
          "description": "Многофункциональный жилет с карманами, для городской среды",
          "material": "канвас",
          "category_id": 4,
          "brand_id": 4,
          "is_active": true,
          "image_urls": [
            "http://localhost:8084/media/vest1.jpg",
            "http://localhost:8084/media/vest2.jpg"
          ],
          "video_urls": [],
          "variants": [
            {
              "variant_id": 7,
              "sku": "UTV-VE-OLIVE-M",
              "price": 54.99,
              "sizes": "M",
              "colors": "Олива",
              "stock": 27,
              "rating": 0,
              "image_urls": []
            }
          ]
        },
      ]
      "total": 150,
      "page": 1,
      "limit": 20,
      "pages": 8
    }
  }
}
```

```
{
  "in": {
    "name": "Shirt",
    "minPrice": 1000,
    "maxPrice": 5000,
    "page": 1,
    "limit": 10
  }
}
```

```
{
  "in": {
    "material": "cotton",
    "color": "red",
    "size": "XL"
  }
}
```


```
{
  "in": {
    "brandID": 5,
    "categoryID": 12
  }
}
```


```
{
  "in": {
    "sku": "TSHIRT-BLK-001",
    "stock": 10
  }
}
```


🛠️ Запуск
docker-compose up --build

Примечания:

Все фильтры в SearchInput опциональны. Если не передавать параметры, возвращаются все товары с пагинацией.

page и limit задают постраничный вывод. Максимальный limit — 100.

Результаты сортируются по релевантности (_score) в Elasticsearch.