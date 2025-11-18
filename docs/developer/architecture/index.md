# Overview

This section contains comprehensive documentation about RedstonePanel's architecture, requirements, and design decisions.

## Contents

- [Architecture](developer/architecture/architecture.md) - Application architecture documentation
- [Requirements](developer/architecture/requirements.md) - Functional and non-functional requirements
- [API Layer](developer/architecture/api.md) - RESTful API interface documentation

## Module Documentation

RedstonePanel follows a hierarchical modular architecture. Each module is documented individually:

- [MineDB](developer/architecture/modules/minedb.md) - Centralized JSON-based data storage
- [InstanceInstallation](developer/architecture/modules/instance-installation.md) - Minecraft server installation and Java management
- [InstanceManifest](developer/architecture/modules/instance-manifest.md) - Server configuration file management
- [InstanceControl](developer/architecture/modules/instance-control.md) - Server lifecycle management (start/stop/status)
- [InstanceMods](developer/architecture/modules/instance-mods.md) - Mod installation and management system

## Quick Start

For developers new to the RedstonePanel architecture:

1. **Start with [Requirements](developer/architecture/requirements.md)** to understand what we're building
2. **Read [Architecture](developer/architecture/architecture.md)** to understand the overall system design
3. **Explore individual module documentation** to understand specific components
4. **Review [API Layer](developer/architecture/api.md)** to understand the RESTful interface layer

## Design Principles

Our architecture is built around these key principles:

- **Hierarchical Dependencies**: Clear dependency flow prevents circular references
- **Independent Modules**: Each module can be developed and tested in isolation
- **Centralized State**: All persistent data flows through MineDB
- **Interface-based Communication**: Well-defined contracts between modules
