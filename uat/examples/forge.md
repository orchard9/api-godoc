# Forge Development System API

HTTP REST API for the Forge Development System - AI-driven software development framework

## Overview

- **API Version**: 1.0.0
- **Specification Type**: OpenAPI 3.0.3
- **Base URL**: http://localhost:50052/api/v1
- **Generated**: 2025-07-15 21:15:41

## API Statistics

- **Total Resources**: 64
- **Total Operations**: 156
- **Total Endpoints**: 83
- **Resource Coverage**: 187%

## Resources

This section groups API endpoints by business resources for better understanding.

### :flagRegret

:flagRegret resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/tasks/{taskId}:flagRegret` | FlagTaskWithRegret flags a task with an implementation regret
Moves the task ... |

### :move

:move resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/tasks/{taskId}:move` | MoveTask changes a task's status (e.g., from todo to in-progress)
Enforces wo... |

### :regenerateIndices

:regenerateIndices resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/projects/{projectId}:regenerateIndices` | RegenerateIndices rebuilds all index files |

### :repair

:repair resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/projects/{projectId}:repair` | RepairProjectStructure repairs missing directories and files in the project s... |

### :resolveRegret

:resolveRegret resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/tasks/{taskId}:resolveRegret` | ResolveTaskRegret resolves a task's implementation regret
Moves the task back... |

### :understand

:understand resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/projects/{projectId}:understand` | UnderstandProject analyzes a project and provides streaming progress updates |

### :validate

:validate resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/projects/{projectId}:validate` | ValidateProjectStructure checks if the project structure is valid and complete |

### :validateTransition

:validateTransition resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/tasks/{taskId}:validateTransition` | ValidateTaskTransition checks if a state transition is valid without performi... |

### Activities

Activities resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/activities` | ListActivities retrieves activities based on filter criteria |

### Activities:recent

Activities:recent resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/activities:recent` | GetRecentActivities retrieves the most recent activities up to a specified limit |

### Agents

Agents resource operations

**Operations**: 3

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/agents` | ListAgents retrieves all available agents with their workload information |
| GET | `/api/v1/agents/{persona}/context` | GetAgentContext retrieves the context for a specific agent persona |
| PUT | `/api/v1/agents/{persona}/context` | UpdateAgentContext updates the context for a specific agent persona |

### Analysis

Analysis resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/documents/analysis/{analysisId}/status` | GetAnalysisStatus retrieves the status of a previously submitted analysis. |

### Answers

Answers resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/projects/{projectId}/answers` | StoreProjectAnswers stores question answers for a project |

### Audits

Audits resource operations

**Operations**: 6

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/audits/{auditId}` | GetAuditRequest retrieves a specific audit request |
| GET | `/api/v1/projects/{projectId}/audits` | ListAudits returns a list of audit requests with filtering support |
| POST | `/api/v1/audits/{auditId}/process` | ProcessAuditResults processes completed audit results and creates tasks |
| POST | `/api/v1/audits/{auditId}/results` | SubmitAuditResults submits audit deliverables and transitions to completed |
| POST | `/api/v1/projects/{projectId}/audits` | CreateAuditRequest creates a new audit request for a project |
| PUT | `/api/v1/audits/{auditId}/status` | UpdateAuditStatus updates the status of an audit request |

### BlockedBy

BlockedBy resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/tasks/{taskId}/blockedBy` | GetBlockedBy retrieves tasks that are blocked by a specific task
Useful for u... |

### Blockers

Blockers resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/tasks/{taskId}/blockers` | GetBlockers retrieves tasks that are blocking a specific task
Useful for unde... |

### Confidence

Confidence resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| PUT | `/api/v1/tasks/{taskId}/confidence` | UpdateConfidenceScore updates a task's confidence score
Includes rationale fo... |

### Containers

Containers resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/docker/containers` |  |

### Containers:build

Containers:build resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/docker/containers:build` |  |

### Containers:publish

Containers:publish resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/docker/containers:publish` |  |

### Containers:validate

Containers:validate resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/docker/containers:validate` |  |

### Context

Context resource operations

