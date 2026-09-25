package amocrm

// The contract is Hexlet's own TypeSpec output of the amoCRM v4 API, pinned
// by commit so a regeneration is reproducible. ogen is pinned by go.mod.
//go:generate go tool ogen --config ogen.yml --target generated --package generated --clean https://raw.githubusercontent.com/Hexlet/amocrm-api/5a24b11242c66c5b64e4a0d1c107a45ee9dfbab3/tsp-output/schema/openapi.yaml
