# =================================================================================
# Blog System Deployment Script (PowerShell)
# Function: Interactive Config -> Generate .env -> Deploy Database -> Build Frontend/Backend
# =================================================================================

# Set output encoding to UTF-8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

# Script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir

# Configuration file paths
$EnvFile = Join-Path $ProjectRoot ".env"
$SqlTemplate = Join-Path $ProjectRoot "backend\manifest\sql\deploy.sql.template"
$SqlOutput = Join-Path $ProjectRoot "backend\manifest\sql\deploy.sql"
$BackendDir = Join-Path $ProjectRoot "backend"
$FrontendDir = Join-Path $ProjectRoot "fronted"

# Color functions
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function Write-Question {
    param([string]$Message)
    Write-Host "[?] $Message" -ForegroundColor Cyan
}

# Read user input with default value
function Read-Input {
    param(
        [string]$Prompt,
        [string]$Default = "",
        [switch]$Password = $false
    )
    
    if ($Default) {
        $fullPrompt = "${Prompt} [${Default}]: "
    } else {
        $fullPrompt = "${Prompt}: "
    }
    
    Write-Question $fullPrompt
    
    if ($Password) {
        $secureInput = Read-Host -AsSecureString
        $input = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
            [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secureInput)
        )
    } else {
        $input = Read-Host
    }
    
    if ([string]::IsNullOrWhiteSpace($input) -and $Default) {
        return $Default
    }
    
    return $input
}

# Check requirements
function Test-Requirements {
    Write-Info "Checking required tools..."
    
    $missingTools = @()
    
    if (-not (Get-Command mysql -ErrorAction SilentlyContinue)) {
        $missingTools += "mysql"
    }
    
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        Write-Warn "Go not found, backend build will be skipped"
        $script:SkipBackendBuild = $true
    }
    
    if (-not (Get-Command npm -ErrorAction SilentlyContinue) -and 
        -not (Get-Command yarn -ErrorAction SilentlyContinue)) {
        Write-Warn "npm/yarn not found, frontend build will be skipped"
        $script:SkipFrontendBuild = $true
    }
    
    if ($missingTools.Count -gt 0) {
        Write-Error "Missing required tools: $($missingTools -join ', ')"
        Write-Info "Installation instructions:"
        Write-Info "  MySQL: https://dev.mysql.com/downloads/"
        Write-Info "  Go: https://golang.org/dl/"
        Write-Info "  Node.js: https://nodejs.org/"
        exit 1
    }
    
    Write-Info "Requirements check completed"
}

