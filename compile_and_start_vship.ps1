param(
    [int]$ApiPort = 3002
)

try { taskkill /F /IM vship-api.exe /T 2>$null } catch {}
Start-Sleep -Seconds 2

$ErrorActionPreference = "Stop"

# Get project directory
$ProjectDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ProjectDir

function Stop-ExistingProcessForExe {
    param([string]$ExePath)

    $exeName = [System.IO.Path]::GetFileNameWithoutExtension($ExePath)
    $procs = @()
    try { $procs = Get-Process -Name $exeName -ErrorAction SilentlyContinue } catch { $procs = @() }

    $stoppedAny = $false
    foreach ($p in $procs) {
        $pPath = $null
        try { $pPath = $p.Path } catch { $pPath = $null }

        if (($pPath -eq $null) -or ($pPath -ieq $ExePath)) {
            try {
                Stop-Process -Id $p.Id -Force -ErrorAction SilentlyContinue
                $stoppedAny = $true
            } catch {}
        }
    }

    if ($stoppedAny) {
        Start-Sleep -Milliseconds 400
    }
}

# Create bin directory (if missing)
$BinDir = Join-Path $ProjectDir "bin"
if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir | Out-Null
    Write-Host "[OK] Created bin directory" -ForegroundColor Green
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  vShip Compile & Start" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check Go environment
Write-Host ">>> Checking Go environment..." -ForegroundColor Yellow
$goVersion = go version
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERR] Go is not installed or not in PATH" -ForegroundColor Red
    exit 1
}
Write-Host "[OK] $goVersion" -ForegroundColor Green
Write-Host ""

# Check .env file
if (-not (Test-Path ".env")) {
    Write-Host "[WARN] .env not found; ensure env vars are configured" -ForegroundColor Yellow
    Write-Host ""
}

# Download dependencies
Write-Host ">>> Downloading Go modules..." -ForegroundColor Yellow
go mod download
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERR] go mod download failed" -ForegroundColor Red
    exit 1
}
Write-Host "[OK] Go modules downloaded" -ForegroundColor Green
Write-Host ""

# Build API server
Write-Host ">>> Building API server..." -ForegroundColor Yellow
$apiExe = Join-Path $BinDir "vship-api.exe"
$apiExeNew = "$apiExe.new"
$apiExeOld = "$apiExe.old"

if (Test-Path $apiExe) {
    Stop-ExistingProcessForExe -ExePath $apiExe
    try { taskkill /F /IM vship-api.exe /T 2>$null } catch {}
    Start-Sleep -Seconds 2
}

go build -o $apiExeNew ./cmd/api/main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERR] API build failed during compilation" -ForegroundColor Red
    exit 1
}

if (Test-Path $apiExe) {
    if (Test-Path $apiExeOld) { try { Remove-Item $apiExeOld -Force } catch {} }
    try {
        Rename-Item -Path $apiExe -NewName $apiExeOld -Force -ErrorAction Stop
        Write-Host "   Renamed old executable to .old" -ForegroundColor Gray
    } catch {
        Write-Host "[WARN] Could not rename old executable. It might be locked." -ForegroundColor Yellow
    }
}

try {
    Rename-Item -Path $apiExeNew -NewName $apiExe -Force -ErrorAction Stop
} catch {
    Write-Host "[ERR] Could not put new executable in place. Please check file permissions or antivirus." -ForegroundColor Red
    exit 1
}
Write-Host "[OK] API build succeeded" -ForegroundColor Green
Write-Host ""

# Start API server
Write-Host ">>> Starting API server..." -ForegroundColor Yellow
Write-Host "    Listening on port: $ApiPort" -ForegroundColor Gray

# Start API in background (redirect output to files)
$apiLogFile = Join-Path $BinDir "vship-api.log"
$apiLogErrFile = Join-Path $BinDir "vship-api.err.log"
$env:SERVER_PORT = "$ApiPort"
$apiProcess = Start-Process -FilePath $apiExe -WorkingDirectory $ProjectDir -PassThru -WindowStyle Hidden -RedirectStandardOutput $apiLogFile -RedirectStandardError $apiLogErrFile

Write-Host "[OK] API started in background (PID: $($apiProcess.Id), log: $apiLogFile)" -ForegroundColor Green
Write-Host ""

# Wait for service to start
Write-Host ">>> Waiting for service..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

# Test health check (try root page as fallback)
Write-Host ">>> Testing API health check..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:$ApiPort/" -TimeoutSec 5 -UseBasicParsing -ErrorAction SilentlyContinue
    if ($response.StatusCode -eq 200) {
        Write-Host "[OK] API is healthy" -ForegroundColor Green
    }
} catch {
    Write-Host "[WARN] Health check failed; check logs in bin/" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  [OK] vShip Started" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Service info:" -ForegroundColor Cyan
Write-Host "  - API:     http://localhost:$ApiPort" -ForegroundColor White
Write-Host "  - Login:   http://localhost:$ApiPort/login" -ForegroundColor White
Write-Host ""
Write-Host "Tips:" -ForegroundColor Cyan
Write-Host "  - API runs in background" -ForegroundColor White
Write-Host "  - Executables and logs are in bin/" -ForegroundColor White
Write-Host "  - Tail API log: Get-Content bin\vship-api.log -Wait -Tail 50" -ForegroundColor White
Write-Host "  - Tail Error log: Get-Content bin\vship-api.err.log -Wait -Tail 50" -ForegroundColor White
Write-Host "  - Stop service: Stop-Process -Id $($apiProcess.Id)" -ForegroundColor White
Write-Host ""
