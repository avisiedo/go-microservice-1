# /tasks Examples

This provide examples for the /tasks resource collection.

## POST /tasks

TODO Add example with curl

### Request

```json
{
    "list": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
    "description": "milk"
}
```

### Response

```json
{
    "uuid": "9349f788-c533-11ee-8b1b-482ae3863d30",
    "list": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
    "description": "milk",
}
```

## GET /tasks

TODO Add example with curl

### Request

### Response

```json
[
    {
        "uuid": "9349f788-c533-11ee-8b1b-482ae3863d30",
        "list": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
        "description": "milk",
    }
]
```

## GET /tasks/:id

TODO Add example with curl

### Request

### Response

```json
{
    "uuid": "9349f788-c533-11ee-8b1b-482ae3863d30",
    "list": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
    "description": "milk",
}
```

## PATCH /tasks/:id

TODO Add example with curl

Partial update; look at as every field is optional.

### Request

```json
{
    "list": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
    "description": "milk",
}
```

### Response

```json
{
    "uuid": "9349f788-c533-11ee-8b1b-482ae3863d30",
    "list": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
    "description": "milk",
}
```

## PUT /tasks/:id

TODO Add example with curl

Completly replace a resource (same as POST but for an existing resource).

### Request

```json
{
    "list": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
    "description": "milk",
}
```

### Response

```json
{
    "uuid": "9349f788-c533-11ee-8b1b-482ae3863d30",
    "list": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
    "description": "milk",
}
```

## DELETE /tasks/:id

TODO Add example with curl

## Request

## Response

```raw
204 (No Content)
```
