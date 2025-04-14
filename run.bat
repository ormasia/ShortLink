chcp 65001
docker build -f Dockerfile.api -t shortlink/api:latest .
docker build -f Dockerfile.user -t shortlink/user:latest .
docker build -f Dockerfile.short -t shortlink/short:latest .
docker-compose up -d