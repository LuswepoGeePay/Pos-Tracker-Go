# API Endpoints - GET vs POST Update Guide

## Overview of Changes

All read-only endpoints have been converted from POST requests with JSON bodies to proper GET requests with query string parameters. This follows REST conventions where:
- **GET** = Retrieve data (no side effects)
- **POST** = Create/modify data (side effects)

## Converted Endpoints

### Users

**Before (POST):**
```bash
curl -X POST http://localhost:8050/v1/users/get \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"page": 1, "pageSize": 10, "searchQuery": "john"}'
```

**After (GET):**
```bash
curl -X GET "http://localhost:8050/v1/users/get?page=1&pageSize=10&search=john" \
  -H "Authorization: Bearer <token>"
```

---

### Dashboard Events

**Before (POST):**
```bash
curl -X POST http://localhost:8050/v1/dashboard/events/get \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"page": 1, "pageSize": 10, "searchQuery": "login"}'
```

**After (GET):**
```bash
curl -X GET "http://localhost:8050/v1/dashboard/events/get?page=1&pageSize=10&search=login" \
  -H "Authorization: Bearer <token>"
```

---

### POS Devices

**Before (POST):**
```bash
curl -X POST http://localhost:8050/v1/pos/devices/get \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"page": 1, "pageSize": 10, "searchQuery": "kitchen"}'
```

**After (GET):**
```bash
curl -X GET "http://localhost:8050/v1/pos/devices/get?page=1&pageSize=10&search=kitchen" \
  -H "Authorization: Bearer <token>"
```

---

### Apps

**Before (POST):**
```bash
curl -X POST http://localhost:8050/v1/apps/get \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"page": 1, "pageSize": 10, "searchQuery": "pos"}'
```

**After (GET):**
```bash
curl -X GET "http://localhost:8050/v1/apps/get?page=1&pageSize=10&search=pos" \
  -H "Authorization: Bearer <token>"
```

---

### App Versions

**Before (POST):**
```bash
curl -X POST http://localhost:8050/v1/app/versions/get \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"page": 1, "pageSize": 10, "searchQuery": "v1.2"}'
```

**After (GET):**
```bash
curl -X GET "http://localhost:8050/v1/app/versions/get?page=1&pageSize=10&search=v1.2" \
  -H "Authorization: Bearer <token>"
```

---

### Locations

**Before (POST):**
```bash
curl -X POST http://localhost:8050/v1/locations/get \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"page": 1, "pageSize": 10, "searchQuery": "warehouse"}'
```

**After (GET):**
```bash
curl -X GET "http://localhost:8050/v1/locations/get?page=1&pageSize=10&search=warehouse" \
  -H "Authorization: Bearer <token>"
```

---

### Businesses

**Before (POST):**
```bash
curl -X POST http://localhost:8050/v1/businesses/get \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"page": 1, "pageSize": 10, "searchQuery": "acme"}'
```

**After (GET):**
```bash
curl -X GET "http://localhost:8050/v1/businesses/get?page=1&pageSize=10&search=acme" \
  -H "Authorization: Bearer <token>"
```

---

## Query Parameters Reference

All GET endpoints support the following query parameters:

| Parameter | Type | Default | Required | Description |
|-----------|------|---------|----------|-------------|
| `page` | integer | 1 | No | Page number for pagination |
| `pageSize` | integer | 10 | No | Number of items per page |
| `search` | string | "" | No | Search query to filter results |

### Examples:

**Get first page with default page size:**
```bash
GET /users/get
```

**Get second page with 20 items per page:**
```bash
GET /users/get?page=2&pageSize=20
```

**Search with default pagination:**
```bash
GET /users/get?search=john
```

**Combination of all parameters:**
```bash
GET /users/get?page=3&pageSize=25&search=john
```

---

## Endpoints Remaining as POST

The following endpoints correctly remain as POST requests (for creating or modifying data):

