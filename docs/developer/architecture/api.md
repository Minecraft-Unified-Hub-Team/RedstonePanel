# API Layer: RESTful Interface Architecture

All code that provided in that docs is just an example. For real code check sources.

This document describes the core architectural component of our application: the **API Layer**. It serves as a unified RESTful interface, enabling communication between external clients and internal modules in a decoupled, scalable, and maintainable way.

The fundamental principle of this architecture is: **All operations flow through the API layer.** External clients interact only with REST endpoints, and the API layer orchestrates calls to appropriate modules.

## API Design Principles

The RedstonePanel API follows RESTful design principles with clear resource-based URLs and standard HTTP methods.

### Core Principles
- **Resource-Based URLs**: Each endpoint represents a specific resource or collection
- **HTTP Methods**: Use standard verbs (GET, POST, PUT, PATCH, DELETE)
- **Status Codes**: Proper HTTP status codes for different scenarios
- **JSON Communication**: All requests and responses use JSON format
- **Idempotency**: Safe methods (GET) and idempotent methods (PUT, DELETE)
- **Error Handling**: Consistent error response format

### API Structure
```
/api/v1/
├── instances/          # Server instance management
├── mods/              # Mod management
├── plugins/           # Plugin management  
├── configs/           # Configuration management
└── system/            # System information and health
```

### Authentication & Authorization
- **API Key Authentication**: Required for all endpoints
- **Role-Based Access**: Different access levels (admin, user, readonly)
- **Rate Limiting**: Prevent abuse and ensure fair usage

---

## API Endpoints Specification

### Instance Management Endpoints

#### List All Instances
```http
GET /api/v1/instances
Authorization: Bearer <api_key>
```

**Response:**
```json
{
    "instances": [
        {
            "id": "server1",
            "name": "Creative Server",
            "path": "/opt/minecraft/instances/server1",
            "version": "1.19.2",
            "server_type": "vanilla",
            "status": "running",
            "player_count": 5,
            "created_at": "2023-01-15T10:30:00Z"
        }
    ]
}
```

#### Create New Instance
```http
POST /api/v1/instances
Content-Type: application/json
Authorization: Bearer <api_key>

{
    "name": "Survival Server",
    "minecraft_version": "1.19.2",
    "server_type": "paper",
    "java_config": {
        "max_memory": "2G",
        "min_memory": "1G"
    }
}
```

#### Get Instance Details
```http
GET /api/v1/instances/{instanceId}
Authorization: Bearer <api_key>
```

#### Start Instance
```http
POST /api/v1/instances/{instanceId}/start
Authorization: Bearer <api_key>
```

#### Stop Instance  
```http
POST /api/v1/instances/{instanceId}/stop
Authorization: Bearer <api_key>

{
    "graceful": true,
    "timeout": 30
}
```

### Configuration Management Endpoints

#### Get Server Configuration
```http
GET /api/v1/instances/{instanceId}/configs/{configType}
Authorization: Bearer <api_key>
```

**Example:** `GET /api/v1/instances/server1/configs/server.properties`

#### Update Server Configuration
```http
PUT /api/v1/instances/{instanceId}/configs/{configType}
Content-Type: application/json
Authorization: Bearer <api_key>

{
    "server-port": "25565",
    "gamemode": "survival",
    "difficulty": "normal",
    "max-players": "20"
}
```

### Mod Management Endpoints

#### List Mods for Instance
```http
GET /api/v1/instances/{instanceId}/mods
Authorization: Bearer <api_key>
```

#### Install Mod
```http
POST /api/v1/instances/{instanceId}/mods
Content-Type: application/json
Authorization: Bearer <api_key>

{
    "name": "JEI",
    "source": "curseforge",
    "minecraft_version": "1.19.2",
    "mod_loader": "forge",
    "include_dependencies": true
}
```

#### Install Mod by URL
```http
POST /api/v1/instances/{instanceId}/mods/install-url
Content-Type: application/json
Authorization: Bearer <api_key>

{
    "url": "https://media.forgecdn.net/files/4567/890/jei-1.19.2-forge-11.6.0.1016.jar",
    "name": "JEI"
}
```

#### Remove Mod
```http
DELETE /api/v1/instances/{instanceId}/mods/{modId}
Authorization: Bearer <api_key>
```

### System Endpoints

#### Health Check
```http
GET /api/v1/system/health
```

**Response:**
```json
{
    "status": "healthy",
    "version": "1.0.0",
    "uptime": "2h30m15s",
    "modules": {
        "minedb": "healthy",
        "instance_control": "healthy",
        "instance_installation": "healthy"
    }
}
```

## API Implementation Architecture

### Handler Structure
```go
// API handlers coordinate between HTTP layer and modules
type InstanceHandler struct {
    installModule   *instance_installation.Module
    controlModule   *instance_control.Module
    manifestModule  *instance_manifest.Module
    modsModule      *instance_mods.Module
}

func (h *InstanceHandler) CreateInstance(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    // 2. Validate input
    // 3. Call appropriate module methods
    // 4. Return structured response
}
```

### Middleware Stack
- **Authentication**: API key validation
- **Authorization**: Role-based access control  
- **Rate Limiting**: Request throttling
- **Logging**: Request/response logging
- **Error Handling**: Consistent error responses
- **Validation**: Input validation and sanitization

