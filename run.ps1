param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "`nStocky Backend Management Commands:" -ForegroundColor Cyan
    Write-Host "  .\run.ps1 migrate    - Run database migrations and seed data" -ForegroundColor Green
    Write-Host "  .\run.ps1 run        - Start the server in production mode" -ForegroundColor Green
    Write-Host "  .\run.ps1 dev        - Start the server in development mode" -ForegroundColor Green
    Write-Host "  .\run.ps1 watch      - Start server with hot-reload (requires air)" -ForegroundColor Green
    Write-Host "  .\run.ps1 build      - Build the application" -ForegroundColor Green
    Write-Host "  .\run.ps1 clean      - Clean build artifacts" -ForegroundColor Green
    Write-Host "  .\run.ps1 test       - Run tests" -ForegroundColor Green
    Write-Host ""
}

switch ($Command.ToLower()) {
    "migrate" {
        Write-Host "Running database migrations..." -ForegroundColor Yellow
        go run cmd/migrate/main.go
    }
    "run" {
        Write-Host "Starting server in production mode..." -ForegroundColor Yellow
        go run main.go
    }
    "dev" {
        Write-Host "Starting server in development mode..." -ForegroundColor Yellow
        $env:GIN_MODE = "debug"
        $env:LOG_LEVEL = "debug"
        go run main.go
    }
    "watch" {
        Write-Host "Checking for air installation..." -ForegroundColor Yellow
        $airInstalled = Get-Command air -ErrorAction SilentlyContinue
        if (-not $airInstalled) {
            Write-Host "air is not installed. Installing..." -ForegroundColor Red
            Write-Host "Run: go install github.com/air-verse/air@latest" -ForegroundColor Cyan
            Write-Host "Make sure GOPATH/bin is in your PATH" -ForegroundColor Cyan
            exit 1
        }
        Write-Host "Starting server with hot-reload..." -ForegroundColor Yellow
        $env:GIN_MODE = "debug"
        $env:LOG_LEVEL = "debug"
        air
    }
    "build" {
        Write-Host "Building application..." -ForegroundColor Yellow
        if (!(Test-Path "bin")) {
            New-Item -ItemType Directory -Path "bin" | Out-Null
        }
        go build -o bin/stocky-backend.exe main.go
        go build -o bin/migrate.exe cmd/migrate/main.go
        Write-Host "Build complete! Executables in bin/" -ForegroundColor Green
    }
    "clean" {
        Write-Host "Cleaning build artifacts..." -ForegroundColor Yellow
        if (Test-Path "bin") {
            Remove-Item -Recurse -Force "bin"
        }
        go clean
        Write-Host "Clean complete!" -ForegroundColor Green
    }
    "test" {
        Write-Host "Running tests..." -ForegroundColor Yellow
        go test -v ./...
    }
    default {
        Show-Help
    }
}
