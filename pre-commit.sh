#!/bin/bash

docker compose run --build --entrypoint "just pre-commit" --rm watch-dev-server
docker compose run --build --entrypoint "npm run test:unit && npm run lint" --rm frontend