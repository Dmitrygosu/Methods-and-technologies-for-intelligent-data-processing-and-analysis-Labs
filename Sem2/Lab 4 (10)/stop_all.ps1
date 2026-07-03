$ErrorActionPreference = 'SilentlyContinue'

Write-Host "Останавливаю погодный советник..."
Get-Process weather_advisor -ErrorAction SilentlyContinue | Stop-Process -Force
Get-Process python -ErrorAction SilentlyContinue |
    Where-Object { $_.Path -and (Get-CimInstance Win32_Process -Filter "ProcessId=$($_.Id)").CommandLine -match "weather_advisor" } |
    Stop-Process -Force

Write-Host "Останавливаю n8n..."
docker stop n8n *> $null

Write-Host "Выгружаю модели из видеопамяти..."
$ollamaExe = Get-Command ollama -ErrorAction SilentlyContinue
if (-not $ollamaExe) {
    $ollamaExe = Join-Path $env:LOCALAPPDATA "Programs\Ollama\ollama.exe"
}
$loaded = & $ollamaExe ps 2>$null | Select-Object -Skip 1
foreach ($line in $loaded) {
    $name = ($line -split '\s+')[0]
    if ($name) {
        & $ollamaExe stop $name *> $null
        Write-Host "  выгружена: $name"
    }
}

Write-Host ""
Write-Host "Готово. Ollama и Docker остаются работать в фоне -"
Write-Host "для полной остановки Ollama используй llm_kill.cmd."
Start-Sleep -Seconds 3
