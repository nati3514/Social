param(
    [string]$action = "help",
    [string]$name = "",
    [string]$steps = "1"
)

# Configuration
$MIGRATIONS_PATH = "cmd\migrate\migrations"
$DB_URL = "postgresql://postgres:12345@localhost:5432/social?sslmode=disable"

function Show-Help {
    Write-Host "Social App Migration Tool"
    Write-Host "========================="
    Write-Host "Usage: .\migrate.ps1 [action] [migration_name] [steps]"
    Write-Host ""
    Write-Host "Actions:"
    Write-Host "  create [name]  - Create new migration files"
    Write-Host "  up [steps]     - Apply migrations (default: 1, use 'all' for all)"
    Write-Host "  down [steps]   - Rollback migrations (default: 1, use 'all' for all)"
    Write-Host "  version       - Show current migration version"
    Write-Host "  help          - Show this help"
}

# Ensure migrations directory exists
if (-not (Test-Path -Path $MIGRATIONS_PATH)) {
    New-Item -ItemType Directory -Path $MIGRATIONS_PATH -Force | Out-Null
    Write-Host "Created migrations directory at $MIGRATIONS_PATH"
}

# Execute migration command
switch ($action.ToLower()) {
    "create" {
        if (-not $name) {
            Write-Host "Error: Migration name is required" -ForegroundColor Red
            Show-Help
            exit 1
        }
        Write-Host "Creating new migration: $name" -ForegroundColor Cyan
        $fullPath = Join-Path -Path $PSScriptRoot -ChildPath $MIGRATIONS_PATH
        $unixPath = $fullPath.Replace("\", "/")
        migrate create -ext sql -dir "$unixPath" -seq $name
    }
    "up" {
        Write-Host "Applying migrations..." -ForegroundColor Cyan
        $dbUrl = "postgres://postgres:12345@localhost:5432/social?sslmode=disable"
        $fullPath = Join-Path -Path $PSScriptRoot -ChildPath $MIGRATIONS_PATH
        $unixPath = $fullPath.Replace("\", "/")
        if ($steps -eq "all") {
            migrate -path="$unixPath" -database "$dbUrl" up
        } else {
            migrate -path="$unixPath" -database "$dbUrl" up $steps
        }
    }
    "down" {
        Write-Host "Rolling back migrations..." -ForegroundColor Yellow
        $dbUrl = "postgres://postgres:12345@localhost:5432/social?sslmode=disable"
        $fullPath = Join-Path -Path $PSScriptRoot -ChildPath $MIGRATIONS_PATH
        $unixPath = $fullPath.Replace("\", "/")
        if ($steps -eq "all") {
            migrate -path="$unixPath" -database "$dbUrl" down
        } else {
            migrate -path="$unixPath" -database "$dbUrl" down $steps
        }
    }
    "version" {
        Write-Host "Current migration version:" -ForegroundColor Cyan
        $dbUrl = "postgres://postgres:12345@localhost:5432/social?sslmode=disable"
        $fullPath = Join-Path -Path $PSScriptRoot -ChildPath $MIGRATIONS_PATH
        $unixPath = $fullPath.Replace("\", "/")
        migrate -path="$unixPath" -database "$dbUrl" version
    }
    "help" {
        Show-Help
    }
    default {
        Write-Host "Unknown action: $action" -ForegroundColor Red
        Show-Help
        exit 1
    }
}

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Migration command failed with code $LASTEXITCODE" -ForegroundColor Red
    exit $LASTEXITCODE
}
