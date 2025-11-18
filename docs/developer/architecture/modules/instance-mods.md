# InstanceMods Module

All code that provided in that docs is just an example. For real code check sources.

## Overview

The InstanceMods module manages Minecraft modifications (mods) for server instances. This document covers the research phase analysis of different mod management approaches and the recommended implementation strategy.

## Purpose and Scope

This module provides:

- **Mod Installation**: Install mods from various sources
- **Mod Management**: Enable/disable, update, and remove mods
- **Dependency Resolution**: Handle mod dependencies automatically
- **Version Compatibility**: Ensure mod compatibility with Minecraft versions
- **Security Validation**: Validate mod files for security

## Research Phase: Mod Management Approaches

### 1. Manual Format

**Description**: Users provide direct download URLs for JAR files.

**Advantages**:
- Simple implementation
- No external API dependencies
- Full user control over mod sources
- Works with any mod distribution method

**Disadvantages**:
- No automatic dependency resolution
- Security risks from unvalidated sources
- Manual version management
- No metadata about mods

**Implementation Requirements**:
- JAR file validation
- Checksum verification
- Malware scanning
- File type validation

### 2. Semi-Automatic Format

**Description**: Users provide URLs from known, trusted sites (CurseForge, Modrinth).

**Advantages**:
- Better security through trusted sources
- Some metadata extraction possible
- Reduced manual work
- Community-trusted platforms

**Disadvantages**:
- Limited to specific platforms
- Still requires manual URL finding
- Partial automation only

**Implementation Requirements**:
- URL pattern validation for trusted domains
- Site-specific parsing logic
- Metadata extraction
- API integration for enhanced features

### 3. Automatic Format

**Description**: Users provide mod names, system automatically finds and downloads from APIs.

**Advantages**:
- Full automation
- Automatic dependency resolution
- Version compatibility checking
- Rich metadata and descriptions
- Update notifications

**Disadvantages**:
- Complex implementation
- API rate limiting
- Dependency on external services
- Potential API changes

**API Research Results**:

#### CurseForge API
- **Endpoint**: `https://api.curseforge.com/`
- **Authentication**: API key required
- **Rate Limits**: 1000 requests/hour (free tier)
- **Features**: Mod search, file downloads, dependency info
- **Reliability**: High, maintained by Overwolf

#### Modrinth API  
- **Endpoint**: `https://api.modrinth.com/`
- **Authentication**: Optional (rate limits apply)
- **Rate Limits**: 300 requests/minute (authenticated)
- **Features**: Mod search, version management, dependencies
- **Reliability**: High, open-source community

#### GitHub Releases API
- **Endpoint**: `https://api.github.com/repos/{owner}/{repo}/releases`
- **Authentication**: Optional (higher rate limits when authenticated)
- **Rate Limits**: 60 requests/hour (unauthenticated), 5000 (authenticated)
- **Features**: Release management, asset downloads
- **Reliability**: Very high (GitHub infrastructure)

## Recommended Architecture: Hybrid Approach

Based on the research, we recommend implementing a **hybrid approach** that supports all three methods with automatic format as the primary method.

### Architecture Design

```go
type ModInstaller interface {
    // Install mod using different strategies
    InstallByURL(serverPath, url string) (*ModInfo, error)
    InstallByName(serverPath, modName string, options InstallOptions) (*ModInfo, error)
    InstallFromFile(serverPath, filePath string) (*ModInfo, error)
    
    // Mod management
    UpdateMod(serverPath, modID string) error
    RemoveMod(serverPath, modID string) error
    ListMods(serverPath string) ([]ModInfo, error)
    
    // Dependencies
    ResolveDependencies(serverPath string, modInfo ModInfo) ([]ModInfo, error)
}

type InstallOptions struct {
    MinecraftVersion string
    ModLoader        ModLoader
    PreferredSource  ModSource
    IncludeDeps      bool
    ForceUpdate      bool
}

type ModInfo struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    Version      string    `json:"version"`
    FileName     string    `json:"file_name"`
    Source       ModSource `json:"source"`
    Dependencies []string  `json:"dependencies"`
    Checksum     string    `json:"checksum"`
    InstallTime  time.Time `json:"install_time"`
}

type ModSource string

const (
    SourceCurseForge ModSource = "curseforge"
    SourceModrinth   ModSource = "modrinth"
    SourceGitHub     ModSource = "github"
    SourceURL        ModSource = "url"
    SourceFile       ModSource = "file"
)

type ModLoader string

const (
    LoaderForge  ModLoader = "forge"
    LoaderFabric ModLoader = "fabric"
    LoaderQuilt  ModLoader = "quilt"
)
```

