# InstanceInstallation Module

All code that provided in that docs is just an example. For real code check sources.

## Overview

The InstanceInstallation module handles the installation of "clean" Minecraft servers without mods, plugins, or shader packs. It manages Java runtime installation, Minecraft server downloads, and maintains installation state.

## Purpose and Scope

This module provides:

- **Java Runtime Management**: Automated Java installation and verification
- **Minecraft Server Installation**: Download and setup of various Minecraft server types
- **Installation State Tracking**: Persistent state management for installed instances
- **Multi-Server Support**: Support for multiple server instances in separate directories

## Architecture

### Location
```
internal/instance_installation/
```

### Dependencies
- **MineDB**: For persistent state storage
- **HttpDownloader**: For downloading server files (from External Services Layer)
- **File System**: For directory and file operations

## Configuration

### InstallConfig Structure
```go
type InstallConfig struct {
    // URL mapping for different server mod loaders
    UrlPerServerType map[ServerModLoader][]string
    
    // Base path where server instances are created  
    StoreInstancePath string
}
```

### ServerModLoader Type
```go
// Defined in internal/common/
type ServerModLoader string

const (
    Vanilla    ServerModLoader = "vanilla"
    Forge      ServerModLoader = "forge"
    Fabric     ServerModLoader = "fabric"
    Quilt      ServerModLoader = "quilt"
    Paper      ServerModLoader = "paper"
    Spigot     ServerModLoader = "spigot"
    Bukkit     ServerModLoader = "bukkit"
)
```

## API Interface

### Constructor
```go
func NewInstanceInstallationModule(db MineDB, config InstallConfig) (*InstanceInstallation, error)
```

### Core Methods

#### Java Management
```go
// Install the latest Java runtime
installJava() error

// Check if Java is installed and available
checkJava() error
```

#### Minecraft Installation
```go
// Install Minecraft server of specified version and type
installMinecraft(targetInstancePath, version string, serverType ServerModLoader) error

// Complete installation process (Java check + Minecraft installation)  
Install(targetInstancePath, version string, serverType ServerModLoader) error
```

## Data Storage

### Storage Key Format
```
instanceinstallation:<targetInstancePath>
```

### Storage Schema
```go
type InstanceInstallationData struct {
    // Minecraft server version
    Version string `json:"version"`
    
    // Server mod loader type
    ServerType ServerModLoader `json:"server_type"`
    
    // Path to Java runtime
    JavaPath string `json:"java_path"`
    
    // Installation timestamp
    InstalledAt time.Time `json:"installed_at"`
    
    // Installation status
    Status InstallationStatus `json:"status"`
    
    // Server JAR file path
    ServerJarPath string `json:"server_jar_path"`
}

type InstallationStatus string

const (
    StatusPending    InstallationStatus = "pending"
    StatusInProgress InstallationStatus = "in_progress"
    StatusCompleted  InstallationStatus = "completed"
    StatusFailed     InstallationStatus = "failed"
)
```

## Implementation Flow

### Installation Process
1. **Java Verification**
   - Check if Java is available (`checkJava()`)
   - Install Java if missing (`installJava()`)
   - Validate Java version compatibility

2. **Directory Setup**
   - Create target instance directory
   - Set up proper directory structure
   - Ensure proper permissions

3. **Server Download**
   - Resolve download URL based on server type and version
   - Download server JAR file
   - Verify file integrity (checksums if available)

4. **State Persistence**
   - Store installation data in MineDB
   - Update installation status
   - Record installation metadata

### Directory Structure
```
<StoreInstancePath>/
├── instance1/
│   ├── server.jar
│   ├── server.properties
│   ├── eula.txt
│   └── logs/
├── instance2/
│   └── ...
```

## Error Handling

### Error Types
- **JavaInstallationError**: Java installation failures
- **DownloadError**: Server file download failures  
- **FileSystemError**: Directory/file operation failures
- **ConfigurationError**: Invalid configuration parameters
- **StateError**: MineDB operation failures

### Error Recovery
- **Partial Installation Cleanup**: Remove incomplete installations
- **Retry Logic**: Configurable retry attempts for network operations
- **State Consistency**: Ensure MineDB state matches file system state

## Usage Examples

### Basic Installation
```go
config := InstallConfig{
    StoreInstancePath: "/opt/minecraft/instances",
    UrlPerServerType: map[ServerModLoader][]string{
        Vanilla: {
            "https://launcher.mojang.com/v1/objects/{version}/server.jar",
        },
        Paper: {
            "https://api.papermc.io/v2/projects/paper/versions/{version}/builds/latest/downloads/paper-{version}-{build}.jar",
        },
    },
}

installer, err := NewInstanceInstallationModule(db, config)
if err != nil {
    log.Fatal(err)
}

err = installer.Install("/opt/minecraft/instances/server1", "1.19.2", Vanilla)
if err != nil {
    log.Printf("Installation failed: %v", err)
}
```

### Checking Installation Status
```go
data, err := db.Get("instanceinstallation:/opt/minecraft/instances/server1")
if err != nil {
    log.Printf("Instance not found: %v", err)
    return
}

installData := data.(InstanceInstallationData)
fmt.Printf("Server %s (v%s) - Status: %s\n", 
    installData.ServerType, 
    installData.Version, 
    installData.Status)
```

## Testing Strategy

### Unit Tests
- Java installation and verification logic
- URL resolution for different server types
- State management operations
- Error handling scenarios

### Integration Tests
- End-to-end installation process
- File system operations
- Network download operations
- MineDB integration

### Mocking Strategy
- Mock HTTP downloader for network operations
- Mock file system operations for testing
- Mock MineDB for state testing

## Security Considerations

### Download Security
- **URL Validation**: Validate download URLs to prevent SSRF attacks
- **File Verification**: Verify downloaded files using checksums
- **Sandboxing**: Isolate download operations

### File System Security
- **Path Validation**: Prevent directory traversal attacks
- **Permission Management**: Set appropriate file/directory permissions
- **Cleanup**: Secure cleanup of temporary files

## Future Enhancements

### Optimization Strategies
- **Concurrent Downloads**: Parallel downloading when installing multiple instances
- **Caching**: Cache frequently downloaded server versions
- **Incremental Updates**: Support for incremental server updates

### Resource Management
- **Disk Space**: Monitor and validate available disk space
- **Network Bandwidth**: Throttle downloads to prevent network saturation
- **Memory Usage**: Efficient handling of large file downloads

### Other Planned Features
- **Version Caching**: Cache downloaded server versions for reuse
- **Integrity Verification**: SHA256 checksum verification
- **Update Management**: In-place server updates
- **Rollback Support**: Ability to rollback failed installations

### Integration Points
- **Backup Integration**: Automatic backup before major updates
- **Monitoring Integration**: Installation progress tracking
- **Notification Integration**: Installation completion notifications
