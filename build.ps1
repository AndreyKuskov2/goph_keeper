# GophKeeper CLI Client Build Script

$VERSION = "1.0.0"
$BUILD_TIME = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
$GIT_COMMIT = "unknown"

Write-Host "Building GophKeeper CLI Client..." -ForegroundColor Green
Write-Host "Version: $VERSION" -ForegroundColor Yellow
Write-Host "Build Time: $BUILD_TIME" -ForegroundColor Yellow

$ldflags = "-X main.version=$VERSION -X main.buildTime='$BUILD_TIME' -X main.gitCommit=$GIT_COMMIT"

go build -ldflags $ldflags -o goph_keeper_client.exe ./cmd/goph_keeper_client

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green
    Write-Host "Binary: goph_keeper_client.exe" -ForegroundColor Cyan
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}