# Interactive configuration
function Start-InteractiveConfig {
    Write-Host ""
    Write-Info "=========================================="
    Write-Info "Interactive Configuration"
    Write-Info "=========================================="
    Write-Host ""
    
    # Database configuration
    Write-Info "--- Database Configuration ---"
    $script:DB_NAME = Read-Input "Database Name" "go_blog"
    $script:DB_USER = Read-Input "Database User" "root"
    $script:DB_PASSWORD = Read-Input "Database Password" -Password
    if ([string]::IsNullOrWhiteSpace($script:DB_PASSWORD)) {
        Write-Error "Database password cannot be empty"
        exit 1
    }
    $script:DB_HOST = Read-Input "Database Host" "127.0.0.1"
    $script:DB_PORT = Read-Input "Database Port" "3306"
    Write-Host ""
    
    # Admin account configuration
    Write-Info "--- Admin Account Configuration ---"
    $script:ADMIN_USERNAME = Read-Input "Admin Username" "admin"
    $script:ADMIN_PASSWORD = Read-Input "Admin Password" -Password
    if ([string]::IsNullOrWhiteSpace($script:ADMIN_PASSWORD)) {
        Write-Error "Admin password cannot be empty"
        exit 1
    }
    $script:ADMIN_NICKNAME = Read-Input "Admin Nickname" "Admin"
    $script:ADMIN_EMAIL = Read-Input "Admin Email" "admin@example.com"
    Write-Host ""
    
    # Website configuration
    Write-Info "--- Website Configuration ---"
    $script:SITE_NAME = Read-Input "Site Name" "Blog System"
    $script:SITE_TITLE = Read-Input "Site Title" "Blog System"
    $script:SITE_SUBTITLE = Read-Input "Site Subtitle" "Share programming insights"
    $script:SITE_DESC = Read-Input "Site Description" "Share programming insights and technical growth"
    $script:SITE_KEYWORDS = Read-Input "Keywords" "Go, Vue, Cloud Native"
    $script:SITE_COPYRIGHT = Read-Input "Copyright" "© 2024 Blog System"
    Write-Host ""
    
    # Production environment configuration
    Write-Info "--- Production Environment Configuration ---"
    $script:FRONTEND_DOMAIN = Read-Input "Frontend Domain (e.g., https://blog.example.com, leave empty to skip)" ""
    $script:BACKEND_DOMAIN = Read-Input "Backend API Domain (e.g., https://api.example.com, leave empty to skip)" ""
    $script:DEPLOY_ENV = Read-Input "Deploy Environment" "development"
    Write-Host ""
    
    # JWT configuration
    Write-Info "--- JWT Configuration ---"
    $jwtSecret = Read-Input "JWT Secret (leave empty to auto-generate)" ""
    if ([string]::IsNullOrWhiteSpace($jwtSecret)) {
        # Generate random secret
        $bytes = New-Object byte[] 32
        [System.Security.Cryptography.RNGCryptoServiceProvider]::Create().GetBytes($bytes)
        $script:JWT_SECRET = [Convert]::ToBase64String($bytes)
    } else {
        $script:JWT_SECRET = $jwtSecret
    }
    Write-Host ""
    
    # OAuth configuration (optional)
    Write-Info "--- OAuth Configuration (Optional, press Enter to skip) ---"
    $script:GITHUB_CLIENT_ID = Read-Input "GitHub Client ID" ""
    $script:GITHUB_CLIENT_SECRET = Read-Input "GitHub Client Secret" ""
    $script:GOOGLE_CLIENT_ID = Read-Input "Google Client ID" ""
    $script:GOOGLE_CLIENT_SECRET = Read-Input "Google Client Secret" ""
    $script:WECHAT_APP_ID = Read-Input "WeChat AppID" ""
    $script:WECHAT_APP_SECRET = Read-Input "WeChat AppSecret" ""
    Write-Host ""
    
    # Redis configuration
    Write-Info "--- Redis Configuration (Optional) ---"
    $script:REDIS_HOST = Read-Input "Redis Host" "127.0.0.1"
    $script:REDIS_PORT = Read-Input "Redis Port" "6379"
    $script:REDIS_DB = Read-Input "Redis Database" "1"
    Write-Host ""
    
    # Set OAuth frontend URL
    if ($script:FRONTEND_DOMAIN) {
        $script:OAUTH_FRONTEND_URL = $script:FRONTEND_DOMAIN
    } else {
        $script:OAUTH_FRONTEND_URL = "http://localhost:5173"
    }
}