**Operations**: 2

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/agents/{persona}/context` | GetAgentContext retrieves the context for a specific agent persona |
| PUT | `/api/v1/agents/{persona}/context` | UpdateAgentContext updates the context for a specific agent persona |

### Cycles

Cycles resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/tasks/cycles` | DetectCycles finds all dependency cycles in the current task graph
Returns de... |

### Dashboard

Dashboard resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/projects/{projectId}/dashboard` | GetDashboardData returns all data needed for the dashboard in a single request |

### Dependencies

Dependencies resource operations

**Operations**: 3

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/tasks/{taskId}/dependencies` | GetDependencies retrieves all dependencies for a specific task |
| POST | `/api/v1/tasks/{taskId}/dependencies` | AddDependency adds a dependency relationship between tasks
Checks for circula... |
| DELETE | `/api/v1/tasks/{taskId}/dependencies/{dependencyId}` | RemoveDependency removes a dependency relationship between tasks |

### DependencyGraph

DependencyGraph resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/tasks/{taskId}/dependencyGraph` | GetTaskDependencies retrieves detailed dependency information for graph visua... |

### DependencyTree

DependencyTree resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/tasks/{taskId}/dependencyTree` | GetDependencyTree builds a hierarchical tree view of task dependencies
Suppor... |

### Docker

Docker resource operations

**Operations**: 5

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/docker/containers` |  |
| GET | `/api/v1/docker/status` | Docker protocol methods |
| POST | `/api/v1/docker/containers:build` |  |
| POST | `/api/v1/docker/containers:publish` |  |
| POST | `/api/v1/docker/containers:validate` |  |

### Documents

Documents resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/documents/analysis/{analysisId}/status` | GetAnalysisStatus retrieves the status of a previously submitted analysis. |

### Documents:analyze

Documents:analyze resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/documents:analyze` | AnalyzeDocuments processes submitted documents using AI to determine
which on... |

### Features

Features resource operations

**Operations**: 6

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/features` | ListFeatures retrieves features with optional filtering |
| GET | `/api/v1/features/{featureId}/tasks` | GetFeatureTasks retrieves all tasks associated with a feature |
| GET | `/api/v1/features/{id}` | GetFeature retrieves a specific feature by ID |
| POST | `/api/v1/features` | CreateFeature creates a new feature group |
| PUT | `/api/v1/features/{id}` | UpdateFeature updates an existing feature |
| DELETE | `/api/v1/features/{id}` | DeleteFeature removes a feature (administrative use only) |

### Files

Files resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/projects/{projectId}/understanding/files` | Get understanding for a specific file |

### Git

Git resource operations

**Operations**: 9

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/git/pullRequests/{branchName}` |  |
| GET | `/api/v1/git/releaseBranches/{version}` |  |
| GET | `/api/v1/git/taskBranches/{taskId}` |  |
| POST | `/api/v1/git/hooks:postCheckout` |  |
| POST | `/api/v1/git/hooks:postMerge` | Hook integration |
| POST | `/api/v1/git/pullRequests/release` |  |
| POST | `/api/v1/git/pullRequests/task` | Pull request operations |
| POST | `/api/v1/git/releaseBranches` | Release branch operations |
| POST | `/api/v1/git/taskBranches` | Task branch operations |

### History

History resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/simulations/{simulationName}/history` |  |

### Hooks:postCheckout

Hooks:postCheckout resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/git/hooks:postCheckout` |  |

### Hooks:postMerge

Hooks:postMerge resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/git/hooks:postMerge` | Hook integration |

### Info

Info resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| PUT | `/api/v1/projects/info` | UpdateProjectInfo updates only mutable fields (name and description) |

### Metadata

Metadata resource operations

**Operations**: 2

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/projects/metadata` | GetProjectMetadata retrieves project metadata |
| PUT | `/api/v1/projects/metadata` | UpdateProjectMetadata initializes or updates project metadata |

### Personas

Personas resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/personas` | ListPersonas retrieves all available personas in the system |

### Process

Process resource operations

**Operations**: 3

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/audits/{auditId}/process` | ProcessAuditResults processes completed audit results and creates tasks |
| POST | `/api/v1/requirements/{requirementId}/process` | ProcessRequirementResults processes completed requirement results |
| POST | `/api/v1/research/{researchId}/process` | ProcessResearchResults processes completed research results |

### Progress

