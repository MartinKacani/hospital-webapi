param (
    $command
)

if (-not $command)  {
    $command = "start"
}

$ProjectRoot = "${PSScriptRoot}/.."

$env:HOSPITAL_API_ENVIRONMENT="Development"
$env:HOSPITAL_API_PORT="8080"

switch ($command) {
    "start" {
        go run ${ProjectRoot}/cmd/hospital-api-service
    }
    "openapi" {
        docker run --rm -ti -v ${ProjectRoot}:/local openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
    }
    default {
        throw "Unknown command: $command"
    }
}