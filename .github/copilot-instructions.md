This is a Wails (Go, Typescript React, SQLite) based repository. It is primarily responsible for connecting timelines of livestreamers' stream and game accounts to find and clip highlights. Please follow these guidelines when contributing:

## Code Standards

### Required Before Each Commit
- Run `make fmt` before committing any changes to ensure proper code formatting
- This will run gofmt on all Go files to maintain consistent style

### Development Flow
- Build: `make build`
- Integration Tests: `make test-integration`
- Run in Dev Mode: `make run`

## Repository Structure
- `cmd/migrate`: Runs SQLite migrations
- `./main.go`: Wails application entry point
- `internal/`: App logic, some of which is exposed to the react frontend by Wails. Bound objects can be found in `./main.go`, in the `Bind` parameter passed to `wails.Run`. 
- `tests/`: Integration tests.

## Key Guidelines
1. Follow Go best practices and idiomatic patterns
2. Maintain existing code structure and organization
3. Use shadcn UI components in the React frontend
4. In the frontend, use models defined in `frontend/wailsjs/models` as much as possible before creating new types. If a required attribute of a model is making it unusable, consider making the attribute optional in the backend `./internal/models` definition. 