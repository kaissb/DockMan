# Advanced Docker Management UI Tool - Development Prompt for DeepSeek

## Project Overview
Create a comprehensive Docker management tool in Go with a modern web UI that provides essential container orchestration features similar to Dokploy/Coolify but with a streamlined focus on core functionality.

## Technical Stack Requirements
- **Backend**: Go (Gin framework recommended)
- **Frontend**: SvelteKit with TypeScript
- **Database**: SQLite for development, PostgreSQL support for production
- **Container Management**: Docker Engine API, Docker Compose
- **Reverse Proxy**: Nginx with automatic SSL via Let's Encrypt
- **Real-time Communication**: WebSockets for logs and terminal access
- **Repository Structure**: Monorepo for easy open-source contribution
- **API Authentication**: JWT-based API keys for programmatic access

## Core Architecture

### 1. Project Structure (Monorepo)
```
docker-manager/
├── apps/
│   ├── api/                 # Go backend API
│   │   ├── cmd/
│   │   │   └── server/
│   │   │       └── main.go
│   │   ├── internal/
│   │   │   ├── api/
│   │   │   ├── docker/
│   │   │   ├── nginx/
│   │   │   ├── ssl/
│   │   │   ├── database/
│   │   │   ├── auth/
│   │   │   └── websocket/
│   │   ├── migrations/
│   │   └── configs/
│   └── web/                 # SvelteKit frontend
│       ├── src/
│       │   ├── lib/
│       │   ├── routes/
│       │   └── app.html
│       ├── static/
│       ├── package.json
│       └── svelte.config.js
├── packages/
│   └── shared/              # Shared types and utilities
│       ├── types/
│       └── utils/
├── templates/               # Docker compose templates
│   ├── databases/
│   └── applications/
├── scripts/                 # Development and build scripts
├── docs/                    # Documentation
├── .github/                 # GitHub workflows and templates
├── docker-compose.dev.yml   # Development environment
├── Makefile                 # Development tasks
├── README.md
├── CONTRIBUTING.md
└── LICENSE

### 2. Data Models
```go
type Project struct {
    ID          uint         `json:"id"`
    Name        string       `json:"name"`
    Description string       `json:"description"`
    Environments []Environment `json:"environments"`
    CreatedAt   time.Time    `json:"created_at"`
}

type Environment struct {
    ID          uint      `json:"id"`
    ProjectID   uint      `json:"project_id"`
    Name        string    `json:"name"`
    Variables   []EnvVar  `json:"variables"`
    Services    []Service `json:"services"`
    IsActive    bool      `json:"is_active"`
}

type Service struct {
    ID            uint      `json:"id"`
    EnvironmentID uint      `json:"environment_id"`
    Name          string    `json:"name"`
    Type          string    `json:"type"` // docker, compose, database
    Image         string    `json:"image,omitempty"`
    ComposeFile   string    `json:"compose_file,omitempty"`
    Domain        string    `json:"domain,omitempty"`
    Port          int       `json:"port,omitempty"`
    Status        string    `json:"status"`
    Variables     []EnvVar  `json:"variables"`
}

type APIKey struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`
    Key         string    `json:"key"`
    ProjectID   *uint     `json:"project_id,omitempty"` // nil for global access
    Permissions []string  `json:"permissions"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty"`
    LastUsedAt  *time.Time `json:"last_used_at"`
    CreatedAt   time.Time `json:"created_at"`
    IsActive    bool      `json:"is_active"`
}