Progress resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/releases/{version}/progress` |  |

### ProjectTypes

ProjectTypes resource operations

**Operations**: 2

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/projectTypes` | GetProjectTypes returns supported project types and languages |
| GET | `/api/v1/projectTypes/{projectType}/questions` | GetProjectQuestions returns questions based on project type |

### Projects

Projects resource operations

**Operations**: 23

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/projects` | ListProjects returns a paginated list of all projects |
| GET | `/api/v1/projects/metadata` | GetProjectMetadata retrieves project metadata |
| GET | `/api/v1/projects/{projectId}` | GetProject retrieves project information |
| GET | `/api/v1/projects/{projectId}/audits` | ListAudits returns a list of audit requests with filtering support |
| GET | `/api/v1/projects/{projectId}/dashboard` | GetDashboardData returns all data needed for the dashboard in a single request |
| GET | `/api/v1/projects/{projectId}/requirements` | ListRequirements returns a list of requirement requests with filtering support |
| GET | `/api/v1/projects/{projectId}/research` | ListResearch returns a list of research requests with filtering support |
| GET | `/api/v1/projects/{projectId}/stats` | GetProjectStats returns project-wide statistics and metrics |
| GET | `/api/v1/projects/{projectId}/understanding` | List all file understandings for a project |
| GET | `/api/v1/projects/{projectId}/understanding/files` | Get understanding for a specific file |
| GET | `/api/v1/projects/{projectId}/understanding/status` | Get current understanding status for a project |
| GET | `/api/v1/projects/{projectId}:validate` | ValidateProjectStructure checks if the project structure is valid and complete |
| POST | `/api/v1/projects/{projectId}/answers` | StoreProjectAnswers stores question answers for a project |
| POST | `/api/v1/projects/{projectId}/audits` | CreateAuditRequest creates a new audit request for a project |
| POST | `/api/v1/projects/{projectId}/requirements` | CreateRequirementRequest creates a new requirement request for a project |
| POST | `/api/v1/projects/{projectId}/research` | CreateResearchRequest creates a new research request for a project |
| POST | `/api/v1/projects/{projectId}/understanding:analyze` | Analyze project files to create AI-powered summaries with streaming progress |
| POST | `/api/v1/projects/{projectId}:regenerateIndices` | RegenerateIndices rebuilds all index files |
| POST | `/api/v1/projects/{projectId}:repair` | RepairProjectStructure repairs missing directories and files in the project s... |
| POST | `/api/v1/projects/{projectId}:understand` | UnderstandProject analyzes a project and provides streaming progress updates |
| PUT | `/api/v1/projects/info` | UpdateProjectInfo updates only mutable fields (name and description) |
| PUT | `/api/v1/projects/metadata` | UpdateProjectMetadata initializes or updates project metadata |
| PUT | `/api/v1/projects/{project.id}` | UpdateProject updates project information |

### Projects:scaffold

Projects:scaffold resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/projects:scaffold` | ScaffoldProject creates a new project with the initial directory structure |

### PullRequests

PullRequests resource operations

**Operations**: 3

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/git/pullRequests/{branchName}` |  |
| POST | `/api/v1/git/pullRequests/release` |  |
| POST | `/api/v1/git/pullRequests/task` | Pull request operations |

### Questions

Questions resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/projectTypes/{projectType}/questions` | GetProjectQuestions returns questions based on project type |

### RegretFlags

RegretFlags resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/tasks/{taskId}/regretFlags` | GetTaskRegretFlags retrieves all regret flags for a task
Returns the complete... |

### Release

Release resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/git/pullRequests/release` |  |

### ReleaseBranches

ReleaseBranches resource operations

**Operations**: 2

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/git/releaseBranches/{version}` |  |
| POST | `/api/v1/git/releaseBranches` | Release branch operations |

### Releases

Releases resource operations

**Operations**: 2

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/releases` | Release management methods |
| GET | `/api/v1/releases/{version}/progress` |  |

### Requirements

Requirements resource operations

