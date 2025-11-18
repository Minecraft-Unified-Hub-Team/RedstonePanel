# RedstonePanel Architecture

**RedstonePanel** is built upon a **hierarchical, modular architecture** with clear dependency management and separation of concerns. This design was chosen to achieve the following key goals:

-   **Independent Module Development:** Each module can be developed, tested, and debugged independently from others.
-   **High Testability:** Every module can be unit tested in isolation with clear interfaces and dependencies.
-   **Clear Dependency Hierarchy:** Dependencies flow in one direction, preventing circular references and making the system easier to understand.
-   **Maintainability:** Well-defined module boundaries make it easier to modify, replace, or extend functionality.

## Architecture Overview

The system follows a layered hierarchical architecture with three main tiers:

### Presentation Layer
- **WebUI**: Web-based user interface
- **API**: RESTful API for all operations and external integrations
- **GUI**: Desktop graphical user interface (optional)

### Core Logic Layer
- **MineDB**: Centralized JSON-based data storage with thread-safety and durability
- **InstanceInstallation**: Minecraft server installation and Java management
- **InstanceManifest**: Server configuration management (server.properties, JVM args)
- **InstanceControl**: Server lifecycle management (start/stop/status)
- **InstanceMods**: Mod installation and management
- **InstancePlugins**: Plugin installation and management
- **InstanceModpacks**: Modpack installation and management

### External Services Layer
- **MineDL**: Minecraft version downloads
- **HttpDownloader**: HTTP-based file downloading
- **Scheduler/TelegramBot**: External integrations and scheduling

## Dependency Rules

The architecture enforces strict dependency rules:

1. **Upward Dependencies Only**: Lower layers can depend on upper layers, but never the reverse
2. **No Horizontal Dependencies**: Modules within the same layer cannot directly depend on each other
3. **Shared Storage**: All modules use MineDB for persistent state management
4. **API-Mediated Access**: All external access flows through the RESTful API layer
5. **Interface-Based Communication**: All inter-module communication happens through well-defined interfaces

## Key Design Principles

### 1. Single Responsibility
Each module handles one specific domain:
- `MineDB` handles all data persistence
- `InstanceInstallation` handles server installation
- `InstanceControl` handles server lifecycle
- `InstanceManifest` handles configuration files

### 2. Dependency Injection
Modules receive their dependencies through constructors, making them easily testable and configurable.

### 3. Thread Safety
All modules are designed to be thread-safe, particularly `MineDB` which handles concurrent access from multiple modules.

### 4. Error Handling
Each module implements comprehensive error handling with clear error propagation up the dependency chain.

### 5. API-First Design
The RESTful API layer serves as the primary interface for all operations, providing:
- **Consistent Interface**: Uniform HTTP-based access to all functionality
- **External Integration**: Easy integration with external tools and services
- **Web Interface Support**: Direct support for web-based user interfaces
- **Authentication & Authorization**: Centralized security controls

## Module Documentation

For detailed information about each module, see the individual module documentation:

- [MineDB](developer/architecture/modules/minedb.md) - Centralized data storage
- [InstanceInstallation](developer/architecture/modules/instance-installation.md) - Server installation
- [InstanceManifest](developer/architecture/modules/instance-manifest.md) - Configuration management
- [InstanceControl](developer/architecture/modules/instance-control.md) - Server lifecycle
- [InstanceMods](developer/architecture/modules/instance-mods.md) - Mod management

## Testing Strategy

The hierarchical architecture enables comprehensive testing at multiple levels:

1. **Unit Tests**: Each module tested in isolation with mocked dependencies
2. **Integration Tests**: Testing module interactions through well-defined interfaces  
3. **System Tests**: End-to-end testing of complete user workflows
4. **Contract Tests**: Ensuring interface compatibility between modules

