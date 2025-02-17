set export := true

build:
    docker compose pull
    docker compose build

up:
    docker compose up -d --build

down:
    docker compose down

reset:
    docker compose down -v --remove-orphans
    just up
