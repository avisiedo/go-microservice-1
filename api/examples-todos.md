# /todos Examples

This provide examples for the /todos resource collection.

## POST /todos

TODO Add example with curl

### Request

```json
{
    "title": "Shopping list",
    "description": "My shopping list for the weekend",
    "tasks":[
        {
            "description": "milk"
        }
    ]
}
```

### Response

```json
{
    "uuid": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
    "title": "shopping list",
    "description": "My shopping list for the weekend",
    "tasks":[
        {
            "uuid": "9349f788-c533-11ee-8b1b-482ae3863d30",
            "description": "milk"
        }
    ]
}
```

## GET /todos

- TODO Add examples with curl.
- TODO Add pagination and a default limit of page when not indicated.

### Request

### Response

```json
[
    {
        "uuid": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
        "title": "shopping list",
        "description": "My shopping list for the weekend",
        "tasks":[
            {
                "uuid": "9349f788-c533-11ee-8b1b-482ae3863d30",
                "description": "milk"
            }
        ]
    }
]
```

## GET /todos/:id

TODO Add example with curl

### Request

### Response

```json
{
    "uuid": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
    "title": "shopping list",
    "description": "My shopping list for the weekend",
    "tasks":[
        {
            "uuid": "9349f788-c533-11ee-8b1b-482ae3863d30",
            "description": "milk"
        }
    ]
}
```

## PATCH /todos/:id

TODO Add example with curl

Partial update; look at as every field is optional.

### Request

```json
{
    "title": "shopping list",
    "description": "My shopping list for the weekend",
    "tasks":[
        {
            "description": "milk"
        }
    ]
}
```

### Response

```json
{
    "uuid": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
    "title": "shopping list",
    "description": "My shopping list for the weekend",
    "tasks":[
        {
            "uuid": "9349f788-c533-11ee-8b1b-482ae3863d30",
            "description": "milk"
        }
    ]
}
```

## PUT /todos/:id

TODO Add example with curl

Completly replace a resource (same as POST but for an existing resource).

### Request

```json
{
    "title": "Shopping list",
    "description": "My shopping list for the weekend",
    "tasks":[
        {
            "description": "milk"
        }
    ]
}
```

### Response

```json
{
    "uuid": "979c1f40-c4e9-11ee-8ba1-482ae3863d30",
    "title": "shopping list",
    "description": "My shopping list for the weekend",
    "tasks":[
        {
            "uuid": "9349f788-c533-11ee-8b1b-482ae3863d30",
            "description": "milk"
        }
    ]
}
```

## DELETE /todos/:id

TODO Add example with curl

## Request

## Response

```raw
204 (No Content)
```