# Generate .env file
function New-EnvFile {
    Write-Info "Generating .env file..."
    
    $envContent = @"
# =================================================================================
# Blog System Environment Variables Configuration
# Auto-generated: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')
# =================================================================================

# ============================================
# Database Configuration
# ============================================
DB_NAME=$script:DB_NAME
DB_USER=$script:DB_USER
DB_PASSWORD=$script:DB_PASSWORD
DB_HOST=$script:DB_HOST
DB_PORT=$script:DB_PORT

# ============================================
# Admin Account Configuration
# ============================================
ADMIN_USERNAME=$script:ADMIN_USERNAME
ADMIN_PASSWORD=$script:ADMIN_PASSWORD
ADMIN_NICKNAME=$script:ADMIN_NICKNAME
ADMIN_EMAIL=$script:ADMIN_EMAIL

# ============================================
# Website Configuration
# ============================================
SITE_NAME=$script:SITE_NAME
SITE_TITLE=$script:SITE_TITLE
SITE_SUBTITLE=$script:SITE_SUBTITLE
SITE_DESC=$script:SITE_DESC
SITE_KEYWORDS=$script:SITE_KEYWORDS
SITE_COPYRIGHT=$script:SITE_COPYRIGHT

# ============================================
# Production Environment Configuration
# ============================================
FRONTEND_DOMAIN=$script:FRONTEND_DOMAIN
BACKEND_DOMAIN=$script:BACKEND_DOMAIN
DEPLOY_ENV=$script:DEPLOY_ENV

# ============================================
# JWT Configuration
# ============================================
JWT_SECRET=$script:JWT_SECRET

# ============================================
# OAuth Configuration (Optional)
# ============================================
OAUTH_FRONTEND_URL=$script:OAUTH_FRONTEND_URL
GITHUB_CLIENT_ID=$script:GITHUB_CLIENT_ID
GITHUB_CLIENT_SECRET=$script:GITHUB_CLIENT_SECRET
GOOGLE_CLIENT_ID=$script:GOOGLE_CLIENT_ID
GOOGLE_CLIENT_SECRET=$script:GOOGLE_CLIENT_SECRET
WECHAT_APP_ID=$script:WECHAT_APP_ID
WECHAT_APP_SECRET=$script:WECHAT_APP_SECRET

# ============================================
# Redis Configuration (Optional)
# ============================================
REDIS_HOST=$script:REDIS_HOST
REDIS_PORT=$script:REDIS_PORT
REDIS_DB=$script:REDIS_DB
"@
    
    $envContent | Out-File -FilePath $EnvFile -Encoding UTF8 -NoNewline
    Write-Info ".env file generated: $EnvFile"
}

# Load environment variables
function Import-EnvFile {
    if (-not (Test-Path $EnvFile)) {
        Write-Error ".env file not found: $EnvFile"
        exit 1
    }
    
    Get-Content $EnvFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            Set-Variable -Name $name -Value $value -Scope Script
        }
    }
    
    # Set defaults
    if (-not $script:DB_HOST) { $script:DB_HOST = "127.0.0.1" }
    if (-not $script:DB_PORT) { $script:DB_PORT = "3306" }
    if (-not $script:ADMIN_NICKNAME) { $script:ADMIN_NICKNAME = "Admin" }
    if (-not $script:ADMIN_EMAIL) { $script:ADMIN_EMAIL = "admin@example.com" }
    if (-not $script:SITE_NAME) { $script:SITE_NAME = "Blog System" }
    if (-not $script:SITE_TITLE) { $script:SITE_TITLE = "Blog System" }
    if (-not $script:SITE_SUBTITLE) { $script:SITE_SUBTITLE = "Share programming insights" }
    if (-not $script:SITE_DESC) { $script:SITE_DESC = "Share programming insights and technical growth" }
    if (-not $script:SITE_KEYWORDS) { $script:SITE_KEYWORDS = "Go, Vue, Cloud Native" }
    if (-not $script:SITE_COPYRIGHT) { $script:SITE_COPYRIGHT = "© 2024 Blog System" }
    if (-not $script:DEPLOY_ENV) { $script:DEPLOY_ENV = "development" }
    if (-not $script:OAUTH_FRONTEND_URL) { $script:OAUTH_FRONTEND_URL = "http://localhost:5173" }
    if (-not $script:REDIS_HOST) { $script:REDIS_HOST = "127.0.0.1" }
    if (-not $script:REDIS_PORT) { $script:REDIS_PORT = "6379" }
    if (-not $script:REDIS_DB) { $script:REDIS_DB = "1" }
}

