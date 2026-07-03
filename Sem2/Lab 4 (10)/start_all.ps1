$ErrorActionPreference = 'SilentlyContinue'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $root

function Test-Site {
    try {
        $r = Invoke-WebRequest -Uri "http://127.0.0.1:5000/" -UseBasicParsing -TimeoutSec 2
        return $r.StatusCode -eq 200
    } catch {
        return $false
    }
}

Write-Host "Проверяю Ollama..."
$ollamaRunning = Get-Process -Name "ollama app" -ErrorAction SilentlyContinue
if (-not $ollamaRunning) {
    $ollamaExe = Join-Path $env:LOCALAPPDATA "Programs\Ollama\ollama app.exe"
    if (Test-Path $ollamaExe) {
        Write-Host "Поднимаю Ollama..."
        Start-Process $ollamaExe
        Start-Sleep -Seconds 3
    } else {
        Write-Host "Ollama не найдена по пути $ollamaExe"
        Write-Host "Установи с https://ollama.com/download и запусти скрипт снова."
    }
}

if (-not (Test-Site)) {
    $exePath = Join-Path $root "Задание 1\weather_advisor\dist\weather_advisor.exe"
    if (Test-Path $exePath) {
        Write-Host "Запускаю погодный советник (exe)..."
        Start-Process -FilePath $exePath -WindowStyle Minimized
    } else {
        Write-Host "exe не найден, запускаю из исходников (python)..."
        $appDir = Join-Path $root "Задание 1\weather_advisor"
        Start-Process -FilePath "cmd.exe" -ArgumentList "/c cd /d `"$appDir`" && python app.py" -WindowStyle Minimized
    }

    Write-Host "Жду, пока сайт поднимется..."
    $tries = 0
    while (-not (Test-Site) -and $tries -lt 60) {
        Start-Sleep -Seconds 1
        $tries++
    }
}

if (Test-Site) {
    Write-Host "Сайт готов, открываю браузер."
    Start-Process "http://127.0.0.1:5000"
} else {
    Write-Host "Сайт не поднялся за 60 секунд."
    Write-Host "Проверь окно 'Погодный советник' и что Ollama установлена и запущена."
    Read-Host "Нажми Enter для выхода"
    exit 1
}

$dockerCmd = Get-Command docker -ErrorAction SilentlyContinue
if (-not $dockerCmd) {
    Write-Host "Docker не установлен - n8n (задания 3-4) недоступны."
    Write-Host "Для их показа: https://www.docker.com/products/docker-desktop"
} else {
    docker info *> $null
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Docker установлен, но не запущен - открой Docker Desktop"
        Write-Host "и запусти этот скрипт снова, если нужны задания 3-4."
    } else {
        docker start n8n *> $null
        if ($LASTEXITCODE -ne 0) {
            $outDir = Join-Path $root "Задание 4\n8n\output"
            docker run -d --name n8n -p 5678:5678 `
                -e N8N_SECURE_COOKIE=false -e GENERIC_TIMEZONE=Europe/Ulyanovsk `
                -e N8N_RESTRICT_FILE_ACCESS_TO=/data/output `
                -v n8n_data:/home/node/.n8n `
                -v "${outDir}:/data/output" n8nio/n8n:latest *> $null
        }
        Write-Host "n8n:   http://localhost:5678  (задания 3-4)"
    }
}

Write-Host ""
Write-Host "Погодный советник: http://127.0.0.1:5000"
Write-Host "Первый совет после запуска может занять до минуты - модель греется."
Start-Sleep -Seconds 5
