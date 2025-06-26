# Catalog Service

A gRPC-based catalog service that provides capabilities for catalog domain.

## Features

> [x] - implemented , [ ] - not implemented yet

- [x] Search product families 
- [ ] Fetch single family and options 
- [ ] Fetch products by criteria
- [ ] Fulltext search suggestions
- [ ] Fetch similar products by criteria 

Each feature is implemented as a separate package. See the respective package README for detailed feature documentation.

## Environment Variables
See .env.example
    
## Building and Running

### Build
```bash
make 
```

This will:
1. Clean the generated protobuf files
2. Generate new protobuf code
3. Build the service

### Run
```bash
bin/app grpc

2024-11-27T01:10:10+03:00 TRC Service has been started 
```

The server will start on `127.0.0.1:8009` by default.

## Project Structure

The project is based on a folder structure based on features ([Folder-by-feature](https://softwareengineering.stackexchange.com/questions/338597/folder-by-type-or-folder-by-feature)).
But that may change in the future.
 
### Root Level
```
.
├── api/          - Protocol buffer definitions and generated code
├── cmd/          - Application entry points
├── config/       - Configuration and dependency injection
├── internal/     - Internal packages
|   ├── adapter   - DEPRECATED External service adapters (gRPC, ES)
|   ├── domain    - DEPRECATED Core business logic
|   └── <feature_package> - package for the feature implementation
├── pkg/          - Shared packages
├── tests/        - DEPRECATED
└── README.md
```

### Feature Package structure

Each feature package follows a clean architecture pattern with these layers:

```
├── <feature_package>/
   ├── app/
   ├── config/
   ├── contr/
   ├── infra/
   └── tests/ 
```
#### `app/` - Application Layer
- Domain models with state: FiltersCollection, Filter
- Data Transfer Objects (DTOs)
- Commands and services implementing business logic
- Use cases and business rules

#### `config/` - Configuration Layer
- External dependency interfaces
- Functions to init the package 
- Environment configuration

#### `contr/` - Controller Layer
Controls I/O operations through various protocols:
- gRPC controllers
- HTTP controllers
- Console command controllers
- Defines the package's external interface

#### `infra/` - Infrastructure Layer
- Storage implementations
- External service clients: GRPC or HTTP clients to another servicec
- Protocol-independent communication code
- Database adapters

#### `tests/` - Integration Tests
- Integration test suites
- Test fixtures and helpers

## Development

### Adding a New Feature
1. Create a new feature package in `internal/`
2. Follow the feature package structure above
3. Implement required interfaces
4. Add integration tests
5. Update the feature checklist in this README

### Testing
```bash
# Run all tests
go test ./...

# Run specific feature tests
go test ./internal/<feature_package>/tests/...
```

### Logs 

Go to https://rke-dev.trgdev.com/dashboard/c/c-m-wwpn6hlr/explorer/apps.deployment/gws-develop/catalog#pods

Click by 3 dots in the pod list, select View logs.


