# layer7lb
layer 7 loadbalancer in go

## Background

An application-layer loadbalancer, typical using HTTP protocol.

- HTTP request comes in.
- Loadbalancer passes request onto a number of backends.

Benefits (over Layer 4):

- Look into HTTP protocol more.
- Can look at a path of URL and choose different types of backends for routing,
access control, etc.

Disadvantages:

- More resources to process all the headers all the way up to layer 7.
- More specific to one application.
- With TCP proxy, don't care about what protocol goes over it.

May have a keep-alive (layer 4 also has these):
- different way to keep the connection open from front-end to lb or from lb to back-end.
- one strategy - send a small packet
- establishing a connection can be time-consuming so efficient.
- often golang client library handles this.
- often can send multiple reuqests on same connection but connections aren't used
concurrently.

## Requirements

clients ---- (send request)/(sends response) <-----> lb --(forwards request)----> backends

- Statically configure list of back-ends.
- 1 Front-end - perhaps different domain names, for example.
- Many Back-ends - different server/host options.
- Load-balancing strategies - for every request, how to choose which backend?
- Translate errors from backend back to client.
- Max number of concurrent requests.

## Pseudocode:

1. Configuration - via flags or file.
2. Have a statically defined list of back-ends.
3. Loadbalancer itself is an HTTP server listening for requests.
4. Global http client is created that leverages keep-alives.
4. Need a handler implementing the serve.HTTP method:
  - HTTP framework in go will spawn off new goroutine and call handler method.
  - a backend is selected through some strategy (maybe randomly).
  - global http client is used to send request to chosen backend.
  - Response is returned by handler.
  - If a backend returns an error, handler returns status code and explicit error message.
      -> 4xx error (client caused) - transparently pass that back
      -> 5xx - (server caused) - return bad gateway or try with another backend
      -> 1xx/2xx/3xx - transparently pass that back
      -> If there's a go error, return bad gateway response to the client. 