**Operations**: 6

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/projects/{projectId}/requirements` | ListRequirements returns a list of requirement requests with filtering support |
| GET | `/api/v1/requirements/{requirementId}` | GetRequirementRequest retrieves a specific requirement request |
| POST | `/api/v1/projects/{projectId}/requirements` | CreateRequirementRequest creates a new requirement request for a project |
| POST | `/api/v1/requirements/{requirementId}/process` | ProcessRequirementResults processes completed requirement results |
| POST | `/api/v1/requirements/{requirementId}/results` | SubmitRequirementResults submits requirement deliverables and transitions to ... |
| PUT | `/api/v1/requirements/{requirementId}/status` | UpdateRequirementStatus updates the status of a requirement request |

### Research

Research resource operations

**Operations**: 6

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/projects/{projectId}/research` | ListResearch returns a list of research requests with filtering support |
| GET | `/api/v1/research/{researchId}` | GetResearchRequest retrieves a specific research request |
| POST | `/api/v1/projects/{projectId}/research` | CreateResearchRequest creates a new research request for a project |
| POST | `/api/v1/research/{researchId}/process` | ProcessResearchResults processes completed research results |
| POST | `/api/v1/research/{researchId}/results` | SubmitResearchResults submits research deliverables and transitions to completed |
| PUT | `/api/v1/research/{researchId}/status` | UpdateResearchStatus updates the status of a research request |

### Results

Results resource operations

**Operations**: 3

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/audits/{auditId}/results` | SubmitAuditResults submits audit deliverables and transitions to completed |
| POST | `/api/v1/requirements/{requirementId}/results` | SubmitRequirementResults submits requirement deliverables and transitions to ... |
| POST | `/api/v1/research/{researchId}/results` | SubmitResearchResults submits research deliverables and transitions to completed |

### Simulations

Simulations resource operations

**Operations**: 2

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/simulations` | Simulation management methods |
| GET | `/api/v1/simulations/{simulationName}/history` |  |

### Stats

Stats resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/projects/{projectId}/stats` | GetProjectStats returns project-wide statistics and metrics |

### Task

Task resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/git/pullRequests/task` | Pull request operations |

### TaskBranches

TaskBranches resource operations

