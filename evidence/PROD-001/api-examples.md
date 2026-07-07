# PROD-001 API Examples

Smoke check ran against local backend `http://localhost:18080` and migrated `marketplace_dev`. JWT values were used only in memory and are not recorded.

## Create Draft As Seller A

`POST /api/v1/seller/products`

```json
{
  "status": 201,
  "body": {
    "id": "85caf760-d87d-4c68-a896-f6ac542b445c",
    "sellerId": "ee852738-5e54-4f3c-bd78-8f6001f46fe2",
    "name": "Smoke product 1783357170491",
    "description": "Smoke description",
    "status": "draft",
    "createdAt": "2026-07-06T16:59:30Z",
    "updatedAt": "2026-07-06T16:59:30Z"
  }
}
```

## Public List Before Publish

`GET /api/v1/products`

```json
{
  "containsSmokeProduct": false
}
```

## Non-Owner Publish

`POST /api/v1/seller/products/{productID}/publish`

```json
{
  "status": 404,
  "body": {
    "error": {
      "code": "not_found",
      "message": "Product not found"
    }
  }
}
```

## Owner Publish

`POST /api/v1/seller/products/{productID}/publish`

```json
{
  "status": 200,
  "body": {
    "id": "85caf760-d87d-4c68-a896-f6ac542b445c",
    "sellerId": "ee852738-5e54-4f3c-bd78-8f6001f46fe2",
    "name": "Smoke product 1783357170491",
    "description": "Smoke description",
    "status": "published",
    "createdAt": "2026-07-06T16:59:30Z",
    "updatedAt": "2026-07-06T16:59:30Z"
  }
}
```

## Public List After Publish

`GET /api/v1/products`

```json
{
  "status": 200,
  "product": {
    "id": "85caf760-d87d-4c68-a896-f6ac542b445c",
    "name": "Smoke product 1783357170491",
    "description": "Smoke description",
    "status": "published",
    "createdAt": "2026-07-06T16:59:30Z",
    "updatedAt": "2026-07-06T16:59:30Z"
  }
}
```

## Public Detail

`GET /api/v1/products/{productID}`

```json
{
  "status": 200,
  "body": {
    "id": "85caf760-d87d-4c68-a896-f6ac542b445c",
    "name": "Smoke product 1783357170491",
    "description": "Smoke description",
    "status": "published",
    "createdAt": "2026-07-06T16:59:30Z",
    "updatedAt": "2026-07-06T16:59:30Z"
  }
}
```

## Repeat Publish

`POST /api/v1/seller/products/{productID}/publish`

```json
{
  "status": 409,
  "body": {
    "error": {
      "code": "already_published",
      "message": "Product is already published"
    }
  }
}
```