### Error Response Format
```json
{
    "error": {
        "code": "INSTANCE_NOT_FOUND",
        "message": "Instance with ID 'server1' not found",
        "details": {
            "instance_id": "server1",
            "timestamp": "2023-01-15T10:30:00Z"
        }
    }
}
```

### WebSocket Events (Optional)
For real-time updates:
```javascript
// Connect to WebSocket
ws://localhost:8080/api/v1/ws?token=<api_key>

// Server events
{
    "event": "instance.status.changed",
    "data": {
        "instance_id": "server1",
        "old_status": "starting",
        "new_status": "running"
    }
}
```

## Integration with Modules

The API layer acts as a facade, orchestrating calls to underlying modules:

```go
func (api *API) handleStartInstance(instanceId string) error {
    // 1. Validate instance exists (InstanceInstallation)
    exists, err := api.installModule.IsInstalled(instanceId)
    if err != nil || !exists {
        return ErrInstanceNotFound
    }
    
    // 2. Check current status (InstanceControl)
    state, err := api.controlModule.GetState(instanceId)
    if err != nil {
        return err
    }
    
    if state.Status == StatusRunning {
        return ErrInstanceAlreadyRunning
    }
    
    // 3. Start the instance (InstanceControl)
    return api.controlModule.Start(instanceId)
}
```

### Module Communication Flow

```
HTTP Request → API Handler → Module Method → MineDB → Response
```

1. **HTTP Request**: Client sends REST request to API endpoint
2. **Authentication & Validation**: Middleware validates request and permissions
3. **API Handler**: Routes request to appropriate handler method
4. **Module Orchestration**: Handler calls one or more module methods
5. **Data Persistence**: Modules interact with MineDB for state management
6. **Response**: Structured JSON response returned to client

### Dependency Injection Pattern

```go
type APIServer struct {
    // Module dependencies
    mineDB          *minedb.DB
    installModule   *instance_installation.Module
    controlModule   *instance_control.Module
    manifestModule  *instance_manifest.Module
    modsModule      *instance_mods.Module
    
    // Infrastructure
    router          *mux.Router
    middleware      *MiddlewareStack
    config          *APIConfig
}

func NewAPIServer(dependencies Dependencies) *APIServer {
    api := &APIServer{
        mineDB:         dependencies.MineDB,
        installModule:  instance_installation.NewModule(dependencies.MineDB, dependencies.InstallConfig),
        controlModule:  instance_control.NewModule(dependencies.MineDB),
        manifestModule: instance_manifest.NewModule(dependencies.MineDB, dependencies.ManifestConfig),
        modsModule:     instance_mods.NewModule(dependencies.MineDB, dependencies.ModsConfig),
    }
    
    api.setupRoutes()
    return api
}
```

### Complete Example: Instance Creation Flow

```go
// POST /api/v1/instances
func (api *APIServer) CreateInstance(w http.ResponseWriter, r *http.Request) {
    var req CreateInstanceRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        api.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
        return
    }
    
    // Validate request
    if err := req.Validate(); err != nil {
        api.writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    instancePath := filepath.Join(api.config.InstancesPath, req.Name)
    
    // Step 1: Install instance (InstanceInstallation module)
    err := api.installModule.Install(instancePath, req.MinecraftVersion, req.ServerType)
    if err != nil {
        api.writeError(w, http.StatusInternalServerError, "Failed to install instance: " + err.Error())
        return
    }
    
    // Step 2: Configure server properties (InstanceManifest module)
    if req.ServerProperties != nil {
        err = api.manifestModule.WriteManifest(instancePath, common.ServerProperties, req.ServerProperties)
        if err != nil {
            // Cleanup on failure
            api.installModule.Remove(instancePath)
            api.writeError(w, http.StatusInternalServerError, "Failed to configure instance: " + err.Error())
            return
        }
    }
    
    // Step 3: Install initial mods if specified (InstanceMods module)
    for _, mod := range req.InitialMods {
        _, err = api.modsModule.InstallByName(instancePath, mod.Name, mod.Options)
        if err != nil {
            log.Printf("Warning: Failed to install mod %s: %v", mod.Name, err)
            // Continue with other mods, don't fail entire creation
        }
    }
    
    // Step 4: Get final instance state
    state, err := api.controlModule.GetState(instancePath)
    if err != nil {
        log.Printf("Warning: Failed to get instance state: %v", err)
        state = &common.ServerState{Status: common.StatusUnknown}
    }
    
    // Response
    response := CreateInstanceResponse{
        ID:          req.Name,
        Name:        req.Name,
        Path:        instancePath,
        Version:     req.MinecraftVersion,
        ServerType:  req.ServerType,
        Status:      string(state.Status),
        CreatedAt:   time.Now(),
    }
    
    api.writeJSON(w, http.StatusCreated, response)
}
```

This design ensures:
- **Clear separation of concerns**: API handles HTTP, modules handle business logic
- **Consistent error handling**: All errors flow through standard HTTP responses
- **Transactional operations**: Failed operations are properly cleaned up
- **Comprehensive logging**: All operations are logged for debugging
- **Type safety**: Strong typing throughout the request/response cycle