## Implementation Strategy

### Phase 1: Core Infrastructure
1. **MineDB Integration**: Store mod metadata and state
2. **File Management**: Mod directory structure and file handling
3. **Validation**: JAR file validation and security checks
4. **Basic Installation**: Manual URL and file-based installation

### Phase 2: API Integration
1. **CurseForge Integration**: Primary mod source
2. **Modrinth Integration**: Alternative mod source
3. **Search and Discovery**: Mod search functionality
4. **Metadata Management**: Rich mod information handling

### Phase 3: Advanced Features
1. **Dependency Resolution**: Automatic dependency handling
2. **Update Management**: Automatic mod updates
3. **Conflict Detection**: Detect and resolve mod conflicts
4. **Backup Integration**: Mod state backups

## Data Storage

### Storage Key Format
```
instancemods:<targetInstancePath>
```

### Storage Schema
```go
type InstanceModsData struct {
    ServerPath    string              `json:"server_path"`
    ModLoader     ModLoader           `json:"mod_loader"`
    MCVersion     string              `json:"minecraft_version"`
    InstalledMods map[string]ModInfo  `json:"installed_mods"`
    ModSources    []ModSourceConfig   `json:"mod_sources"`
    LastUpdate    time.Time           `json:"last_update"`
}

type ModSourceConfig struct {
    Source  ModSource `json:"source"`
    APIKey  string    `json:"api_key,omitempty"`
    BaseURL string    `json:"base_url"`
    Enabled bool      `json:"enabled"`
}
```

## API Clients

### CurseForge Client
```go
type CurseForgeClient struct {
    apiKey  string
    baseURL string
    client  *http.Client
}

func (c *CurseForgeClient) SearchMods(query string, gameVersion string, modLoader ModLoader) ([]ModSearchResult, error)
func (c *CurseForgeClient) GetModInfo(modID string) (*ModInfo, error)
func (c *CurseForgeClient) DownloadMod(modID string, fileID string) ([]byte, error)
func (c *CurseForgeClient) GetDependencies(modID string, fileID string) ([]Dependency, error)
```

### Modrinth Client
```go
type ModrinthClient struct {
    token   string
    baseURL string
    client  *http.Client
}

func (m *ModrinthClient) SearchProjects(query SearchQuery) ([]ProjectResult, error)
func (m *ModrinthClient) GetProject(id string) (*Project, error)
func (m *ModrinthClient) GetVersions(projectID string, filters VersionFilters) ([]Version, error)
func (m *ModrinthClient) DownloadVersion(versionID string) ([]byte, error)
```

## Security Implementation

### File Validation
```go
type ModValidator struct {
    allowedExtensions []string
    maxFileSize       int64
    scanners          []SecurityScanner
}

func (v *ModValidator) ValidateFile(filePath string) ValidationResult {
    // 1. Extension validation
    // 2. File size validation  
    // 3. JAR structure validation
    // 4. Malware scanning
    // 5. Checksum verification
}

type SecurityScanner interface {
    ScanFile(filePath string) ScanResult
}
```

