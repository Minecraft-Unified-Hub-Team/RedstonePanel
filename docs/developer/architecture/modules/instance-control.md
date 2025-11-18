# InstanceControl Module

All code that provided in that docs is just an example. For real code check sources.

## Overview

The InstanceControl module manages the lifecycle of Minecraft server instances, providing functionality to start, stop, and monitor server processes. It maintains server state information and ensures proper process management across the application lifecycle.

## Purpose and Scope

This module provides:

- **Server Lifecycle Management**: Start and stop Minecraft server processes
- **Process Monitoring**: Track server status and health
- **State Management**: Maintain server state in MineDB
- **Resource Management**: Monitor server resource usage
- **Graceful Shutdown**: Handle clean server shutdowns

## Architecture

### Location
```
internal/instance_control/
```

### Dependencies
- **MineDB**: For server state persistence
- **InstanceManifest**: For server configuration access
- **Process Management**: OS-level process control
- **Logging**: Server output capture and logging

## API Interface

### Constructor
```go
func NewInstanceControlModule(db MineDB) (*InstanceControl, error)
```

### Core Methods

```go
// Start a Minecraft server instance
Start(targetServerPath string) error

// Stop a running Minecraft server instance  
Stop(targetServerPath string) error

// Get current state of a server instance
GetState(targetServerPath string) (ServerState, error)
```

## Server State Management

### Storage Key Format
```
instancecontrol:<targetInstancePath>
```

### State Schema
```go
type ServerState struct {
    // Server identification
    ServerPath string `json:"server_path"`
    
    // Current server status
    Status ServerStatus `json:"status"`
    
    // Process information
    ProcessID int `json:"process_id,omitempty"`
    
    // Server startup time
    StartTime time.Time `json:"start_time,omitempty"`
    
    // Server shutdown time
    StopTime time.Time `json:"stop_time,omitempty"`
    
    // Server metrics
    Metrics ServerMetrics `json:"metrics"`
    
    // Last error information
    LastError string `json:"last_error,omitempty"`
}

type ServerStatus string

const (
    StatusStopped  ServerStatus = "stopped"
    StatusStarting ServerStatus = "starting"
    StatusRunning  ServerStatus = "running"
    StatusStopping ServerStatus = "stopping"
    StatusError    ServerStatus = "error"
    StatusUnknown  ServerStatus = "unknown"
)

type ServerMetrics struct {
    // Memory usage in bytes
    MemoryUsage uint64 `json:"memory_usage"`
    
    // CPU usage percentage
    CPUUsage float64 `json:"cpu_usage"`
    
    // Number of connected players
    PlayerCount int `json:"player_count"`
    
    // Ticks per second
    TPS float64 `json:"tps"`
    
    // Last updated timestamp
    LastUpdate time.Time `json:"last_update"`
}
```

## Implementation Details

### Server Startup Process

#### Startup Validation
```go
func (ic *InstanceControl) validateStartup(serverPath string) error {
    // Check if server is already running
    state, err := ic.GetState(serverPath)
    if err == nil && state.Status == StatusRunning {
        return ErrServerAlreadyRunning
    }
    
    // Verify server installation
    if !ic.isServerInstalled(serverPath) {
        return ErrServerNotInstalled
    }
    
    // Check Java availability
    if !ic.isJavaAvailable(serverPath) {
        return ErrJavaNotFound
    }
    
    return nil
}
```

#### Process Execution
```go
type ServerProcess struct {
    cmd    *exec.Cmd
    stdin  io.WriteCloser
    stdout io.ReadCloser
    stderr io.ReadCloser
}

func (ic *InstanceControl) startProcess(serverPath string) (*ServerProcess, error) {
    // Build command with proper arguments
    cmd := exec.Command("java", ic.buildJavaArgs(serverPath)...)
    cmd.Dir = serverPath
    
    // Setup pipes for communication
    stdin, _ := cmd.StdinPipe()
    stdout, _ := cmd.StdoutPipe()
    stderr, _ := cmd.StderrPipe()
    
    // Start the process
    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("failed to start server: %w", err)
    }
    
    return &ServerProcess{
        cmd:    cmd,
        stdin:  stdin,
        stdout: stdout,
        stderr: stderr,
    }, nil
}
```

### Server Shutdown Process

#### Graceful Shutdown
```go
func (ic *InstanceControl) gracefulShutdown(serverPath string, timeout time.Duration) error {
    state, err := ic.GetState(serverPath)
    if err != nil {
        return err
    }
    
    if state.Status != StatusRunning {
        return ErrServerNotRunning
    }
    
    // Send stop command through stdin
    process := ic.getProcess(state.ProcessID)
    if process != nil {
        process.stdin.Write([]byte("stop\n"))
        
        // Wait for graceful shutdown with timeout
        done := make(chan error, 1)
        go func() {
            done <- process.cmd.Wait()
        }()
        
        select {
        case err := <-done:
            return err
        case <-time.After(timeout):
            // Force kill if timeout exceeded
            return ic.forceShutdown(state.ProcessID)
        }
    }
    
    return ErrProcessNotFound
}
```

#### Force Shutdown
```go
func (ic *InstanceControl) forceShutdown(processID int) error {
    process, err := os.FindProcess(processID)
    if err != nil {
        return err
    }
    
    return process.Kill()
}
```

### Process Monitoring

