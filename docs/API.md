# Midtools API Documentation

## Introduction

The Midtools API is a Go service that provides vehicle-related services (brown card, sticker, USSD check, policy verification) and product/risk type management for the Motor Insurance Database (MID) system.

- **Base URL:** `http://localhost:8000`
- **Content-Type:** `application/json` (for all POST requests with a body)
- **Response Format:** All responses are JSON-encoded

## Response Format

### Success Responses

Success responses return the payload directly (an array or object as appropriate) with HTTP status `200 OK`.

### Error Responses

Error responses return a JSON-encoded string containing the error message. The response body is the raw error message (e.g. `"cars is required"`), not a structured object.

| Status Code | Description |
|-------------|-------------|
| 400 | Bad Request - Invalid or missing request body, validation failure |
| 500 | Internal Server Error - Server-side processing error |
| 503 | Service Unavailable - Downstream service (e.g. database, external API) unavailable |

---

## Endpoints

### Vehicle Services

The following endpoints accept a list of Ghana license plate numbers and return results per vehicle. The `cars` field accepts multiple plate numbers separated by **comma**, **newline**, or **tab**. Invalid plates are validated client-side and return an item with `statusCode: false` and a validation message.

---

### POST /browncard

Retrieve brown card information for one or more vehicles.

**Request Body**

```json
{
  "cars": "GR1234-22, GR5678AD"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| cars | string | Yes | Comma, newline, or tab-separated list of Ghana license plate numbers |

**Response** (200 OK)

Array of objects:

```json
[
  {
    "statusCode": true,
    "brownCardNumber": "string",
    "url": "string",
    "message": "string",
    "carNumber": "string"
  }
]
```

| Field | Type | Description |
|-------|------|-------------|
| statusCode | boolean | Whether the lookup succeeded |
| brownCardNumber | string | Brown card number (or validation message if invalid plate) |
| url | string | URL to brown card document (or validation message if invalid plate) |
| message | string | Status or error message |
| carNumber | string | The vehicle registration number |

---

### POST /sticker

Retrieve sticker information for one or more vehicles.

**Request Body**

```json
{
  "cars": "GR1234-22, GR5678AD"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| cars | string | Yes | Comma, newline, or tab-separated list of Ghana license plate numbers |

**Response** (200 OK)

Array of objects:

```json
[
  {
    "statusCode": true,
    "stickerLink": "string",
    "stickerNumber": "string",
    "message": "string",
    "carNumber": "string"
  }
]
```

| Field | Type | Description |
|-------|------|-------------|
| statusCode | boolean | Whether the lookup succeeded |
| stickerLink | string | URL/link to sticker (or validation message if invalid plate) |
| stickerNumber | string | Sticker number (or validation message if invalid plate) |
| message | string | Status or error message |
| carNumber | string | The vehicle registration number |

---

### POST /ussd_check

Perform USSD check for one or more vehicles.

**Request Body**

```json
{
  "cars": "GR1234-22, GR5678AD"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| cars | string | Yes | Comma, newline, or tab-separated list of Ghana license plate numbers |

**Response** (200 OK)

Array of objects:

```json
[
  {
    "statusCode": true,
    "message": "string",
    "carNumber": "string"
  }
]
```

| Field | Type | Description |
|-------|------|-------------|
| statusCode | boolean | Whether the check succeeded |
| message | string | Status or error message |
| carNumber | string | The vehicle registration number |

---

### POST /policy_verification

Verify policy status for one or more vehicles.

**Request Body**

```json
{
  "cars": "GR1234-22, GR5678AD"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| cars | string | Yes | Comma, newline, or tab-separated list of Ghana license plate numbers |

**Response** (200 OK)

Array of objects:

```json
[
  {
    "statusCode": true,
    "ProductName": "string",
    "startDate": "string",
    "endDate": "string",
    "message": "string",
    "carNumber": "string"
  }
]
```

| Field | Type | Description |
|-------|------|-------------|
| statusCode | boolean | Whether the verification succeeded |
| ProductName | string | Name of the insurance product |
| startDate | string | Policy start date |
| endDate | string | Policy end date |
| message | string | Status or error message |
| carNumber | string | The vehicle registration number |

---

### Product Endpoints

---

### POST /products

Seed or create products from the external NIC API. No request body is required.

**Request Body**

None.

**Response** (200 OK)

```json
"Product created"
```

---

### GET /products

List all products stored in the database.

**Request**

No request body. Query parameters are not used.

**Response** (200 OK)

Array of product objects:

```json
[
  {
    "ID": "550e8400-e29b-41d4-a716-446655440000",
    "Name": "string",
    "ProductCode": "string",
    "ProductId": 0,
    "Description": "string",
    "CreatedAt": "2025-02-17T12:00:00Z"
  }
]
```

| Field | Type | Description |
|-------|------|-------------|
| ID | string (UUID) | Unique identifier |
| Name | string | Product name |
| ProductCode | string | Product code |
| ProductId | integer | External product ID |
| Description | string | Product description |
| CreatedAt | string (ISO8601) | Creation timestamp |

---

### Risk Type Endpoints

---

### POST /risk_type

Seed or create risk types from the external NIC API. No request body is required.

**Request Body**

None.

**Response** (200 OK)

```json
"Risk Type created"
```

---

### GET /risk_type

List all risk types stored in the database.

**Request**

No request body. Query parameters are not used.

**Response** (200 OK)

Array of risk type objects:

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "string",
    "risk_type_id": 0,
    "description": "string",
    "riskCategory": "string",
    "riskTypeCode": "string",
    "createdAt": "2025-02-17T12:00:00Z"
  }
]
```

| Field | Type | Description |
|-------|------|-------------|
| id | string (UUID) | Unique identifier |
| name | string | Risk type name |
| risk_type_id | integer | External risk type ID |
| description | string | Risk type description |
| riskCategory | string | Risk category |
| riskTypeCode | string | Risk type code |
| createdAt | string (ISO8601) | Creation timestamp |