# Generate SQL file
function New-SqlFile {
    Write-Info "Generating deployment SQL file..."
    
    if (-not (Test-Path $SqlTemplate)) {
        Write-Error "SQL template file not found: $SqlTemplate"
        exit 1
    }
    
    $templateContent = Get-Content $SqlTemplate -Raw -Encoding UTF8
    
    # Replace environment variables
    $templateContent = $templateContent -replace '\$\{DB_NAME\}', $script:DB_NAME
    $templateContent = $templateContent -replace '\$\{ADMIN_USERNAME\}', $script:ADMIN_USERNAME
    $templateContent = $templateContent -replace '\$\{ADMIN_PASSWORD\}', $script:ADMIN_PASSWORD
    $templateContent = $templateContent -replace '\$\{ADMIN_NICKNAME\}', $script:ADMIN_NICKNAME
    $templateContent = $templateContent -replace '\$\{ADMIN_EMAIL\}', $script:ADMIN_EMAIL
    $templateContent = $templateContent -replace '\$\{SITE_NAME\}', $script:SITE_NAME
    $templateContent = $templateContent -replace '\$\{SITE_TITLE\}', $script:SITE_TITLE
    $templateContent = $templateContent -replace '\$\{SITE_SUBTITLE\}', $script:SITE_SUBTITLE
    $templateContent = $templateContent -replace '\$\{SITE_DESC\}', $script:SITE_DESC
    $templateContent = $templateContent -replace '\$\{SITE_KEYWORDS\}', $script:SITE_KEYWORDS
    $templateContent = $templateContent -replace '\$\{SITE_COPYRIGHT\}', $script:SITE_COPYRIGHT
    
    $templateContent | Out-File -FilePath $SqlOutput -Encoding UTF8 -NoNewline
    Write-Info "SQL file generated: $SqlOutput"
}

