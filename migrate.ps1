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
    Write-Host "                 Use 'up all' to apply all pending migrations"
    Write-Host "  down [steps]   - Rollback migrations (default: 1, use 'all' for all)"
    Write-Host "  seed           - Seed the database with test data"
    Write-Host "  version        - Show current migration version"
    Write-Host "  force [version] - Force mark a specific version as complete"
    Write-Host "  gen-docs       - Generate API documentation using swag"
    Write-Host "  help           - Show this help"
}

function Invoke-GenDocs {
    Write-Host "Generating API documentation with swag..." -ForegroundColor Cyan
    $swagCmd = "swag"
    
    # Check if swag is installed
    if (-not (Get-Command $swagCmd -ErrorAction SilentlyContinue)) {
        Write-Host "swag not found. Installing swag..." -ForegroundColor Yellow
        go install github.com/swaggo/swag/cmd/swag@latest
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Error: Failed to install swag" -ForegroundColor Red
            exit 1
        }
        # Refresh PATH to include Go bin directory
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
        $swagCmd = "$env:GOPATH\bin\swag"
    }

    # Save current directory
    $originalDir = Get-Location
    $projectRoot = $PSScriptRoot
    $docsDir = Join-Path $projectRoot "docs"

    try {
        # Ensure docs directory exists
        if (-not (Test-Path $docsDir)) {
            New-Item -ItemType Directory -Path $docsDir -Force | Out-Null
        }

        # Change to the cmd/api directory where main.go is located
        $apiDir = Join-Path $projectRoot "cmd\api"
        if (-not (Test-Path $apiDir)) {
            Write-Host "Error: API directory not found at $apiDir" -ForegroundColor Red
            exit 1
        }
        
        Set-Location -Path $apiDir
        Write-Host "Changed to API directory: $(Get-Location)" -ForegroundColor Cyan

        # Check if main.go exists
        if (-not (Test-Path "main.go")) {
            Write-Host "Error: main.go not found in $(Get-Location)" -ForegroundColor Red
            exit 1
        }

        # Run swag init with relative paths
        Write-Host "`nRunning: swag init -g main.go -o ../../docs" -ForegroundColor DarkGray
        & $swagCmd init -g "main.go" -o "../../docs"
        
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Error: Failed to generate documentation" -ForegroundColor Red
            exit 1
        }
        
        # Verify files were created
        $files = Get-ChildItem -Path "../../docs" -Recurse -ErrorAction SilentlyContinue
        if ($files.Count -gt 0) {
            Write-Host "`nGenerated files:" -ForegroundColor Green
            $files | ForEach-Object {
                Write-Host "  $($_.FullName.Substring($projectRoot.Length + 1))" -ForegroundColor Cyan
            }
            
            Write-Host "`nAPI documentation generated successfully!" -ForegroundColor Green
            Write-Host "Documentation will be available at: http://localhost:8080/swagger/index.html" -ForegroundColor Cyan
            Write-Host "Make sure your API server is running on port 8080" -ForegroundColor Cyan
        } else {
            Write-Host "Error: No documentation files were generated" -ForegroundColor Red
            Write-Host "Please check that your main.go file has the proper Swagger annotations" -ForegroundColor Yellow
            exit 1
        }
    }
    catch {
        Write-Host "An error occurred: $_" -ForegroundColor Red
        Write-Host $_.ScriptStackTrace -ForegroundColor DarkGray
        exit 1
    }
    finally {
        # Restore the original directory
        Set-Location -Path $originalDir
    }
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

    $migrateArgs = @(
        "-path", $MIGRATIONS_PATH,
        "-database", $DB_URL,
        $direction
    )

    if ($steps -ne "all") {
        $migrateArgs += $steps
    }

    Write-Host "Running migration ($direction $steps)..." -ForegroundColor Cyan
    & $migrateCmd $migrateArgs

    if ($LASTEXITCODE -ne 0) {
        Write-Host "Migration failed" -ForegroundColor Red
        exit 1
    }
}

function Invoke-CreateMigration {
    param(
        [string]$name
    )

    $migrateCmd = "migrate"
    if (-not (Get-Command $migrateCmd -ErrorAction SilentlyContinue)) {
        $migrateCmd = "$env:GOPATH\bin\migrate"
    }

    $migrateArgs = @(
        "create",
        "-ext", "sql",
        "-dir", $MIGRATIONS_PATH,
        "-seq", $name
    )

    Write-Host "Creating new migration: $name" -ForegroundColor Cyan
    & $migrateCmd $migrateArgs

    if ($LASTEXITCODE -ne 0) {
        Write-Host "Failed to create migration" -ForegroundColor Red
        exit 1
    }
}

function Invoke-ForceVersion {
    param(
        [string]$version
    )

    $migrateCmd = "migrate"
    if (-not (Get-Command $migrateCmd -ErrorAction SilentlyContinue)) {
        $migrateCmd = "$env:GOPATH\bin\migrate"
    }

    Write-Host "Forcing version to $version..." -ForegroundColor Yellow
    & $migrateCmd -path $MIGRATIONS_PATH -database $DB_URL force $version

    if ($LASTEXITCODE -ne 0) {
        Write-Host "Failed to force version" -ForegroundColor Red
        exit 1
    }
    Write-Host "Successfully forced version $version" -ForegroundColor Green
}

try {
    switch ($action.ToLower()) {
        "create" {
            if ([string]::IsNullOrEmpty($name)) {
                Write-Host "Error: Migration name is required" -ForegroundColor Red
                Show-Help
                exit 1
            }
            Invoke-CreateMigration $name
        }
        "up" {
            Invoke-Migrate -direction "up" -steps $steps
        }
        "down" {
            Invoke-Migrate -direction "down" -steps $steps
        }
        "seed" {
            Invoke-Seed
        }
        "version" {
            Invoke-Migrate -direction "version" -steps "1"
        }
        "force" {
            if ([string]::IsNullOrEmpty($name)) {
                Write-Host "Error: Version number is required" -ForegroundColor Red
                Show-Help
                exit 1
            }
            Invoke-ForceVersion $name
        }
        "gen-docs" {
            Invoke-GenDocs
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
} catch {
    Write-Host "An error occurred: $_" -ForegroundColor Red
    Write-Host $_.ScriptStackTrace -ForegroundColor DarkGray
    exit 1
}