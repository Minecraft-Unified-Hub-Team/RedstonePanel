# MineDB Module

All code that provided in that docs is just an example. For real code check sources.

## Overview

MineDB is the centralized data storage module for RedstonePanel. It provides a unified, thread-safe, and durable JSON-based storage system that all other modules use for persistent state management.

## Purpose and Motivation

Instead of each module maintaining separate state files, MineDB provides:

- **Centralized Storage**: Single point of data management
- **Thread Safety**: Concurrent access from multiple modules
- **Durability**: Guaranteed persistence across application restarts
- **Isolation**: Namespaced data to prevent module interference

## Architecture

### Location
```
internal/mine_db/
```

### Key Features

- **Thread-Safe Operations**: Concurrent read/write operations without data corruption
- **Durability**: All committed changes survive application crashes
- **Isolation**: Module data is stored in separate namespaces using key prefixes
- **JSON Format**: Human-readable storage format for debugging and maintenance

## API Interface

### Core Operations

```go
type MineDB interface {
    Get(key string) (any, error)
    Set(key string, value any) error
}
```

#### Get Operation
```go
Get(key string) (any, error)
```
- **Purpose**: Retrieve value by key
- **Returns**: Value if found, error if key doesn't exist
- **Thread Safety**: Safe for concurrent reads

#### Set Operation  
```go
Set(key string, value any) error
```
- **Purpose**: Store or update value by key
- **Returns**: Error if operation fails
- **Thread Safety**: Safe for concurrent writes
- **Durability**: Changes are immediately persisted

## Key Naming Convention

Each module uses a specific key prefix to ensure data isolation:

| Module | Key Prefix | Example |
|--------|------------|---------|
| InstanceInstallation | `instanceinstallation:` | `instanceinstallation:/path/to/instance` |
| InstanceManifest | `instancemanifest:` | `instancemanifest:/path/to/instance` |
| InstanceControl | `instancecontrol:` | `instancecontrol:/path/to/instance` |
| InstanceMods | `instancemods:` | `instancemods:/path/to/instance` |

## Implementation Requirements

### Thread Safety
- Use appropriate synchronization primitives (mutexes, channels)
- Ensure atomic read/write operations
- Handle concurrent access without blocking unnecessary operations

### Durability
- Implement write-ahead logging or similar mechanism
- Ensure data consistency during concurrent writes
- Graceful handling of disk I/O errors

### Error Handling
- Clear error messages for debugging
- Proper handling of file system errors
- Validation of input data types

## Usage Examples

### Storing Instance Installation Data
```go
type InstanceData struct {
    Version    string `json:"version"`
    ServerType string `json:"server_type"`
    JavaPath   string `json:"java_path"`
}

data := InstanceData{
    Version:    "1.19.2",
    ServerType: "vanilla",
    JavaPath:   "/usr/bin/java",
}

err := db.Set("instanceinstallation:/instances/server1", data)
if err != nil {
    return fmt.Errorf("failed to save instance data: %w", err)
}
```

### Retrieving Instance Data
```go
rawData, err := db.Get("instanceinstallation:/instances/server1")
if err != nil {
    return fmt.Errorf("instance not found: %w", err)
}

data, ok := rawData.(InstanceData)
if !ok {
    return fmt.Errorf("invalid data type")
}
```

## Testing Strategy

### Unit Tests
- Test concurrent read/write operations
- Verify data persistence across restarts
- Test error conditions (invalid keys, disk full, etc.)
- Validate key isolation between modules

### Integration Tests
- Test with actual file system operations
- Verify behavior under high concurrency
- Test recovery from corrupted data files

## Dependencies

### Internal Dependencies
- Standard Go libraries for JSON handling
- File system operations
- Synchronization primitives

### External Dependencies
- None (pure Go implementation)

## Security Considerations

- **Input Validation**: Validate all keys and values before storage
- **Path Safety**: Ensure key names don't allow path traversal
- **Data Sanitization**: Sanitize JSON data to prevent injection attacks

## Performance Considerations

- **Caching**: Consider in-memory caching for frequently accessed data
- **Batch Operations**: Optimize for bulk read/write scenarios
- **File Locking**: Minimize lock contention for better concurrent performance

## Future Enhancements

- **Backup/Restore**: Automated backup functionality
- **Data Migration**: Support for schema migrations
- **Compression**: Optional data compression for large datasets
- **Encryption**: Optional data encryption at rest