type EnvVar struct {
    Key     string `json:"key"`
    Value   string `json:"value"`
    IsSecret bool  `json:"is_secret"`
}
```

## Feature Implementation Requirements

### 1. Docker Container Management
- **Image Operations**: Pull, list, remove images with progress tracking
- **Container Lifecycle**: Start, stop, restart, remove containers
- **Real-time Logs**: Stream container logs via WebSocket
- **Interactive Terminal**: Web-based terminal access to running containers
- **Resource Monitoring**: CPU, memory, disk usage for containers

### 2. Docker Compose Management
- **Compose File Editor**: Syntax-highlighted YAML editor
- **Service Operations**: Up, down, restart, rebuild individual services
- **Multi-service Deployment**: Deploy entire compose stacks
- **Health Checks**: Monitor service health and dependencies
- **Log Aggregation**: Centralized logging for all compose services

### 3. Environment Variable Management
- **Hierarchical Variables**: Project → Environment → Service level variables
- **Secret Management**: Encrypted storage for sensitive variables
- **Variable Inheritance**: Environment variables cascade to services
- **Template Variables**: Support for variable substitution in configs
- **Bulk Operations**: Import/export variables via JSON/YAML

### 4. Domain Management & SSL
- **Automatic Nginx Configuration**: Generate reverse proxy configs
- **Let's Encrypt Integration**: Automatic SSL certificate provisioning
- **Domain Validation**: DNS validation before certificate issuance
- **Certificate Renewal**: Automated certificate renewal process
- **Multi-domain Support**: Handle multiple domains per service

### 5. Template System
```go
type Template struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Category    string `json:"category"` // database, application, monitoring
    Description string `json:"description"`
    ComposeFile string `json:"compose_file"`
    Variables   []TemplateVar `json:"variables"`
    Icon        string `json:"icon"`
}

type TemplateVar struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Type        string `json:"type"` // string, number, boolean, select
    Default     string `json:"default"`
    Required    bool   `json:"required"`
    Options     []string `json:"options,omitempty"`
}
```

**Pre-built Templates**:
- **Databases**: PostgreSQL, MySQL, MongoDB, Redis, InfluxDB
- **Web Servers**: Nginx, Apache, Caddy
- **Applications**: WordPress, Ghost, Grafana, Prometheus
- **Development**: Code Server, GitLab, Jenkins

### 6. Database Management
- **Database Deployment**: One-click database deployment from templates
- **Connection Management**: Built-in database connection testing
- **Backup Operations**: Automated backup scheduling and restoration
- **User Management**: Database user creation and permission management
- **Monitoring**: Database performance metrics and alerts

### 8. API Key Management
- **API Key Generation**: Create scoped API keys with specific permissions
- **Permission System**: Fine-grained permissions (read, write, deploy, admin)
- **Project Scoping**: Keys can be scoped to specific projects or global
- **Rate Limiting**: Configurable rate limits per API key
- **Usage Analytics**: Track API key usage and statistics
- **Key Rotation**: Support for key rotation and expiration

```go
type APIKeyService struct {
    permissions map[string][]string // key -> permissions
}

func (s *APIKeyService) ValidatePermission(apiKey, action string) bool {
    perms, exists := s.permissions[apiKey]
    if !exists {
        return false
    }
    return contains(perms, action) || contains(perms, "admin")
}
```
- **Project Creation**: Initialize projects with default environments
- **Environment Cloning**: 
  ```go
  func (s *EnvironmentService) CloneEnvironment(sourceEnvID, targetProjectID uint, newName string) error {
      // Copy all services, variables, and configurations
      // Deploy all services to new environment
      // Maintain service dependencies and order
  }
  ```
- **Environment Comparison**: Side-by-side environment configuration comparison
- **Deployment History**: Track deployment history and rollback capability

## API Design

### REST Endpoints
```
# Projects
GET    /api/projects
POST   /api/projects
GET    /api/projects/:id
PUT    /api/projects/:id
DELETE /api/projects/:id

# Environments
GET    /api/projects/:id/environments
POST   /api/projects/:id/environments
POST   /api/environments/:id/clone
GET    /api/environments/:id/services

# Services
POST   /api/environments/:id/services
GET    /api/services/:id
PUT    /api/services/:id
DELETE /api/services/:id
POST   /api/services/:id/deploy
POST   /api/services/:id/restart
GET    /api/services/:id/logs

# Docker Operations
GET    /api/docker/images
POST   /api/docker/images/pull
GET    /api/docker/containers
POST   /api/docker/containers/:id/start
POST   /api/docker/containers/:id/stop

# Templates
GET    /api/templates
GET    /api/templates/:id
POST   /api/templates/:id/deploy