### Trusted Sources
```go
var TrustedDomains = []string{
    "files.minecraftforge.net",
    "maven.minecraftforge.net", 
    "api.modrinth.com",
    "cdn.modrinth.com",
    "media.forgecdn.net",
    "github.com",
}

func (m *ModInstaller) validateURL(url string) error {
    parsedURL, err := url.Parse(url)
    if err != nil {
        return err
    }
    
    for _, domain := range TrustedDomains {
        if strings.HasSuffix(parsedURL.Host, domain) {
            return nil
        }
    }
    
    return ErrUntrustedSource
}
```

## Usage Examples

### Automatic Installation
```go
installer, err := NewInstanceModsModule(db, config)
if err != nil {
    log.Fatal(err)
}

options := InstallOptions{
    MinecraftVersion: "1.19.2",
    ModLoader:        LoaderForge,
    PreferredSource:  SourceCurseForge,
    IncludeDeps:      true,
}

modInfo, err := installer.InstallByName("/opt/minecraft/instances/server1", "JEI", options)
if err != nil {
    log.Printf("Failed to install mod: %v", err)
} else {
    log.Printf("Installed %s v%s", modInfo.Name, modInfo.Version)
}
```

### Manual URL Installation
```go
modInfo, err := installer.InstallByURL(
    "/opt/minecraft/instances/server1",
    "https://media.forgecdn.net/files/4567/890/jei-1.19.2-forge-11.6.0.1016.jar",
)
if err != nil {
    log.Printf("Failed to install from URL: %v", err)
}
```

### Dependency Resolution
```go
mods, err := installer.ListMods("/opt/minecraft/instances/server1")
if err != nil {
    log.Fatal(err)
}

for _, mod := range mods {
    deps, err := installer.ResolveDependencies("/opt/minecraft/instances/server1", mod)
    if err != nil {
        log.Printf("Failed to resolve dependencies for %s: %v", mod.Name, err)
        continue
    }
    
    for _, dep := range deps {
        log.Printf("Installing dependency: %s", dep.Name)
        _, err := installer.InstallByName("/opt/minecraft/instances/server1", dep.Name, options)
        if err != nil {
            log.Printf("Failed to install dependency %s: %v", dep.Name, err)
        }
    }
}
```

## Testing Strategy

### Unit Tests
- API client implementations
- File validation logic
- Dependency resolution algorithms
- Error handling scenarios

### Integration Tests
- End-to-end mod installation
- API rate limiting behavior
- File system operations
- Security validation

### Mock Implementations
```go
type MockCurseForgeClient struct {
    searchResults map[string][]ModSearchResult
    modInfos      map[string]*ModInfo
    downloads     map[string][]byte
}

func (m *MockCurseForgeClient) SearchMods(query string, gameVersion string, modLoader ModLoader) ([]ModSearchResult, error) {
    return m.searchResults[query], nil
}
```

## Configuration

### Module Configuration
```go
type ModsConfig struct {
    StoreInstancePath string               `json:"store_instance_path"`
    ModSources        []ModSourceConfig    `json:"mod_sources"`
    Security          SecurityConfig       `json:"security"`
    Cache             CacheConfig          `json:"cache"`
}

type SecurityConfig struct {
    EnableValidation bool     `json:"enable_validation"`
    TrustedDomains   []string `json:"trusted_domains"`
    MaxFileSize      int64    `json:"max_file_size"`
    EnableScanning   bool     `json:"enable_scanning"`
}

type CacheConfig struct {
    EnableCache   bool          `json:"enable_cache"`
    CacheDir      string        `json:"cache_dir"`
    MaxCacheSize  int64         `json:"max_cache_size"`
    CacheTTL      time.Duration `json:"cache_ttl"`
}
```

## Future Enhancements

### Planned Features
- **Mod Pack Support**: Install complete mod packs
- **Version Management**: Track and manage multiple mod versions
- **Conflict Resolution**: Automatic mod conflict detection and resolution
- **Performance Monitoring**: Monitor mod performance impact

### Advanced Features
- **AI-Powered Recommendations**: Suggest compatible mods
- **Automated Testing**: Test mod combinations automatically
- **Cloud Sync**: Synchronize mod configurations across instances
- **Marketplace Integration**: Direct integration with mod marketplaces
