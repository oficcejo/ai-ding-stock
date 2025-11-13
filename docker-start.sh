#!/bin/bash

# ================================
# Stock Analysis System Docker Start Script
# ================================

set -e

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Log functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
}

# Check configuration file
check_config() {
    if [ ! -f "config_stock.json" ]; then
        log_warn "Configuration file config_stock.json does not exist"
        if [ -f "config_stock.json.example" ]; then
            log_info "Copying example configuration file..."
            cp config_stock.json.example config_stock.json
            log_warn "Please edit config_stock.json to fill in your configuration"
            log_warn "Required configuration:"
            echo "  - TDX API URL"
            echo "  - AI API configuration (DeepSeek or Qwen)"
            echo "  - Stock codes to monitor"
            read -p "Press Enter to continue after configuration..."
        else
            log_error "Cannot find example configuration file"
            exit 1
        fi
    fi
}

# Create necessary directories
create_dirs() {
    log_info "Creating necessary directories..."
    mkdir -p logs
    log_success "Directories created"
}

# Start services
start() {
    log_info "Starting stock analysis system..."
    
    check_docker
    check_config
    create_dirs
    
    # Use docker compose or docker-compose
    if docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    else
        COMPOSE_CMD="docker-compose"
    fi
    
    log_info "Building and starting containers..."
    $COMPOSE_CMD up -d --build
    
    log_success "Services started successfully!"
    echo ""
    echo "Stock Analysis System is now running"
    echo ""
    echo "Access URLs:"
    echo "  - Web UI: http://localhost"
    echo "  - API: http://localhost:8080/api/stocks"
    echo ""
    echo "View logs:"
    echo "  - Real-time: $0 logs"
    echo "  - Files: ./logs/"
    echo ""
}

# Stop services
stop() {
    log_info "Stopping stock analysis system..."
    
    if docker compose version &> /dev/null; then
        docker compose down
    else
        docker-compose down
    fi
    
    log_success "Services stopped"
}

# Restart services
restart() {
    log_info "Restarting stock analysis system..."
    stop
    sleep 2
    start
}

# View logs
logs() {
    if docker compose version &> /dev/null; then
        docker compose logs -f --tail=100
    else
        docker-compose logs -f --tail=100
    fi
}

# View status
status() {
    log_info "Service status:"
    if docker compose version &> /dev/null; then
        docker compose ps
    else
        docker-compose ps
    fi
    
    echo ""
    log_info "Container health status:"
    docker ps --filter "name=stock-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
}

# Clean up
clean() {
    log_warn "This will remove all containers, images and volume data!"
    read -p "Are you sure you want to continue? (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        log_info "Cleaning up..."
        
        if docker compose version &> /dev/null; then
            docker compose down -v --rmi all
        else
            docker-compose down -v --rmi all
        fi
        
        log_success "Cleanup completed"
    else
        log_info "Cancelled"
    fi
}

# Update
update() {
    log_info "Updating stock analysis system..."
    
    if docker compose version &> /dev/null; then
        docker compose pull
        docker compose up -d --build
    else
        docker-compose pull
        docker-compose up -d --build
    fi
    
    log_success "Update completed"
}

# Enter container shell
shell() {
    log_info "Entering backend container..."
    docker exec -it stock-analyzer sh
}

# Usage
usage() {
    echo "Stock Analysis System Docker Management Script"
    echo ""
    echo "Usage: $0 {start|stop|restart|logs|status|clean|update|shell|help}"
    echo ""
    echo "Commands:"
    echo "  start   - Start services"
    echo "  stop    - Stop services"
    echo "  restart - Restart services"
    echo "  logs    - View real-time logs"
    echo "  status  - View service status"
    echo "  clean   - Clean all data (dangerous)"
    echo "  update  - Update and restart services"
    echo "  shell   - Enter backend container"
    echo "  help    - Show this help message"
    echo ""
}

# Main function
main() {
    case "$1" in
        start)
            start
            ;;
        stop)
            stop
            ;;
        restart)
            restart
            ;;
        logs)
            logs
            ;;
        status)
            status
            ;;
        clean)
            clean
            ;;
        update)
            update
            ;;
        shell)
            shell
            ;;
        help|--help|-h)
            usage
            ;;
        *)
            usage
            exit 1
            ;;
    esac
}

# Execute main function
main "$@"
