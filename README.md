# nuage

nuage (French for cloud) is a cloud-native, API-first framework for Go designed
to build robust, scalable, and standards-compliant services.

nuage is built on official RFCs and industry best practices, ensuring
consistency, interoperability, and long-term maintainability. It focuses on
developer productivity while embracing modern cloud architectures from the
ground up.

## Philisophy

nuage is a opinioated framework trying to force the developer to implement clean
and robust APIs and make it hard to develop bad patterns. For that it is based
on the following principles.

### JSON Only

To make APIs compatible to various other tools it is best practice to
communicate using [JSON](https://datatracker.ietf.org/doc/html/rfc8259). Other
formats such as XML are not supported. Furhter using JSON allows nuage to
leverage JSON Schema and OpenAPI standards for documentation and validation.

### Configuration

Because nuage is a cloud-native first framework and implements matching factors
of the [12 Factor App](https://12factor.net/). This includes allowing to
configure nuage using only environment variables.

### Logging

Because of the 12FA principles logging will be only possible to the stdout and
in structured format.

### HTTP Errors

Errors are part of every software and have to be first-class citizens in the
design. RFC9457 is describing an official standard for returning HTTP errors
which is the only error response from nuage.

### OpenAPI

OpenAPI is a standard to document APIs and make them compatible to many tools to
generate REST clients, CLIs etc. This is making it a central design document for
many developers. Therefore nuage is trying to generate as much of the OpenAPI
documentation from your code with the possibility of extending it for further
customization.

### Compile time over Runtime

Reflection is a nice tool in Go allowing for powerful analysis of types but the
operations are expensive and create large overheads. Therefore any reflection in
the hot-path (e.g. requests) are outsourced to compile time by generating the
required code for the runtime beforehand.
