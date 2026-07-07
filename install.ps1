param(
    [string]$Version = "v0.1.0"
)

$Repo = "Fenomen-Alex/go-chat"
$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "arm64" }
$Binary = "chat-windows-$Arch.exe"
$Url = "https://github.com/$Repo/releases/download/$Version/$Binary"
$ChecksumsUrl = "https://github.com/$Repo/releases/download/$Version/checksums.txt"

Write-Output "Downloading $Binary $Version..."

$OutDir = "$env:TEMP\go-chat-install"
New-Item -ItemType Directory -Force -Path $OutDir | Out-Null
$OutFile = "$OutDir\$Binary"

[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
Invoke-WebRequest -Uri $Url -OutFile $OutFile -UseBasicParsing

try {
    $Checksums = Invoke-WebRequest -Uri $ChecksumsUrl -UseBasicParsing
    $Expected = ($Checksums.Content -split "`n" | Where-Object { $_ -match $Binary } | ForEach-Object { $_ -split "\s+" | Select-Object -First 1 })
    $Actual = (Get-FileHash $OutFile -Algorithm SHA256).Hash.ToLower()
    if ($Expected -and $Actual -ne $Expected) {
        Write-Error "Checksum mismatch"
        exit 1
    }
    Write-Output "Checksum verified"
} catch {
    Write-Warning "Could not verify checksum"
}

$InstallDir = "$env:ProgramFiles\go-chat"
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
Move-Item -Force $OutFile "$InstallDir\chat.exe"

$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
}

Write-Output "Installed: $InstallDir\chat.exe"
Write-Output "Run: chat"