**Operations**: 2

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/git/taskBranches/{taskId}` |  |
| POST | `/api/v1/git/taskBranches` | Task branch operations |

### Tasks

Tasks resource operations

**Operations**: 20

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/features/{featureId}/tasks` | GetFeatureTasks retrieves all tasks associated with a feature |
| GET | `/api/v1/tasks` | List tasks |
| GET | `/api/v1/tasks/cycles` | DetectCycles finds all dependency cycles in the current task graph
Returns de... |
| GET | `/api/v1/tasks/{taskId}` | Get task by ID |
| GET | `/api/v1/tasks/{taskId}/blockedBy` | GetBlockedBy retrieves tasks that are blocked by a specific task
Useful for u... |
| GET | `/api/v1/tasks/{taskId}/blockers` | GetBlockers retrieves tasks that are blocking a specific task
Useful for unde... |
| GET | `/api/v1/tasks/{taskId}/dependencies` | GetDependencies retrieves all dependencies for a specific task |
| GET | `/api/v1/tasks/{taskId}/dependencyGraph` | GetTaskDependencies retrieves detailed dependency information for graph visua... |
| GET | `/api/v1/tasks/{taskId}/dependencyTree` | GetDependencyTree builds a hierarchical tree view of task dependencies
Suppor... |
| GET | `/api/v1/tasks/{taskId}/regretFlags` | GetTaskRegretFlags retrieves all regret flags for a task
Returns the complete... |
| POST | `/api/v1/tasks` | CreateTask creates a fully specified task
Allows creating a task with all fie... |
| POST | `/api/v1/tasks/{taskId}/dependencies` | AddDependency adds a dependency relationship between tasks
Checks for circula... |
| POST | `/api/v1/tasks/{taskId}:flagRegret` | FlagTaskWithRegret flags a task with an implementation regret
Moves the task ... |
| POST | `/api/v1/tasks/{taskId}:move` | MoveTask changes a task's status (e.g., from todo to in-progress)
Enforces wo... |
| POST | `/api/v1/tasks/{taskId}:resolveRegret` | ResolveTaskRegret resolves a task's implementation regret
Moves the task back... |
| POST | `/api/v1/tasks/{taskId}:validateTransition` | ValidateTaskTransition checks if a state transition is valid without performi... |
| PUT | `/api/v1/tasks/{task.id}` | UpdateTask updates a task's content
Modifies task fields but doesn't change s... |
| PUT | `/api/v1/tasks/{taskId}/confidence` | UpdateConfidenceScore updates a task's confidence score
Includes rationale fo... |
| DELETE | `/api/v1/tasks/{taskId}` | DeleteTask removes a task (rarely used, primarily for administrative purposes... |
| DELETE | `/api/v1/tasks/{taskId}/dependencies/{dependencyId}` | RemoveDependency removes a dependency relationship between tasks |

#### GET /api/v1/tasks

Retrieve a paginated list of tasks with optional filtering by status, persona, version, and confidence score

---

#### GET /api/v1/tasks/{taskId}

Retrieves a specific task by its ID, including all details and activity log

---

### Tasks:batchGet

Tasks:batchGet resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/tasks:batchGet` | BatchGetTasks retrieves multiple tasks in a single request
Accepts up to 100 ... |

### Tasks:pluckNext

Tasks:pluckNext resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/tasks:pluckNext` | PluckNextTask finds the best next task to work on
Uses the task selection alg... |

### Tasks:scaffold

Tasks:scaffold resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/tasks:scaffold` | Create a task from template |

#### POST /api/v1/tasks:scaffold

Creates a new task with default values and places it in 'todo' status. This is a quick way to create tasks with minimal input.

---

### Understanding

Understanding resource operations

**Operations**: 3

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/projects/{projectId}/understanding` | List all file understandings for a project |
| GET | `/api/v1/projects/{projectId}/understanding/files` | Get understanding for a specific file |
| GET | `/api/v1/projects/{projectId}/understanding/status` | Get current understanding status for a project |

### Understanding:analyze

Understanding:analyze resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| POST | `/api/v1/projects/{projectId}/understanding:analyze` | Analyze project files to create AI-powered summaries with streaming progress |

### Watch

Watch resource operations

**Operations**: 1

| Method | Path | Summary |
|--------|------|----------|
| GET | `/api/v1/health/watch` | Watch streams health status changes |

## Resource Relationships

This section shows how resources relate to each other.

### Relationship Diagram

```mermaid
graph TD
    releases
    info
    :repair
    :validateTransition
    audits
    tasks:scaffold
    dashboard
    confidence
    taskBranches
    metadata
    results
    personas
    release
    simulations
    :resolveRegret
    containers:publish
    blockedBy
    watch
    containers
    understanding
    hooks:postMerge
    dependencyGraph
    hooks:postCheckout
    answers
    docker
    agents
    task
    :understand
    dependencyTree
    stats
    dependencies
    tasks:batchGet
    tasks
    activities:recent
    projectTypes
    :flagRegret
    progress
    :move
    activities
    features
    regretFlags
    analysis
    documents:analyze
    git
    pullRequests
    releaseBranches
    cycles
    history
    blockers
    files
    containers:validate
    context
    research
    projects:scaffold
    understanding:analyze
    documents
    :validate
    requirements
    :regenerateIndices
    projects
    process
    tasks:pluckNext
    questions
    containers:build
    releases -->|has_many| progress
    :validateTransition -.->|references| task
    audits -->|has_many| status
    audits -->|belongs_to| projects
    audits -->|has_many| process
    audits -->|has_many| results
    dashboard -->|belongs_to| projects
    confidence -->|belongs_to| tasks
    confidence -.->|references| task
    taskBranches -.->|references| task
    results -->|belongs_to| research
    results -->|belongs_to| requirements
    results -->|belongs_to| audits
    results -.->|references| research
    simulations -->|has_many| history
    :resolveRegret -.->|references| task
    blockedBy -->|belongs_to| tasks
    blockedBy -.->|references| task
    understanding -->|belongs_to| projects
    dependencyGraph -->|belongs_to| tasks
    dependencyGraph -.->|references| task
    answers -->|belongs_to| projects
    agents -->|has_many| context
    task -.->|has_many| dependencies
    task -.->|has_many| blockers
    task -.->|has_many| regretFlags
    task -.->|has_many| :flagRegret
    task -.->|has_many| dependencyGraph
    task -.->|has_many| confidence
    task -.->|has_many| dependencyTree
    task -.->|has_many| tasks
    task -.->|has_many| :resolveRegret
    task -.->|has_many| :move
    task -.->|has_many| :validateTransition
    task -.->|has_many| git
    task -.->|has_many| taskBranches
    task -.->|has_many| blockedBy
    dependencyTree -->|belongs_to| tasks
    dependencyTree -.->|references| task
    stats -->|belongs_to| projects
    dependencies -->|belongs_to| tasks
    dependencies -.->|references| task
    tasks -->|has_many| confidence
    tasks -->|has_many| dependencyTree
    tasks -->|has_many| blockers
    tasks -->|has_many| regretFlags
    tasks -->|belongs_to| features
    tasks -->|has_many| dependencies
    tasks -->|has_many| blockedBy
    tasks -->|has_many| dependencyGraph
    tasks -.->|references| task
    projectTypes -->|has_many| questions
    :flagRegret -.->|references| task
    progress -->|belongs_to| releases
    :move -.->|references| task
    features -->|has_many| tasks
    regretFlags -->|belongs_to| tasks
    regretFlags -.->|references| task
    analysis -->|has_many| status
    analysis -.->|has_many| analysis
    analysis -.->|references| analysis
    analysis -.->|has_many| documents
    git -.->|references| task
    history -->|belongs_to| simulations
    blockers -->|belongs_to| tasks
    blockers -.->|references| task
    context -->|belongs_to| agents
    research -->|belongs_to| projects
    research -->|has_many| results
    research -->|has_many| status
    research -->|has_many| process
    research -.->|has_many| research
    research -.->|references| research
    research -.->|has_many| results
    research -.->|has_many| process
    understanding:analyze -->|belongs_to| projects
    documents -.->|references| analysis
    requirements -->|has_many| status
    requirements -->|has_many| process
    requirements -->|has_many| results
    requirements -->|belongs_to| projects
    projects -->|has_many| understanding
    projects -->|has_many| stats
    projects -->|has_many| answers
    projects -->|has_many| research
    projects -->|has_many| audits
    projects -->|has_many| dashboard
    projects -->|has_many| understanding:analyze
    projects -->|has_many| requirements
    process -->|belongs_to| requirements
    process -->|belongs_to| research
    process -->|belongs_to| audits
    process -.->|references| research
    questions -->|belongs_to| projectTypes
```

### Relationship Details

### Releases Relationships

- **has_many** progress (strong strength via `path hierarchy`)
  - releases contains multiple progress resources

### :validateTransition Relationships

- **references** task (medium strength via `taskId`)
  - :validateTransition references a task resource

### Audits Relationships

- **has_many** status (strong strength via `path hierarchy`)
  - audits contains multiple status resources
- **belongs_to** projects (strong strength via `path hierarchy`)
  - audits belongs to a projects resource
- **has_many** process (strong strength via `path hierarchy`)
  - audits contains multiple process resources
- **has_many** results (strong strength via `path hierarchy`)
  - audits contains multiple results resources

### Dashboard Relationships

- **belongs_to** projects (strong strength via `path hierarchy`)
  - dashboard belongs to a projects resource

### Confidence Relationships

- **belongs_to** tasks (strong strength via `path hierarchy`)
  - confidence belongs to a tasks resource
- **references** task (medium strength via `taskId`)
  - confidence references a task resource

### TaskBranches Relationships

- **references** task (medium strength via `taskId`)
  - taskBranches references a task resource

### Results Relationships

- **belongs_to** research (strong strength via `path hierarchy`)
  - results belongs to a research resource
- **belongs_to** requirements (strong strength via `path hierarchy`)
  - results belongs to a requirements resource
- **belongs_to** audits (strong strength via `path hierarchy`)
  - results belongs to a audits resource
- **references** research (medium strength via `researchId`)
  - results references a research resource

### Simulations Relationships

- **has_many** history (strong strength via `path hierarchy`)
  - simulations contains multiple history resources

### :resolveRegret Relationships

- **references** task (medium strength via `taskId`)
  - :resolveRegret references a task resource

### BlockedBy Relationships

- **belongs_to** tasks (strong strength via `path hierarchy`)
  - blockedBy belongs to a tasks resource
- **references** task (medium strength via `taskId`)
  - blockedBy references a task resource

### Understanding Relationships

- **belongs_to** projects (strong strength via `path hierarchy`)
  - understanding belongs to a projects resource

### DependencyGraph Relationships

- **belongs_to** tasks (strong strength via `path hierarchy`)
  - dependencyGraph belongs to a tasks resource
- **references** task (medium strength via `taskId`)
  - dependencyGraph references a task resource

### Answers Relationships

- **belongs_to** projects (strong strength via `path hierarchy`)
  - answers belongs to a projects resource

### Agents Relationships

- **has_many** context (strong strength via `path hierarchy`)
  - agents contains multiple context resources

### Task Relationships

- **has_many** dependencies (medium strength via `taskId`)
  - task contains multiple dependencies resources
- **has_many** blockers (medium strength via `taskId`)
  - task contains multiple blockers resources
- **has_many** regretFlags (medium strength via `taskId`)
  - task contains multiple regretFlags resources
- **has_many** :flagRegret (medium strength via `taskId`)
  - task contains multiple :flagRegret resources
- **has_many** dependencyGraph (medium strength via `taskId`)
  - task contains multiple dependencyGraph resources
- **has_many** confidence (medium strength via `taskId`)
  - task contains multiple confidence resources
- **has_many** dependencyTree (medium strength via `taskId`)
  - task contains multiple dependencyTree resources
- **has_many** tasks (medium strength via `taskId`)
  - task contains multiple tasks resources
- **has_many** :resolveRegret (medium strength via `taskId`)
  - task contains multiple :resolveRegret resources
- **has_many** :move (medium strength via `taskId`)
  - task contains multiple :move resources
- **has_many** :validateTransition (medium strength via `taskId`)
  - task contains multiple :validateTransition resources
- **has_many** git (medium strength via `taskId`)
  - task contains multiple git resources
- **has_many** taskBranches (medium strength via `taskId`)
  - task contains multiple taskBranches resources
- **has_many** blockedBy (medium strength via `taskId`)
  - task contains multiple blockedBy resources

### DependencyTree Relationships

- **belongs_to** tasks (strong strength via `path hierarchy`)
  - dependencyTree belongs to a tasks resource
- **references** task (medium strength via `taskId`)
  - dependencyTree references a task resource

### Stats Relationships

- **belongs_to** projects (strong strength via `path hierarchy`)
  - stats belongs to a projects resource

### Dependencies Relationships

- **belongs_to** tasks (strong strength via `path hierarchy`)
  - dependencies belongs to a tasks resource
- **references** task (medium strength via `taskId`)
  - dependencies references a task resource

### Tasks Relationships

- **has_many** confidence (strong strength via `path hierarchy`)
  - tasks contains multiple confidence resources
- **has_many** dependencyTree (strong strength via `path hierarchy`)
  - tasks contains multiple dependencyTree resources
- **has_many** blockers (strong strength via `path hierarchy`)
  - tasks contains multiple blockers resources
- **has_many** regretFlags (strong strength via `path hierarchy`)
  - tasks contains multiple regretFlags resources
- **belongs_to** features (strong strength via `path hierarchy`)
  - tasks belongs to a features resource
- **has_many** dependencies (strong strength via `path hierarchy`)
  - tasks contains multiple dependencies resources
- **has_many** blockedBy (strong strength via `path hierarchy`)
  - tasks contains multiple blockedBy resources
- **has_many** dependencyGraph (strong strength via `path hierarchy`)
  - tasks contains multiple dependencyGraph resources
- **references** task (medium strength via `taskId`)
  - tasks references a task resource

### ProjectTypes Relationships

- **has_many** questions (strong strength via `path hierarchy`)
  - projectTypes contains multiple questions resources

### :flagRegret Relationships

- **references** task (medium strength via `taskId`)
  - :flagRegret references a task resource

### Progress Relationships

- **belongs_to** releases (strong strength via `path hierarchy`)
  - progress belongs to a releases resource

### :move Relationships

- **references** task (medium strength via `taskId`)
  - :move references a task resource

### Features Relationships

- **has_many** tasks (strong strength via `path hierarchy`)
  - features contains multiple tasks resources

### RegretFlags Relationships

- **belongs_to** tasks (strong strength via `path hierarchy`)
  - regretFlags belongs to a tasks resource
- **references** task (medium strength via `taskId`)
  - regretFlags references a task resource

### Analysis Relationships

- **has_many** status (strong strength via `path hierarchy`)
  - analysis contains multiple status resources
- **has_many** analysis (medium strength via `analysisId`)
  - analysis contains multiple analysis resources
- **references** analysis (medium strength via `analysisId`)
  - analysis references a analysis resource
- **has_many** documents (medium strength via `analysisId`)
  - analysis contains multiple documents resources

### Git Relationships

- **references** task (medium strength via `taskId`)
  - git references a task resource

### History Relationships

- **belongs_to** simulations (strong strength via `path hierarchy`)
  - history belongs to a simulations resource

### Blockers Relationships

- **belongs_to** tasks (strong strength via `path hierarchy`)
  - blockers belongs to a tasks resource
- **references** task (medium strength via `taskId`)
  - blockers references a task resource

### Context Relationships

- **belongs_to** agents (strong strength via `path hierarchy`)
  - context belongs to a agents resource

### Research Relationships

- **belongs_to** projects (strong strength via `path hierarchy`)
  - research belongs to a projects resource
- **has_many** results (strong strength via `path hierarchy`)
  - research contains multiple results resources
- **has_many** status (strong strength via `path hierarchy`)
  - research contains multiple status resources
- **has_many** process (strong strength via `path hierarchy`)
  - research contains multiple process resources
- **has_many** research (medium strength via `researchId`)
  - research contains multiple research resources
- **references** research (medium strength via `researchId`)
  - research references a research resource
- **has_many** results (medium strength via `researchId`)
  - research contains multiple results resources
- **has_many** process (medium strength via `researchId`)
  - research contains multiple process resources

### Understanding:analyze Relationships

- **belongs_to** projects (strong strength via `path hierarchy`)
  - understanding:analyze belongs to a projects resource

### Documents Relationships

- **references** analysis (medium strength via `analysisId`)
  - documents references a analysis resource

### Requirements Relationships

- **has_many** status (strong strength via `path hierarchy`)
  - requirements contains multiple status resources
- **has_many** process (strong strength via `path hierarchy`)
  - requirements contains multiple process resources
- **has_many** results (strong strength via `path hierarchy`)
  - requirements contains multiple results resources
- **belongs_to** projects (strong strength via `path hierarchy`)
  - requirements belongs to a projects resource

### Projects Relationships

- **has_many** understanding (strong strength via `path hierarchy`)
  - projects contains multiple understanding resources
- **has_many** stats (strong strength via `path hierarchy`)
  - projects contains multiple stats resources
- **has_many** answers (strong strength via `path hierarchy`)
  - projects contains multiple answers resources
- **has_many** research (strong strength via `path hierarchy`)
  - projects contains multiple research resources
- **has_many** audits (strong strength via `path hierarchy`)
  - projects contains multiple audits resources
- **has_many** dashboard (strong strength via `path hierarchy`)
  - projects contains multiple dashboard resources
- **has_many** understanding:analyze (strong strength via `path hierarchy`)
  - projects contains multiple understanding:analyze resources
- **has_many** requirements (strong strength via `path hierarchy`)
  - projects contains multiple requirements resources

### Process Relationships

- **belongs_to** requirements (strong strength via `path hierarchy`)
  - process belongs to a requirements resource
- **belongs_to** research (strong strength via `path hierarchy`)
  - process belongs to a research resource
- **belongs_to** audits (strong strength via `path hierarchy`)
  - process belongs to a audits resource
- **references** research (medium strength via `researchId`)
  - process references a research resource

### Questions Relationships

- **belongs_to** projectTypes (strong strength via `path hierarchy`)
  - questions belongs to a projectTypes resource

## Detected Patterns

### Versioning

**Confidence**: high  
**Impact**: Clients should be aware of API version compatibility

API uses URL path versioning. Versions found: v1

**Examples**:
- /api/v1/tasks/{taskId}

### Batch_operations

**Confidence**: low  
**Impact**: Clients can perform bulk operations for better performance

API supports batch operations for bulk create/update/delete

**Examples**:
- /api/v1/tasks:batchGet

