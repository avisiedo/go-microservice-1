# Examples

Here are added examples of the API which will be helpful
to bootstrap the micro-service by extracting schema from
the examples in apicurio.

## General

- Content-Type: application/json
- Location: /{todos,tasks}/:id

When returning record timestamp:

- Last-Modified: [Last-Modified][rfc9110_8_8_2]

Status Codes on Success

- POST /path: HTTP 201 Created
- PUT /path: HTTP 200 Ok | HTTP 201 Created
- PATCH /path: HTTP 200 Ok
- DELETE /path: HTTP 204 No content

Status Codes on Error

- 404 Not found

Final tips:

- The more accurate is your openapi, the better to automate
  the validation.
- Security is not managed on this repository; it use to be
  implemented as a middleware or in a more abstract way
  wrapping the artifact deployed in the cloud.
- Design the API focus on the resources.
- Avoid nested resources in the paths if possible.
- Use plural form to name the resource on the endpoint as
  it is referenced the resource collection; when a specific
  resource is requested is the id provided after the resource.

## Resources

- [/todos][examples_todos]
- [/tasks][examples_tasks] belongs to one /todos resource.

[rfc9110_8_8_2]: https://www.rfc-editor.org/rfc/rfc9110#section-8.8.2
[examples_todos]: examples-todos.md
[examples_tasks]: examples-tasks.md