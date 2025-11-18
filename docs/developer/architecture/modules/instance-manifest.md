# InstanceManifest Module

All code that provided in that docs is just an example. For real code check sources.

## Overview

The InstanceManifest module manages server configuration files such as `server.properties`, `jvm.args`, and other configuration files. It provides a unified interface for reading and writing server configurations while maintaining state synchronization with the file system.

## Purpose and Scope

This module provides:

- **Configuration File Management**: Read/write server configuration files
- **Format Abstraction**: Handle different configuration file formats uniformly
- **State Synchronization**: Keep MineDB state in sync with file system
- **Validation**: Ensure configuration values are valid and compatible

## Architecture

### Location
```
internal/instance_manifest/
```

### Dependencies
- **MineDB**: For configuration state storage and caching
- **File System**: For reading/writing configuration files
- **Validation**: Configuration value validation

## Configuration

### ManifestConfig Structure
```go
type ManifestConfig struct {
    // Base path where server instances are located
    StoreInstancePath string
}
```

### ServerManifestType
```go
// Defined in internal/common/
type ServerManifestType string

const (
    ServerProperties ServerManifestType = "server.properties"
    JvmArgs          ServerManifestType = "jvm.args"
    EulaText         ServerManifestType = "eula.txt"
    Whitelist        ServerManifestType = "whitelist.json"
    BannedPlayers    ServerManifestType = "banned-players.json"
    BannedIPs        ServerManifestType = "banned-ips.json"
    Ops              ServerManifestType = "ops.json"
)
```

## API Interface

### Constructor
```go
func NewInstanceManifestModule(db MineDB, config ManifestConfig) (*InstanceManifest, error)
```

### Core Methods

```go
// Write configuration manifest to file system
WriteManifest(targetServerPath string, mtype ServerManifestType, manifest map[string]string) error

// Read configuration manifest from file system
ReadManifest(targetServerPath string, mtype ServerManifestType) (map[string]string, error)
```

## Data Storage

### Storage Key Format
```
instancemanifest:<targetInstancePath>
```

### Storage Schema
```go
type InstanceManifestData struct {
    // Server path identifier
    ServerPath string `json:"server_path"`
    
    // Configuration files and their content
    Manifests map[ServerManifestType]ManifestContent `json:"manifests"`
    
    // Last modification timestamps
    LastModified map[ServerManifestType]time.Time `json:"last_modified"`
    
    // Configuration validation status
    ValidationStatus map[ServerManifestType]ValidationResult `json:"validation_status"`
}

type ManifestContent struct {
    // Configuration key-value pairs
    Properties map[string]string `json:"properties"`
    
    // File format metadata
    Format FileFormat `json:"format"`
    
    // Checksum for integrity verification
    Checksum string `json:"checksum"`
}

type ValidationResult struct {
    IsValid bool     `json:"is_valid"`
    Errors  []string `json:"errors,omitempty"`
}

type FileFormat string

const (
    PropertiesFormat FileFormat = "properties"  // key=value format
    JsonFormat       FileFormat = "json"        // JSON format
    TextFormat       FileFormat = "text"        // Plain text
)
```

## Implementation Details

### File Format Handlers

#### Properties Format (server.properties)
```properties
# Minecraft Server Properties
server-port=25565
gamemode=survival
difficulty=normal
max-players=20
```

#### JVM Arguments Format
```
-Xmx2G
-Xms1G
-XX:+UseG1GC
-XX:+ParallelRefProcEnabled
```

#### JSON Format (whitelist.json, ops.json)
```json
[
  {
    "uuid": "f430dbb6-0a58-4c64-9d79-7e97d69c7e29",
    "name": "PlayerName"
  }
]
```

### Configuration Validation

#### Server Properties Validation
```go
type ServerPropertiesValidator struct {
    // Port range validation
    ValidPortRange struct {
        Min int
        Max int
    }
    
    // Valid game modes
    ValidGameModes []string
    
    // Valid difficulties  
    ValidDifficulties []string
}

func (v *ServerPropertiesValidator) Validate(properties map[string]string) ValidationResult
```

#### JVM Arguments Validation
```go
type JvmArgsValidator struct {
    // Memory allocation limits
    MaxMemory string
    MinMemory string
    
    // Allowed JVM flags
    AllowedFlags []string
}
```

## Usage Examples

