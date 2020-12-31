# Templates

A collection of templates in miscellaneous programming languages and technologies for quickly starting a new project!

## IDL

_Interface Description Language_ or _Interface Definition Language_ is a formal language for defining interfaces between components.
It also allows you to define types that can be used for input and output as well as validations.
By adopting an IDL (_gRPC_, _Thrift_, etc.) you can specify well-defined and strongly-typed interfaces for your services.

## Slicing Your Domain

When it comes to organizing your service code into different packages or modules, there are generally two approaches.
You can either slice your domain horizontally or vertically.

### Horizontally

This approach is also known as _layered architecture_, _onion architecture_, etc.
In this approach, you will slice your service into horizontal layers (packages or modules).
Each layer is focused on handling one aspect of a request.
For a request to be fulfilled, it should go through all layers from top to bottom (or outside to inside).
Layers are abstracted away from each other using interfaces.
Each layer only depends on the layer below and the flow between layers is one-way from top to bottom (or outside to inside).

Typically, we have the following layers:

| Layer      | Description                                                                                   |
|------------|-----------------------------------------------------------------------------------------------|
| Entity     | Your domain-specific representation of request and response objects.                          |
| Mapper     | Mapper functions for mapping between domain-specific and protocol-specific models.            |
| Gateway    | Providing access to external services.                                                        |
| Repository | Providing access to external data stores.                                                     |
| Controller | Implementing the core domain functionality and logic agnostic of any transport protocol.      |
| Handler    | Translating protocol-specific requests to domain-specific ones using mappers and controllers. |

Controllers can follow a _composable architecture_.
You can additionally have a separate package or module for modeling requests and responses to external services used by gateways.

### Vertically

Another approach is to slice your service vertically.
In this approach, you will break your service into a disjoint set of packages or modules formed around your domain functionalities.
You will have a package or module per each domain functionality. These packages together partition your request space.
For a request to be fulfilled, it should go only through one package or module and all aspects of handling the request will be carried out by that package or module.

### Comparison

Like any other layered architecture, horizontal slicing gives a lot of flexibility for changing a layer or completely replacing one.
You can easily replace your handler layer to switch to gRPC protocol from HTTP without changing your entities and controllers.
Also since your controller only has your domain logic and everything else is moved out, your domain logic is isolated and thus easier to follow.
However, like all other layered architectures, this comes at the cost of more overhead, more time-consuming development, harder debuggability, and a lot of boilerplate codes.

This vertical slicing has a better and more clear semantic comparing to the horizontal approach.
For following a request through your service, you only need to look at one function in one package or module.
The development experience is better in terms of velocity, testing, debugging, and monitoring.

In general, although horizontal slicing promises better flexibility and abstraction, in reality, we rarely leverage and benefit from that flexibility.
How many times did have to switch the transport protocol in your last job?
Vertical slicing works best when it is used alongside an IDL and a set of solid and robust tools.
The models defined in IDL can be directly used as universal models for our domain and we will not need to have entities and mappers anymore.
The IDL tools can take care of tasks not related to domain logic such as validation, serializing, deserializing, etc.
This way we can eliminate boilerplate codes and keep our vertical packages extremely lean and slick.
As a developer, you now only need to implement your domain logic in one place, and switching to a different protocol also becomes a very low-cost task.

## Control Flow

When it comes to composing and building an application, we have two options.
We can let developers import and define everything they need explicitly or we can ask them to only _fill in the blanks_ (inverted).
We use the terms _explicit_ and _inverted_ loosely here.

### Explicit

In an explicit control flow, you will have full control over the final application and how it is built.
You will create all the required files and structures and implement the required definitions and procedures for your application.
You will import built-in or third-party libraries for brevity, reusability, modularity, consistency, security, etc.
You will decide how your application proceeds from start to end.

Needless to say, you can leverage available tools to increase your velocity and productivity.
For example, you can use a scaffolding tool to generate a template application in your programming language of choice.

### Inverted

This approach is based on the [inversion of control](https://en.wikipedia.org/wiki/Inversion_of_control)
and [dependency inversion](https://en.wikipedia.org/wiki/Dependency_inversion_principle) principles.

In an inverted control flow, you will only fill in the blanks provided by a framework.
You will have to provide implementations for some well-defined abstractions (interfaces).
Dependencies (logger, database access, etc.) are also provided to you as interfaces (so you can easily mock them).
Using a _code generation_ or a _runtime_ framework, your implementations alongside dependency implementations will be injected into a generated or predefined application.
In other words, the framework will provide the dependencies and call into your implementations.

Using this approach, non-functional qualities such as security, observability, resiliency, consistency, etc. automatically can be ensured.

## Similar Projects

### Go Kit

Go kit is a toolkit for microservices. It takes a similar approach to slicing domain and organizing a service.
Go kit introduces similar concepts such as requests and responses, endpoints, transports, middleware, etc.
It also provides standard libraries for logging, metrics, tracing, authentication, and so forth.

  - https://gokit.io
  - https://github.com/go-kit/kit
  - https://pkg.go.dev/github.com/go-kit/kit

## Reading More

  - **OpenAPI**
    - [OpenAPI Specification](https://swagger.io/specification)
    - [OpenAPI 3.0 Tutorial](https://app.swaggerhub.com/help/tutorials/openapi-3-tutorial)
  - **Protocol Buffers & gRPC**
    - [Language Guide (proto3)](https://developers.google.com/protocol-buffers/docs/proto3)
    - [Introduction to gRPC](https://grpc.io/docs/what-is-grpc/introduction)
    - [Core Concepts, Architecture and Lifecycle](https://grpc.io/docs/what-is-grpc/core-concepts)
    - [Go Generated Code Reference](https://grpc.io/docs/languages/go/generated-code)
    - [Package grpc](https://pkg.go.dev/google.golang.org/grpc)
  - **Architecture**
    - [Domain-Driven Design: Tackling Complexity in the Heart of Software](https://www.amazon.com/Domain-Driven-Design-Tackling-Complexity-Software/dp/0321125215)
    - [Patterns of Enterprise Application Architecture](https://www.amazon.com/Patterns-Enterprise-Application-Architecture-Martin/dp/0321127420)
    - [The Onion Architecture](https://jeffreypalermo.com/2008/07/the-onion-architecture-part-1)
    - [Catalog of Patterns of Enterprise Application Architecture](https://martinfowler.com/eaaCatalog)
  - **Go**
    - [Protocol Buffer Basics: Go](https://developers.google.com/protocol-buffers/docs/gotutorial)
  - **JavaScript**
    - [Classes](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Classes)
    - [ECMAScript Modules](https://nodejs.org/api/esm.html)
  - **Misc**
    - [Using submodules in Git](https://www.vogella.com/tutorials/GitSubmodules/article.html)
    - [GitHub Actions](https://docs.github.com/en/actions)
