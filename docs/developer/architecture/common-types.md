# Common Types

All code that provided in that docs is just an example. For real code check sources.

## Overview

This document defines the common types and interfaces used across RedstonePanel modules. These types are defined in the `internal/common/` package and provide consistency across the application.

## Server Types

### ServerModLoader
```go
// ServerModLoader represents different Minecraft server mod loading systems
type ServerModLoader string

const (
    Vanilla    ServerModLoader = "vanilla"    // Vanilla Minecraft server
    Forge      ServerModLoader = "forge"      // Minecraft Forge
    Fabric     ServerModLoader = "fabric"     // Fabric mod loader
    Quilt      ServerModLoader = "quilt"      // Quilt mod loader (Fabric fork)
    Paper      ServerModLoader = "paper"      // Paper server (Bukkit/Spigot fork)
    Spigot     ServerModLoader = "spigot"     // Spigot server (Bukkit fork)
    Bukkit     ServerModLoader = "bukkit"     // Original Bukkit server
    Purpur     ServerModLoader = "purpur"     // Purpur server (Paper fork)
    Mohist     ServerModLoader = "mohist"     // Mohist (Forge + Bukkit hybrid)
)
```

### ServerManifestType
```go
// ServerManifestType represents different configuration file types
type ServerManifestType string

const (
    // Core server configuration
    ServerProperties ServerManifestType = "server.properties"
    EulaText         ServerManifestType = "eula.txt"
    
    // JVM configuration
    JvmArgs          ServerManifestType = "jvm.args"
    
    // Player management
    Whitelist        ServerManifestType = "whitelist.json"
    BannedPlayers    ServerManifestType = "banned-players.json"
    BannedIPs        ServerManifestType = "banned-ips.json"
    Ops              ServerManifestType = "ops.json"
    
    // World configuration (for Bukkit-based servers)
    BukkitYml        ServerManifestType = "bukkit.yml"
    SpigotYml        ServerManifestType = "spigot.yml"
    PaperYml         ServerManifestType = "paper.yml"
    
    // Plugin configurations
    PluginConfig     ServerManifestType = "plugin.yml"
    
    // Mod configurations (for modded servers)
    ModConfig        ServerManifestType = "mod.conf"
)
```

## Data Structures

### InstallationInfo
```go
// InstallationInfo contains information about a server installation
type InstallationInfo struct {
    ServerPath    string          `json:"server_path"`
    Version       string          `json:"version"`
    ServerType    ServerModLoader `json:"server_type"`
    JavaPath      string          `json:"java_path"`
    InstallTime   time.Time       `json:"install_time"`
    LastStartTime *time.Time      `json:"last_start_time,omitempty"`
    Status        ServerStatus    `json:"status"`
}

// ServerStatus represents the current state of a server instance
type ServerStatus string

const (
    StatusStopped   ServerStatus = "stopped"
    StatusStarting  ServerStatus = "starting"
    StatusRunning   ServerStatus = "running"
    StatusStopping  ServerStatus = "stopping"
    StatusError     ServerStatus = "error"
    StatusUnknown   ServerStatus = "unknown"
    StatusInstalling ServerStatus = "installing"
    StatusUpdating   ServerStatus = "updating"
)
```

### ModInfo
```go
// ModInfo contains information about an installed mod
type ModInfo struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    Version      string    `json:"version"`
    FileName     string    `json:"file_name"`
    FilePath     string    `json:"file_path"`
    Source       ModSource `json:"source"`
    SourceURL    string    `json:"source_url,omitempty"`
    Dependencies []string  `json:"dependencies"`
    Checksum     string    `json:"checksum"`
    ModLoader    ModLoader `json:"mod_loader"`
    MCVersions   []string  `json:"minecraft_versions"`
    InstallTime  time.Time `json:"install_time"`
    Enabled      bool      `json:"enabled"`
}

// ModSource represents where a mod was installed from
type ModSource string

const (
    SourceCurseForge ModSource = "curseforge"
    SourceModrinth   ModSource = "modrinth"
    SourceGitHub     ModSource = "github"
    SourceURL        ModSource = "url"
    SourceFile       ModSource = "file"
    SourceUnknown    ModSource = "unknown"
)

// ModLoader represents the mod loading system
type ModLoader string

const (
    LoaderForge  ModLoader = "forge"
    LoaderFabric ModLoader = "fabric"
    LoaderQuilt  ModLoader = "quilt"
    LoaderNone   ModLoader = "none"    // For vanilla servers
)
```

## Configuration Structures

### DatabaseConfig
```go
// DatabaseConfig configures the MineDB storage system
type DatabaseConfig struct {
    StoragePath   string        `json:"storage_path"`
    BackupDir     string        `json:"backup_dir"`
    BackupInterval time.Duration `json:"backup_interval"`
    AutoCleanup   bool          `json:"auto_cleanup"`
    MaxBackups    int           `json:"max_backups"`
}
```

### JavaConfig
```go
// JavaConfig contains Java runtime configuration
type JavaConfig struct {
    JavaHome        string   `json:"java_home,omitempty"`
    JavaVersion     string   `json:"java_version"`
    DefaultJvmArgs  []string `json:"default_jvm_args"`
    MaxMemory       string   `json:"max_memory"`
    MinMemory       string   `json:"min_memory"`
    AutoInstall     bool     `json:"auto_install"`
    DownloadSources []string `json:"download_sources"`
}
```