### Writing Server Properties
```go
manifest := map[string]string{
    "server-port":    "25565",
    "gamemode":       "survival", 
    "difficulty":     "normal",
    "max-players":    "20",
    "online-mode":    "true",
    "enable-rcon":    "false",
}

err := manifestModule.WriteManifest(
    "/opt/minecraft/instances/server1",
    ServerProperties,
    manifest,
)
if err != nil {
    log.Printf("Failed to write server properties: %v", err)
}
```

### Reading Server Properties
```go
manifest, err := manifestModule.ReadManifest(
    "/opt/minecraft/instances/server1",
    ServerProperties,
)
if err != nil {
    log.Printf("Failed to read server properties: %v", err)
    return
}

port := manifest["server-port"]
gamemode := manifest["gamemode"]
fmt.Printf("Server running on port %s in %s mode\n", port, gamemode)
```

### Managing JVM Arguments
```go
jvmArgs := map[string]string{
    "0": "-Xmx2G",
    "1": "-Xms1G", 
    "2": "-XX:+UseG1GC",
    "3": "-XX:+ParallelRefProcEnabled",
}

err := manifestModule.WriteManifest(
    "/opt/minecraft/instances/server1",
    JvmArgs,
    jvmArgs,
)
```

## File Operations

### File Path Resolution
```go
func (m *InstanceManifest) getConfigFilePath(serverPath string, mtype ServerManifestType) string {
    switch mtype {
    case ServerProperties:
        return filepath.Join(serverPath, "server.properties")
    case JvmArgs:
        return filepath.Join(serverPath, "jvm.args")
    case EulaText:
        return filepath.Join(serverPath, "eula.txt")
    // ... other types
    }
}
```

### Atomic File Operations
```go
func (m *InstanceManifest) writeConfigFile(filePath string, content string) error {
    // Write to temporary file first
    tmpFile := filePath + ".tmp"
    
    if err := ioutil.WriteFile(tmpFile, []byte(content), 0644); err != nil {
        return err
    }
    
    // Atomic rename
    return os.Rename(tmpFile, filePath)
}
```

## Error Handling

### Error Types
- **FileNotFoundError**: Configuration file doesn't exist
- **ValidationError**: Invalid configuration values
- **FormatError**: Malformed configuration file
- **PermissionError**: Insufficient file system permissions
- **StateError**: MineDB synchronization failures

### Error Recovery
- **Default Configuration**: Generate default configs for missing files
- **Backup and Restore**: Create backups before modifications
- **Validation Rollback**: Revert invalid configurations

## Testing Strategy

### Unit Tests
- File format parsing and generation
- Configuration validation logic
- Error handling scenarios
- State synchronization

### Integration Tests
- End-to-end configuration management
- File system operations
- MineDB integration
- Cross-platform compatibility

### Test Data
```go
var testServerProperties = map[string]string{
    "server-port":     "25565",
    "gamemode":        "survival",
    "difficulty":      "normal",
    "max-players":     "20",
    "online-mode":     "true",
    "spawn-protection": "16",
}

var invalidServerProperties = map[string]string{
    "server-port": "99999",  // Invalid port
    "gamemode":    "invalid", // Invalid gamemode
}
```

## Security Considerations

### Input Validation
- **Value Sanitization**: Sanitize all configuration values
- **Type Validation**: Ensure values match expected types
- **Range Validation**: Validate numeric ranges

### File System Security
- **Path Validation**: Prevent directory traversal
- **Permission Management**: Set appropriate file permissions
- **Backup Security**: Secure storage of configuration backups

## Future Enhancements

### Optimization Strategies
- **Caching**: Cache frequently accessed configurations
- **Batch Operations**: Support bulk configuration updates
- **Lazy Loading**: Load configurations only when needed

### Resource Management
- **Memory Usage**: Efficient handling of large configuration files
- **File Locking**: Minimize file lock duration
- **Concurrent Access**: Support concurrent read operations

### Other Planned Features
- **Configuration Templates**: Predefined configuration templates
- **History Tracking**: Track configuration change history
- **Hot Reloading**: Apply configuration changes without restart
- **Validation Plugins**: Extensible validation system

### Advanced Features
- **Configuration Diffing**: Show configuration differences
- **Bulk Operations**: Manage configurations across multiple instances
- **Import/Export**: Configuration import/export functionality
- **Automated Optimization**: AI-driven configuration optimization
