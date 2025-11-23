# nuage

nuage is a minimal Go REST API framework based on best practices and official
standards. It is inspired by huma (currently not mantained) and FastAPI with a
more explicit and standard library conform approach.

# Conforming to many standards

nuage is built from the ground up with the idea of conforming to as many
standards and best practices of the web and industry to make your API as
compatible as possible. The following Standards are implemented:

- RFC 9110 HTTP Semantics
- RFC 9111 HTTP Caching
- RFC 9112 HTTP/1.1
- RFC 9114 HTTP/3 (quic)
- RFC 3986 URI Syntax
- RFC 8288 Web Linking
- RFC 9114 .12 Content Negotiation
- RFC 8259 JSON Data Format
- JSON Schema (2020-12)
- RFC 6902 JSON Patch
- RFC 7386 JSON Merge PATCH
- RFC 9457 HTTP Problem Details
- JSON:API Spec
- RFC 9421 HTTP Message Signatures
- RFC 9110 Conditional Requests
- RFC 9333 RateLimit Headers
- OpenAPI Spec
- AsyncAPI Spec

## TODOs

Custom format support
Require /livez and /readyz endpoints to make it k8s compatibale for liveness and readiness probe