#### Health Check
```go
func (ic *InstanceControl) performHealthCheck(serverPath string) error {
    state, err := ic.GetState(serverPath)
    if err != nil {
        return err
    }
    
    if state.Status != StatusRunning {
        return nil
    }
    
    // Check if process is still alive
    if !ic.isProcessAlive(state.ProcessID) {
        // Update state to stopped
        state.Status = StatusStopped
        state.StopTime = time.Now()
        return ic.updateState(serverPath, state)
    }
    
    // Update metrics
    metrics, err := ic.collectMetrics(state.ProcessID)
    if err == nil {
        state.Metrics = metrics
        return ic.updateState(serverPath, state)
    }
    
    return nil
}
```

#### Metrics Collection
```go
func (ic *InstanceControl) collectMetrics(processID int) (ServerMetrics, error) {
    // Collect process metrics
    memUsage, cpuUsage, err := ic.getProcessMetrics(processID)
    if err != nil {
        return ServerMetrics{}, err
    }
    
    // Parse server logs for game-specific metrics
    playerCount, tps := ic.parseServerLogs()
    
    return ServerMetrics{
        MemoryUsage: memUsage,
        CPUUsage:    cpuUsage,
        PlayerCount: playerCount,
        TPS:         tps,
        LastUpdate:  time.Now(),
    }, nil
}
```

## Usage Examples

### Starting a Server
```go
control, err := NewInstanceControlModule(db)
if err != nil {
    log.Fatal(err)
}

// Start server
err = control.Start("/opt/minecraft/instances/server1")
if err != nil {
    log.Printf("Failed to start server: %v", err)
} else {
    log.Println("Server started successfully")
}
```

### Monitoring Server State
```go
state, err := control.GetState("/opt/minecraft/instances/server1") 
if err != nil {
    log.Printf("Failed to get server state: %v", err)
    return
}

fmt.Printf("Server Status: %s\n", state.Status)
if state.Status == StatusRunning {
    fmt.Printf("PID: %d\n", state.ProcessID)
    fmt.Printf("Uptime: %v\n", time.Since(state.StartTime))
    fmt.Printf("Players: %d\n", state.Metrics.PlayerCount)
    fmt.Printf("Memory: %d MB\n", state.Metrics.MemoryUsage/1024/1024)
}
```

### Graceful Server Shutdown
```go
err = control.Stop("/opt/minecraft/instances/server1")
if err != nil {
    log.Printf("Failed to stop server: %v", err)
} else {
    log.Println("Server stopped successfully")
}
```

## Error Handling

### Error Types
```go
var (
    ErrServerNotInstalled   = errors.New("server is not installed")
    ErrServerAlreadyRunning = errors.New("server is already running")
    ErrServerNotRunning     = errors.New("server is not running")
    ErrJavaNotFound         = errors.New("java runtime not found")
    ErrProcessNotFound      = errors.New("server process not found")
    ErrStartupTimeout       = errors.New("server startup timeout")
    ErrShutdownTimeout      = errors.New("server shutdown timeout")
)
```

### Error Recovery Strategies
- **State Reconciliation**: Automatically fix inconsistent states
- **Process Cleanup**: Clean up zombie processes
- **Retry Logic**: Retry failed operations with backoff
- **Rollback**: Revert to previous known good state

## Configuration Management

### Server Arguments
```go
func (ic *InstanceControl) buildJavaArgs(serverPath string) []string {
    args := []string{
        "-server",
        "-Xmx2G",
        "-Xms1G",
        "-XX:+UseG1GC",
        "-XX:+ParallelRefProcEnabled",
        "-jar", "server.jar",
        "--nogui",
    }
    
    // Load custom JVM args from manifest
    if customArgs := ic.loadJvmArgs(serverPath); customArgs != nil {
        args = append(customArgs, args[len(customArgs):]...)
    }
    
    return args
}
```

## Testing Strategy

### Unit Tests
- State management operations
- Process lifecycle management
- Error handling scenarios
- Metrics collection

### Integration Tests
- End-to-end server lifecycle
- Process monitoring accuracy
- State persistence verification
- Resource cleanup validation

### Mock Implementations
```go
type MockProcessManager struct {
    processes map[int]*MockProcess
}

func (m *MockProcessManager) Start(serverPath string) (*MockProcess, error) {
    // Mock process startup logic
}
```

## Security Considerations

### Process Security
- **User Isolation**: Run servers under dedicated users
- **Resource Limits**: Enforce CPU and memory limits
- **Network Security**: Restrict network access if needed

### Command Injection Prevention
- **Argument Sanitization**: Validate all command arguments
- **Path Validation**: Prevent directory traversal
- **Input Validation**: Validate server paths and commands

## Future Enhancements

### Resource Management
- **Memory Monitoring**: Track memory usage and prevent leaks
- **CPU Throttling**: Implement CPU usage limits
- **Disk I/O**: Monitor disk usage and performance

### Scalability
- **Concurrent Operations**: Support multiple server management
- **Resource Pooling**: Efficient resource allocation
- **Background Tasks**: Asynchronous monitoring and cleanup

### Other Planned Features
- **Auto-restart**: Automatic server restart on crashes
- **Load Balancing**: Distribute load across multiple servers
- **Performance Tuning**: Automatic JVM argument optimization
- **Backup Integration**: Automatic backups before restarts

### Advanced Features
- **Container Support**: Docker-based server isolation
- **Cluster Management**: Manage server clusters
- **Health Checks**: Advanced server health monitoring
- **Auto-scaling**: Dynamic server scaling based on load
