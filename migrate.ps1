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

# Ensure migrations directory exists
if (-not (Test-Path -Path $MIGRATIONS_PATH)) {
    New-Item -ItemType Directory -Path $MIGRATIONS_PATH -Force | Out-Null
    Write-Host "Created migrations directory at $MIGRATIONS_PATH"
}

function Invoke-CreateMigration {
    param([string]$name)
    if ([string]::IsNullOrEmpty($name)) {
        Write-Host "Error: Migration name is required" -ForegroundColor Red
        Show-Help
        exit 1
    }
    $timestamp = Get-Date -Format "yyyyMMddHHmmss"
    $fileName = "${timestamp}_${name}.sql"
    $upPath = Join-Path $MIGRATIONS_PATH "..\migrations\$($timestamp)_${name}.up.sql"
    $downPath = Join-Path $MIGRATIONS_PATH "..\migrations\$($timestamp)_${name}.down.sql"
    
    # Create empty migration files
    Set-Content -Path $upPath -Value "-- Write your UP migration SQL here"
    Set-Content -Path $downPath -Value "-- Write your DOWN migration SQL here"
    
    Write-Host "Created migration files:" -ForegroundColor Green
    Write-Host "  UP:   $upPath"
    Write-Host "  DOWN: $downPath"
}

function Invoke-MigrateUp {
    param([string]$steps = "1")
    $stepParam = if ($steps -eq "all") { "" } else { "-step $steps" }
    migrate -path=$MIGRATIONS_PATH -database=$DB_URL up $stepParam
}

function Invoke-MigrateDown {
    param([string]$steps = "1")
    $stepParam = if ($steps -eq "all") { "" } else { "-step $steps" }
    migrate -path=$MIGRATIONS_PATH -database=$DB_URL down $stepParam
}

function Get-CurrentVersion {
    migrate -path=$MIGRATIONS_PATH -database=$DB_URL version
}

# Main script execution
switch ($action.ToLower()) {
    "create" { Invoke-CreateMigration -name $name }
    "up" { Invoke-MigrateUp -steps $steps }
    "down" { Invoke-MigrateDown -steps $steps }
    "version" { Get-CurrentVersion }
    "seed" { Invoke-Seed }
    "help" { Show-Help }
    default {
        Write-Host "Unknown action: $action" -ForegroundColor Red
        Show-Help
        exit 1
    }
}