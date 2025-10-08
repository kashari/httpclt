Go HTTP CLI Utility
A simple, configurable command-line utility written in Go to make HTTP requests. It's a lightweight alternative to tools like curl for basic API testing and load testing directly from your terminal.

Features
Specify the request URL, method, body, and headers via flags.

Supports sending multiple headers by repeating the --header flag.

Run multiple requests concurrently for load testing.

Rate limit requests on a per-second basis.

Prints a clear summary of both the outgoing request (for single requests) and the incoming response status.

Cross-platform (compiles for Windows, macOS, and Linux).

Installation
To use this utility, you need to have Go installed on your system.

./httpclt --help

Command-Line Flags
Flag

Description

Default

--url

The URL for the HTTP request. (Required)

""

--method

The HTTP method to use (e.g., GET, POST).

"GET"

--body

The request body for POST, PUT, etc.

""

--header

A request header in 'Key: Value' format. Can be specified multiple times.

[]

--requests

The total number of requests to send.

1

--per-second

The maximum number of requests per second. 0 means no limit.

0

Examples
1. Simple GET Request

When making a single request, the tool provides a verbose output showing the request and response details.

./httpclt --url https://jsonplaceholder.typicode.com/todos/1

2. POST Request with JSON Body and Headers

To send data, use the --method, --header, and --body flags. Make sure to wrap the JSON body in single quotes to prevent shell interpretation.

./httpclt \
  --url https://jsonplaceholder.typicode.com/posts \
  --method POST \
  --header "Content-Type: application/json; charset=UTF-8" \
  --body '{"title": "foo", "body": "bar", "userId": 1}'

3. Concurrent GET Requests (Load Testing)

Send 50 requests at a rate of 10 requests per second. The output for each request is condensed to show just the status.

./httpclt \
  --url https://jsonplaceholder.typicode.com/todos/1 \
  --requests 50 \
  --per-second 10

4. Concurrent POST Requests without Rate Limiting

Send 100 requests as fast as possible.

./httpclt \
  --url https://jsonplaceholder.typicode.com/posts \
  --method POST \
  --header "Content-Type: application/json" \
  --body '{"title": "concurrent post"}' \
  --requests 100

5. DELETE Request

This request does not typically require a body.

./httpclt --url https://jsonplaceholder.typicode.com/posts/1 --method DELETE