# API Keys
GET    /api/keys
POST   /api/keys
DELETE /api/keys/:id
PUT    /api/keys/:id/toggle
GET    /api/keys/:id/usage
```

### WebSocket Endpoints
```
/ws/logs/:serviceId          - Real-time log streaming
/ws/terminal/:containerId    - Interactive terminal
/ws/deployments/:envId       - Deployment progress
/ws/system/stats             - System resource monitoring
```

## SvelteKit Frontend Architecture

### 1. Project Structure
```
apps/web/src/
├── lib/
│   ├── components/        # Reusable UI components
│   │   ├── ui/           # Base UI components
│   │   ├── forms/        # Form components
│   │   ├── charts/       # Chart components
│   │   └── terminal/     # Terminal emulator
│   ├── stores/           # Svelte stores
│   │   ├── auth.ts
│   │   ├── projects.ts
│   │   ├── services.ts
│   │   └── websocket.ts
│   ├── api/              # API client
│   ├── types/            # TypeScript types
│   └── utils/            # Utility functions
├── routes/
│   ├── +layout.svelte
│   ├── +page.svelte
│   ├── projects/
│   ├── templates/
│   └── settings/
└── app.html
```

### 2. Real-time Features
```typescript
// Real-time service monitoring
import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export const serviceStats = writable<Map<string, ServiceStats>>(new Map());
export const deploymentProgress = writable<DeploymentStatus[]>([]);

if (browser) {
    const ws = new WebSocket('ws://localhost:8080/ws/system/stats');
    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        serviceStats.update(stats => stats.set(data.serviceId, data));
    };
}
```

## UI/UX Requirements (SvelteKit)

### 1. Dashboard Layout
- **Sidebar Navigation**: Projects, Templates, System Settings
- **Main Content Area**: Environment overview, service grid
- **Status Bar**: System health, active deployments, resource usage
- **Real-time Updates**: Live status updates without page refresh

### 2. Service Management Interface
- **Service Cards**: Visual representation of service status and health
- **Quick Actions**: Deploy, restart, stop, logs, terminal buttons
- **Resource Graphs**: CPU, memory, network usage charts
- **Log Viewer**: Searchable, filterable log interface with timestamps

### 3. Environment Management
- **Environment Tabs**: Switch between environments within a project
- **Service Grid**: Drag-and-drop service organization
- **Variable Editor**: Tabular interface for environment variables
- **Deployment Pipeline**: Visual deployment progress indicator

### 4. Template Deployment Wizard
- **Template Gallery**: Categorized template browser with search
- **Configuration Form**: Dynamic form generation based on template variables
- **Deployment Preview**: Show generated docker-compose before deployment
- **Progress Tracking**: Real-time deployment progress with detailed steps

## Open Source Strategy & Proprietary Features

### Open Source Core (Community Edition)
**Included Features:**
- Basic Docker container management
- Docker Compose orchestration
- Local environment management
- Basic template system
- Self-hosted domain management
- Basic monitoring and logs
- API access with rate limiting
- Single-server deployment

### Proprietary Features (Cloud/Enterprise Edition)
**Premium Features to Monetize:**
1. **Multi-Server Management**: Deploy across multiple Docker hosts/clusters
2. **Advanced Monitoring**: Detailed analytics, alerting, performance insights
3. **Team Collaboration**: User management, permissions, audit logs
4. **Cloud Integrations**: AWS/GCP/Azure deployment, managed databases
5. **Enterprise Templates**: Advanced application stacks, compliance templates
6. **Backup & Disaster Recovery**: Automated backups, point-in-time recovery
7. **Advanced Security**: SSO integration, compliance reports, vulnerability scanning
8. **Custom Branding**: White-label solutions for agencies
9. **Priority Support**: 24/7 support, dedicated account management
10. **Advanced API**: Higher rate limits, webhooks, advanced automation

### Implementation Strategy
```go
// Feature flag system for open source vs proprietary
type FeatureFlag struct {
    Name        string `json:"name"`
    Enabled     bool   `json:"enabled"`
    Tier        string `json:"tier"` // "community", "cloud", "enterprise"
    Description string `json:"description"`
}

