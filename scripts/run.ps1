param (
    $command
)

if (-not $command)  {
    $command = "start"
}

$ProjectRoot = "${PSScriptRoot}/.."

$env:HOSPITAL_API_ENVIRONMENT="Development"
$env:HOSPITAL_API_PORT="8080"
$env:HOSPITAL_API_MONGODB_USERNAME="root"
$env:HOSPITAL_API_MONGODB_PASSWORD="neUhaDnes"
$env:HOSPITAL_API_MONGODB_DATABASE="tjmk-hospital-wl"

function mongo {
    docker compose --file ${ProjectRoot}/deployments/docker-compose/compose.yaml $args
}

switch ($command) {
    "openapi" {
        docker run --rm -ti  -v ${ProjectRoot}:/local openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
    }
    "start" {
        try {
            mongo up --detach
            go run ${ProjectRoot}/cmd/hospital-api-service
        } finally {
            mongo down
        }
    }
    "mongo" {
        mongo up
    }
    default {
        throw "Unknown command: $command"
    }
}