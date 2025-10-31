param(
    [string]$action = "help",
    [string]$name = "",
    [string]$steps = "1"
)

# Configuration
$MIGRATIONS_PATH = "cmd/migrate/migrations"
$DB_URL = "postgres://postgres:12345@localhost:5432/social?sslmode=disable"

function Show-Help {
    Write-Host "Social App Migration Tool"
    Write-Host "========================="
    Write-Host "Usage: .\migrate.ps1 [action] [migration_name] [steps]"
    Write-Host ""
    Write-Host "Actions:"
    Write-Host "  create [name]  - Create new migration files"
    Write-Host "  up [steps]     - Apply migrations (default: 1, use 'all' for all)"
    Write-Host "  down [steps]   - Rollback migrations (default: 1, use 'all' for all)"
    Write-Host "  seed           - Seed the database with test data"
    Write-Host "  version        - Show current migration version"
    Write-Host "  help           - Show this help"
}

function Invoke-Seed {
    Write-Host "Seeding database..." -ForegroundColor Cyan
    go run cmd/migrate/seed/main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error seeding database" -ForegroundColor Red
        exit 1
    }
    Write-Host "Database seeded successfully" -ForegroundColor Green
}

function Invoke-Migrate {
    param(
        [string]$direction,
        [string]$steps
    )

    # Ensure migrations directory exists
    if (-not (Test-Path -Path $MIGRATIONS_PATH)) {
        New-Item -ItemType Directory -Path $MIGRATIONS_PATH -Force | Out-Null
        Write-Host "Created migrations directory at $MIGRATIONS_PATH"
    }

    $migrateCmd = "migrate"
    if (-not (Get-Command $migrateCmd -ErrorAction SilentlyContinue)) {
        $migrateCmd = "$env:GOPATH\bin\migrate"
    }

    # Get absolute path and normalize slashes
    $absPath = (Resolve-Path $MIGRATIONS_PATH -ErrorAction SilentlyContinue)
    if (-not $absPath) {
        $absPath = $MIGRATIONS_PATH
    } else {
        $absPath = $absPath.Path
    }
    $absPath = $absPath -replace '\\', '/'

    # Build the command
    $cmd = "$migrateCmd -path=`"$absPath`" -database=`"$DB_URL`" $direction"
    if ($steps -ne "all") {
        $cmd += " $steps"
    }

    Write-Host "Running: $cmd" -ForegroundColor Cyan
    $ErrorActionPreference = "Stop"
    try {
        Invoke-Expression $cmd
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Error running migration" -ForegroundColor Red
            exit 1
        }
    } catch {
        Write-Host "Migration error: $_" -ForegroundColor Red
        exit 1
    }
}

function Invoke-CreateMigration {
    param([string]$name)
    if ([string]::IsNullOrEmpty($name)) {
        Write-Host "Error: Migration name is required" -ForegroundColor Red
        Show-Help
        exit 1
    }
    $timestamp = Get-Date -Format "yyyyMMddHHmmss"
    $upPath = Join-Path $MIGRATIONS_PATH "${timestamp}_${name}.up.sql"
    $downPath = Join-Path $MIGRATIONS_PATH "${timestamp}_${name}.down.sql"
    
    "# Add your SQL for migrating up" | Out-File -FilePath $upPath -Encoding utf8
    "# Add your SQL for rolling back" | Out-File -FilePath $downPath -Encoding utf8
    
    Write-Host "Created migration files:" -ForegroundColor Green
    Write-Host "  $upPath"
    Write-Host "  $downPath"
}

# Main execution
try {
    switch ($action.ToLower()) {
        "create" { Invoke-CreateMigration $name }
        "up" { Invoke-Migrate "up" $steps }
        "down" { Invoke-Migrate "down" $steps }
        "seed" { Invoke-Seed }
        "version" { 
            $absPath = (Resolve-Path $MIGRATIONS_PATH -ErrorAction SilentlyContinue)
            if (-not $absPath) { $absPath = $MIGRATIONS_PATH }
            $absPath = $absPath -replace '\\', '/'
            migrate -path="$absPath" -database="$DB_URL" version 
        }
        "help" { Show-Help }
        default {
            Write-Host "Unknown action: $action" -ForegroundColor Red
            Show-Help
            exit 1
        }
    }
} catch {
    Write-Host "An error occurred: $_" -ForegroundColor Red
    exit 1
}