func (f *FeatureService) IsEnabled(feature string) bool {
    if flag, exists := f.flags[feature]; exists {
        return flag.Enabled && f.hasAccess(flag.Tier)
    }
    return false
}
```

## Contributing Guidelines & Developer Experience

### 1. Development Setup
```bash
# Quick start script
make setup          # Install dependencies
make dev            # Start development servers
make test           # Run all tests
make build          # Build for production
make docker-dev     # Start with Docker
```

### 2. Contributing Workflow
```
.github/
├── workflows/
│   ├── ci.yml           # Run tests on PR
│   ├── build.yml        # Build and release
│   └── security.yml     # Security scanning
├── ISSUE_TEMPLATE/
│   ├── bug_report.md
│   ├── feature_request.md
│   └── improvement.md
├── PULL_REQUEST_TEMPLATE.md
└── CONTRIBUTING.md
```

### 3. Code Quality Standards
- **Linting**: ESLint for TypeScript, golangci-lint for Go
- **Testing**: Vitest for frontend, Go testing framework
- **Type Safety**: Full TypeScript coverage, Go struct validation
- **Documentation**: JSDoc for components, Go doc for packages
- **Conventional Commits**: Standardized commit messages

### 4. Development Tools
```json
// package.json scripts
{
  "scripts": {
    "dev": "vite dev",
    "build": "vite build",
    "test": "vitest run",
    "test:watch": "vitest",
    "lint": "eslint . && prettier --check .",
    "format": "prettier --write .",
    "api:dev": "cd ../api && go run cmd/server/main.go"
  }
}
```

## Security Considerations
- **Authentication**: JWT-based authentication system
- **Authorization**: Role-based access control (Admin, Developer, Viewer)
- **Secret Management**: Encrypt sensitive environment variables
- **Docker Socket Security**: Secure Docker daemon access
- **SSL/TLS**: HTTPS enforcement with proper certificate validation
- **Input Validation**: Sanitize all user inputs, especially YAML/JSON

## Performance Requirements
- **Concurrent Operations**: Handle multiple simultaneous deployments
- **Resource Monitoring**: Efficient system resource tracking
- **Database Optimization**: Indexed queries for large numbers of services
- **Caching**: Redis caching for frequently accessed data
- **WebSocket Management**: Efficient connection pooling and cleanup

## Error Handling & Logging
- **Structured Logging**: JSON-formatted logs with contextual information
- **Error Recovery**: Graceful handling of Docker daemon disconnections
- **Deployment Rollback**: Automatic rollback on deployment failures
- **Health Checks**: Proactive service health monitoring
- **Audit Trail**: Complete audit log of all system operations

## Development Guidelines
1. **Code Organization**: Clean architecture with dependency injection
2. **Testing**: Unit tests for all business logic, integration tests for Docker operations
3. **Documentation**: OpenAPI/Swagger documentation for all endpoints
4. **Configuration**: Environment-based configuration management
5. **Monitoring**: Prometheus metrics integration for system monitoring
6. **Database Migrations**: Versioned database schema migrations
7. **Docker Best Practices**: Multi-stage builds, minimal base images

## Advanced Features (Phase 2)
- **Git Integration**: Deploy directly from Git repositories
- **CI/CD Pipelines**: Basic build and deployment pipelines
- **Monitoring Integration**: Grafana dashboard auto-generation
- **Backup Automation**: Scheduled backups with retention policies
- **Multi-server Support**: Manage Docker across multiple hosts
- **Resource Quotas**: CPU/memory limits per project/environment

## Implementation Priority
1. **Core Docker Management** (Week 1-2)
2. **Monorepo Setup & SvelteKit Frontend** (Week 3)
3. **Project/Environment Structure** (Week 4)
4. **API Key System** (Week 5)
5. **Template System** (Week 6)
6. **Domain/SSL Management** (Week 7)
7. **Database Management** (Week 8)
8. **Documentation & Contributing Setup** (Week 9)
9. **Testing & Polish** (Week 10)
10. **Feature Flag System for Monetization** (Week 11-12)

## Monorepo Development Workflow

### Package Management
```json
// Root package.json for workspace management
{
  "name": "docker-manager",
  "private": true,
  "workspaces": [
    "apps/*",
    "packages/*"
  ],
  "scripts": {
    "dev": "concurrently \"npm run dev:api\" \"npm run dev:web\"",
    "dev:api": "cd apps/api && go run cmd/server/main.go",
    "dev:web": "cd apps/web && npm run dev",
    "build": "npm run build:web && npm run build:api",
    "test": "npm run test:web && npm run test:api"
  }
}
```

### Shared Types Package
```typescript
// packages/shared/types/index.ts
export interface Project {
  id: number;
  name: string;
  environments: Environment[];
}

export interface DeploymentRequest {
  environmentId: number;
  services: ServiceConfig[];
}
```

This tool should provide a production-ready, open-source foundation with clear monetization opportunities through advanced features while maintaining an excellent developer experience for contributors.