### InstanceConfig
```go
// InstanceConfig contains server instance configuration
type InstanceConfig struct {
    Name              string          `json:"name"`
    ServerPath        string          `json:"server_path"`
    MinecraftVersion  string          `json:"minecraft_version"`
    ServerType        ServerModLoader `json:"server_type"`
    JavaConfig        JavaConfig      `json:"java_config"`
    AutoStart         bool            `json:"auto_start"`
    AutoRestart       bool            `json:"auto_restart"`
    BackupBeforeStart bool            `json:"backup_before_start"`
}
```

## Interface Definitions

### MineDB Interface
```go
// MineDB provides centralized data storage for all modules
type MineDB interface {
    // Basic operations
    Get(key string) (any, error)
    Set(key string, value any) error
    Delete(key string) error
    
    // Batch operations
    GetBatch(keys []string) (map[string]any, error)
    SetBatch(data map[string]any) error
    
    // Utility operations
    Exists(key string) (bool, error)
    Keys(pattern string) ([]string, error)
    
    // Transaction support
    BeginTx() (Transaction, error)
}

// Transaction provides transactional database operations
type Transaction interface {
    Get(key string) (any, error)
    Set(key string, value any) error
    Delete(key string) error
    Commit() error
    Rollback() error
}
```

### ModuleInterface
```go
// Module defines the basic interface that all modules must implement
type Module interface {
    // Name returns the module name
    Name() string
    
    // Initialize performs module initialization
    Initialize() error
    
    // Start starts the module (if applicable)
    Start() error
    
    // Stop stops the module gracefully
    Stop() error
    
    // Health returns the current health status
    Health() HealthStatus
}

// HealthStatus represents module health
type HealthStatus struct {
    Status  string            `json:"status"`  // healthy, degraded, unhealthy
    Message string            `json:"message,omitempty"`
    Details map[string]string `json:"details,omitempty"`
}
```

## Error Types

### Common Errors
```go
// Common error types used across modules
var (
    // Generic errors
    ErrNotFound       = errors.New("resource not found")
    ErrAlreadyExists  = errors.New("resource already exists")
    ErrInvalidInput   = errors.New("invalid input")
    ErrPermission     = errors.New("permission denied")
    
    // Server-specific errors
    ErrServerNotInstalled   = errors.New("server is not installed")
    ErrServerAlreadyRunning = errors.New("server is already running")
    ErrServerNotRunning     = errors.New("server is not running")
    ErrInvalidServerPath    = errors.New("invalid server path")
    
    // Java-specific errors
    ErrJavaNotFound     = errors.New("java runtime not found")
    ErrJavaNotSupported = errors.New("java version not supported")
    
    // Mod-specific errors
    ErrModNotFound      = errors.New("mod not found")
    ErrModAlreadyExists = errors.New("mod already installed")
    ErrModIncompatible  = errors.New("mod incompatible with current setup")
    ErrDependencyFailed = errors.New("dependency resolution failed")
    
    // Configuration errors
    ErrInvalidConfig    = errors.New("invalid configuration")
    ErrConfigNotFound   = errors.New("configuration not found")
)
```

### Custom Error Types
```go
// ValidationError represents configuration validation errors
type ValidationError struct {
    Field   string `json:"field"`
    Value   any    `json:"value"`
    Message string `json:"message"`
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// ModuleError wraps errors with module context
type ModuleError struct {
    Module string
    Op     string
    Err    error
}

func (e ModuleError) Error() string {
    return fmt.Sprintf("%s.%s: %v", e.Module, e.Op, e.Err)
}

func (e ModuleError) Unwrap() error {
    return e.Err
}
```

## Utility Functions

### Key Generation
```go
// GenerateModuleKey creates a standardized key for module data storage
func GenerateModuleKey(moduleName, instancePath string) string {
    return fmt.Sprintf("%s:%s", strings.ToLower(moduleName), instancePath)
}

// ParseModuleKey extracts module name and instance path from a storage key
func ParseModuleKey(key string) (moduleName, instancePath string, err error) {
    parts := strings.SplitN(key, ":", 2)
    if len(parts) != 2 {
        return "", "", fmt.Errorf("invalid module key format: %s", key)
    }
    return parts[0], parts[1], nil
}
```

### Validation Helpers
```go
// ValidateServerPath validates a server installation path
func ValidateServerPath(path string) error {
    if path == "" {
        return errors.New("server path cannot be empty")
    }
    
    if !filepath.IsAbs(path) {
        return errors.New("server path must be absolute")
    }
    
    // Additional validation logic...
    return nil
}

// ValidateMinecraftVersion validates a Minecraft version string
func ValidateMinecraftVersion(version string) error {
    // Regex pattern for valid Minecraft versions (e.g., 1.19.2, 1.20-pre1, etc.)
    pattern := `^1\.\d+(\.\d+)?(-\w+\d*)?$`
    matched, err := regexp.MatchString(pattern, version)
    if err != nil {
        return err
    }
    if !matched {
        return fmt.Errorf("invalid Minecraft version format: %s", version)
    }
    return nil
}
```

## Usage in Modules

These common types should be imported and used consistently across all modules:

```go
package instance_installation

import (
    "github.com/redstone-panel/redstone-panel/internal/common"
)

func (ii *InstanceInstallation) Install(targetPath, version string, serverType common.ServerModLoader) error {
    // Use common types for consistency
    if err := common.ValidateServerPath(targetPath); err != nil {
        return err
    }
    
    if err := common.ValidateMinecraftVersion(version); err != nil {
        return err
    }
    
    // ... implementation
}
```