- `POST /v1/create-user` - User registration
- `POST /v1/login` - User login
- `POST /user/update` - Update user info
- `POST /v1/pos/register` - Register POS device
- `POST /pos/device/update` - Update device
- `POST /app/register` - Register app
- `POST /v1/app/update` - Check app update
- `POST /app/info/update` - Update app info
- `POST /app/version/register` - Register app version
- `POST /app/version/update` - Update app version
- `POST /v1/location/register` - Register location
- `POST /business/create` - Create business
- `POST /business/update` - Update business
- `POST /v1/pos/device/heartbeat` - Device heartbeat

---

## Code Examples

### JavaScript/Fetch API

**Old way (POST):**
```javascript
const response = await fetch('http://localhost:8050/v1/users/get', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    page: 1,
    pageSize: 10,
    searchQuery: 'john'
  })
});
```

**New way (GET):**
```javascript
const params = new URLSearchParams({
  page: 1,
  pageSize: 10,
  search: 'john'
});

const response = await fetch(`http://localhost:8050/v1/users/get?${params}`, {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

---

### React Query Example

```typescript
import { useQuery } from '@tanstack/react-query';

const useGetUsers = (page = 1, pageSize = 10, search = '') => {
  return useQuery({
    queryKey: ['users', page, pageSize, search],
    queryFn: async () => {
      const params = new URLSearchParams({ page, pageSize, search });
      const response = await fetch(`/api/users/get?${params}`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      return response.json();
    }
  });
};

// Usage
const { data, isLoading } = useGetUsers(1, 10, 'john');
```

---

### Axios Example

```typescript
import axios from 'axios';

const getUsers = async (page = 1, pageSize = 10, search = '') => {
  const response = await axios.get('/v1/users/get', {
    params: {
      page,
      pageSize,
      search
    },
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.data;
};

// Usage
getUsers(1, 10, 'john').then(users => console.log(users));
```

---

## Migration Checklist

If you're updating your frontend:

- [ ] Change all `GET /users/get` requests from POST to GET
- [ ] Change all `GET /pos/devices/get` requests from POST to GET
- [ ] Change all `GET /apps/get` requests from POST to GET
- [ ] Change all `GET /app/versions/get` requests from POST to GET
- [ ] Change all `GET /locations/get` requests from POST to GET
- [ ] Change all `GET /businesses/get` requests from POST to GET
- [ ] Change all `GET /dashboard/events/get` requests from POST to GET
- [ ] Update request headers (no need for `Content-Type: application/json`)
- [ ] Move request body to query string parameters
- [ ] Test all endpoints in staging environment
- [ ] Update API documentation
- [ ] Update client code/SDKs

---

## Benefits

1. **REST Compliance** - Follows REST conventions where GET is for retrieving data
2. **Simpler Testing** - Can test with browser URL bar or curl without body payload
3. **Caching** - GET requests can be cached by proxies and browsers
4. **Standard Behavior** - Consistent with HTTP specifications
5. **Better Semantics** - Clearer intent that the operation is read-only
6. **Reduced Bandwidth** - No need to send request body, query string is more concise

---

## Troubleshooting

### Still getting 400 errors?

**Problem:** Sending JSON body with GET request
**Solution:** Use query parameters instead (`?page=1&pageSize=10&search=term`)

**Problem:** Parameter values not being read
**Solution:** Ensure parameters are URL encoded if they contain special characters
```javascript
const params = new URLSearchParams({ search: 'john doe' });
// Automatically encodes as: search=john+doe
```

**Problem:** Case sensitivity issues
**Solution:** Use lowercase parameter names: `page`, `pageSize`, `search` (not `pageSize` or `SearchQuery`)

---

## API Documentation

For your API documentation, update all GET endpoint examples:

```yaml
GET /v1/users/get
Parameters:
  - name: page
    in: query
    type: integer
    default: 1
  - name: pageSize
    in: query
    type: integer
    default: 10
  - name: search
    in: query
    type: string
    required: false
```

---

**Last Updated:** March 2026
**Version:** 1.0