# Deploy database
function Deploy-Database {
    Write-Info "Starting database deployment..."
    
    # Test database connection
    Write-Info "Testing database connection..."
    $mysqlCmd = "mysql -h$($script:DB_HOST) -P$($script:DB_PORT) -u$($script:DB_USER) -p$($script:DB_PASSWORD) -e `"SELECT 1`""
    $result = cmd /c $mysqlCmd 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Database connection failed, please check configuration"
        exit 1
    }
    Write-Info "Database connection successful"
    
    # Execute SQL script
    Write-Info "Executing SQL script..."
    $mysqlCmd = "mysql -h$($script:DB_HOST) -P$($script:DB_PORT) -u$($script:DB_USER) -p$($script:DB_PASSWORD) < `"$SqlOutput`""
    $result = cmd /c $mysqlCmd 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Database deployment failed"
        exit 1
    }
    Write-Info "Database deployed successfully!"
    
    # Verify deployment
    Write-Info "Verifying deployment..."
    $mysqlCmd = "mysql -h$($script:DB_HOST) -P$($script:DB_PORT) -u$($script:DB_USER) -p$($script:DB_PASSWORD) -D$($script:DB_NAME) -e `"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='$($script:DB_NAME)'`""
    $tableCount = (cmd /c $mysqlCmd 2>&1 | Select-String -Pattern '\d+').Matches.Value
    if ($tableCount) {
        Write-Info "Database tables created successfully, total: $tableCount tables"
    } else {
        Write-Warn "Could not detect database tables, please check SQL execution result"
    }
}

# Build backend
function Build-Backend {
    Write-Info "Starting backend build..."
    
    Push-Location $BackendDir
    
    if (-not (Test-Path "go.mod")) {
        Write-Error "go.mod file not found, please check backend directory"
        exit 1
    }
    
    Write-Info "Downloading Go dependencies..."
    go mod download
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to download Go dependencies"
        exit 1
    }
    
    Write-Info "Building backend program..."
    go build -o blog.exe -ldflags "-s -w" main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Backend build failed"
        exit 1
    }
    
    Write-Info "Backend build successful: $BackendDir\blog.exe"
    Pop-Location
}

# Build frontend
function Build-Frontend {
    Write-Info "Starting frontend build..."
    
    Push-Location $FrontendDir
    
    if (-not (Test-Path "package.json")) {
        Write-Error "package.json file not found, please check frontend directory"
        exit 1
    }
    
    Write-Info "Installing frontend dependencies..."
    npm install
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to install frontend dependencies"
        exit 1
    }
    
    # Create production environment configuration
    if ($script:FRONTEND_DOMAIN -and $script:BACKEND_DOMAIN) {
        Write-Info "Creating production environment configuration..."
        $prodEnvContent = @"
VITE_API_BASE_URL=$($script:BACKEND_DOMAIN)/api
VITE_FRONTEND_URL=$($script:FRONTEND_DOMAIN)
"@
        $prodEnvContent | Out-File -FilePath ".env.production" -Encoding UTF8 -NoNewline
    }
    
    Write-Info "Building frontend application..."
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Frontend build failed"
        exit 1
    }
    
    if (Test-Path "dist") {
        Write-Info "Frontend build successful: $FrontendDir\dist"
    } else {
        Write-Error "Frontend build failed, dist directory not found"
        exit 1
    }
    
    Pop-Location
}

# Show deployment info
function Show-DeploymentInfo {
    Write-Host ""
    Write-Info "=========================================="
    Write-Info "Deployment Completed!"
    Write-Info "=========================================="
    Write-Host ""
    Write-Info "Database Information:"
    Write-Host "  Database Name: $($script:DB_NAME)"
    Write-Host "  Host: $($script:DB_HOST):$($script:DB_PORT)"
    Write-Host "  User: $($script:DB_USER)"
    Write-Host ""
    Write-Info "Admin Account:"
    Write-Host "  Username: $($script:ADMIN_USERNAME)"
    Write-Host "  Password: $($script:ADMIN_PASSWORD)"
    Write-Host ""
    if ($script:FRONTEND_DOMAIN) {
        Write-Info "Production Environment Configuration:"
        Write-Host "  Frontend Domain: $($script:FRONTEND_DOMAIN)"
        Write-Host "  Backend Domain: $($script:BACKEND_DOMAIN)"
        Write-Host "  Frontend Build Directory: $FrontendDir\dist"
        Write-Host "  Backend Executable: $BackendDir\blog.exe"
        Write-Host ""
    }
    Write-Warn "Please keep the admin password and JWT secret safe!"
    Write-Host ""
    Write-Info "Next Steps:"
    Write-Host "  1. Configure web server (Nginx/IIS) to point to frontend dist directory"
    Write-Host "  2. Configure reverse proxy to forward API requests to backend"
    Write-Host "  3. Start backend service: cd backend && .\blog.exe"
    Write-Host ""
}

# Main function
function Main {
    Write-Host ""
    Write-Info "=========================================="
    Write-Info "Blog System Deployment Script"
    Write-Info "=========================================="
    Write-Host ""
    
    # Check if .env file exists
    if (Test-Path $EnvFile) {
        Write-Warn ".env file already exists"
        $rebuild = Read-Input "Reconfigure? (y/N)" "N"
        if ($rebuild -ne "y" -and $rebuild -ne "Y") {
            Write-Info "Using existing configuration..."
            Import-EnvFile
        } else {
            Start-InteractiveConfig
            New-EnvFile
        }
    } else {
        Start-InteractiveConfig
        New-EnvFile
    }
    
    # Load environment variables
    Import-EnvFile
    
    # Deployment steps
    Test-Requirements
    New-SqlFile
    Deploy-Database
    
    # Ask about building
    Write-Host ""
    if (-not $script:SkipBackendBuild) {
        $buildBackend = Read-Input "Build backend? (Y/n)" "Y"
        if ($buildBackend -ne "n" -and $buildBackend -ne "N") {
            Build-Backend
        }
    }
    
    if (-not $script:SkipFrontendBuild) {
        $buildFrontend = Read-Input "Build frontend? (Y/n)" "Y"
        if ($buildFrontend -ne "n" -and $buildFrontend -ne "N") {
            Build-Frontend
        }
    }
    
    Show-DeploymentInfo
}

# Execute main function
Main

