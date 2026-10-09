param([ValidateSet('run','build','test')][string]$Task='run')
$ErrorActionPreference='Stop'
Push-Location $PSScriptRoot
try {
    switch ($Task) {
        'run' { go run ./cmd/server }
        'build' {
            New-Item -ItemType Directory -Path (Join-Path $PSScriptRoot 'bin') -Force | Out-Null
            go build -o ./bin/server.exe ./cmd/server
        }
        'test' { go test ./... }
    }
    if ($LASTEXITCODE -ne 0) { throw "Go exited with code $LASTEXITCODE" }
} finally { Pop-Location }
