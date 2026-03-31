#!/bin/bash
set -ve

#=============================================================
# ForgetURL Database Initialization Script
# Reads MySQL config from app/conf/local.toml, auto-creates database and tables
#=============================================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CONF_FILE="$SCRIPT_DIR/app/conf/local.toml"
SQL_DIR="$SCRIPT_DIR/app/dal/gensql"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC}  $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; exit 1; }

#---------- Parse local.toml config ----------
toml_val() {
    awk -F'"' -v key="$1" '$0 ~ "^[[:space:]]*"key"[[:space:]]*=" {print $2; exit}' "$CONF_FILE"
}

parse_conf() {
    if [[ ! -f "$CONF_FILE" ]]; then
        error "Config file not found: $CONF_FILE"
    fi

    MYSQL_HOST=$(toml_val "Host")
    MYSQL_PORT=$(toml_val "Port")
    MYSQL_DB=$(toml_val "Name")
    MYSQL_USER=$(toml_val "User")
    MYSQL_PASS=$(toml_val "Passwd")

    REDIS_HOST=$(toml_val "host")
    REDIS_PORT=$(toml_val "port" | tr -d ':')

    info "MySQL config: ${MYSQL_USER}@${MYSQL_HOST}:${MYSQL_PORT}/${MYSQL_DB}"
    info "Redis config: ${REDIS_HOST}:${REDIS_PORT}"
}

#---------- Check dependencies ----------
check_deps() {
    if ! command -v mysql &>/dev/null; then
        warn "mysql client command not found"
        warn "macOS install: brew install mysql-client && echo 'export PATH=\"/opt/homebrew/opt/mysql-client/bin:\$PATH\"' >> ~/.zshrc"
        error "Please install mysql client and try again"
    fi
}

#---------- Build mysql connection args ----------
build_mysql_args() {
    MYSQL_ARGS=(-u "$MYSQL_USER" -p"$MYSQL_PASS")
    # Try TCP first, fall back to socket on failure
    if mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASS" -e "SELECT 1" &>/dev/null; then
        MYSQL_ARGS=(-h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASS")
        MYSQL_MODE="TCP (${MYSQL_HOST}:${MYSQL_PORT})"
    elif mysql -u "$MYSQL_USER" -p"$MYSQL_PASS" -e "SELECT 1" &>/dev/null; then
        MYSQL_ARGS=(-u "$MYSQL_USER" -p"$MYSQL_PASS")
        MYSQL_MODE="Socket"
    else
        MYSQL_ARGS=()
        MYSQL_MODE=""
    fi
}

#---------- Check MySQL connectivity ----------
check_mysql_connection() {
    info "Checking MySQL connection..."
    build_mysql_args
    if [[ -n "$MYSQL_MODE" ]]; then
        info "MySQL connected successfully (${MYSQL_MODE})"
    else
        warn "Unable to connect to MySQL"
        echo ""
        echo "  You can start MySQL using one of the following methods:"
        echo ""
        echo "  Method 1 - Docker (recommended):"
        echo "    docker run -d --name forgeturl-mysql \\"
        echo "      -p ${MYSQL_PORT}:3306 \\"
        echo "      -e MYSQL_ROOT_PASSWORD=${MYSQL_PASS} \\"
        echo "      -e MYSQL_DATABASE=${MYSQL_DB} \\"
        echo "      mysql:8.0"
        echo ""
        echo "  Method 2 - Homebrew:"
        echo "    brew install mysql && brew services start mysql"
        echo ""
        error "Please start MySQL service and re-run this script"
    fi
}

#---------- Check Redis connectivity ----------
check_redis_connection() {
    info "Checking Redis connection..."
    if command -v redis-cli &>/dev/null; then
        if redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" ping 2>/dev/null | grep -q PONG; then
            info "Redis connected successfully"
            return
        fi
    fi
    warn "Unable to connect to Redis (${REDIS_HOST}:${REDIS_PORT}), please make sure Redis is running"
    echo ""
    echo "  Method 1 - Docker:"
    echo "    docker run -d --name forgeturl-redis -p ${REDIS_PORT}:6379 redis:7"
    echo ""
    echo "  Method 2 - Homebrew:"
    echo "    brew install redis && brew services start redis"
    echo ""
    warn "Redis not ready, but will continue with database initialization..."
}

#---------- Create database ----------
create_database() {
    info "Creating database ${MYSQL_DB} (if not exists)..."
    mysql "${MYSQL_ARGS[@]}" \
        -e "CREATE DATABASE IF NOT EXISTS \`${MYSQL_DB}\` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"
    info "Database ${MYSQL_DB} is ready"
}

#---------- Create tables ----------
create_tables() {
    local SQL_FILES=("user.sql" "unique_pid.sql" "page.sql" "user_page.sql" "openclaw_api_key.sql" "tmp_bookmark.sql")

    for sql_file in "${SQL_FILES[@]}"; do
        local filepath="${SQL_DIR}/${sql_file}"
        if [[ ! -f "$filepath" ]]; then
            error "SQL file not found: $filepath"
        fi

        local table_name
        table_name=$(awk '/CREATE TABLE/ {gsub(/`/,"",$3); print $3; exit}' "$filepath")

        local exists
        exists=$(mysql "${MYSQL_ARGS[@]}" \
            -D "$MYSQL_DB" -N -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='${MYSQL_DB}' AND table_name='${table_name}';" 2>/dev/null)

        if [[ "$exists" == "1" ]]; then
            warn "Table ${table_name} already exists, skipping"
        else
            info "Creating table ${table_name} ..."
            # Filter out comment lines (starting with #) and DROP TABLE statements
            grep -v '^\s*#' "$filepath" | grep -iv 'DROP TABLE' | \
                mysql "${MYSQL_ARGS[@]}" -D "$MYSQL_DB"
            info "Table ${table_name} created successfully"
        fi
    done
}

#---------- Verify results ----------
verify() {
    info "Verifying database tables..."
    echo ""
    mysql "${MYSQL_ARGS[@]}" -D "$MYSQL_DB" -e "SHOW TABLES;" 2>/dev/null
    echo ""
    info "All tables created successfully!"
}

#---------- Main ----------
main() {
    echo "=========================================="
    echo "  ForgetURL Database Initialization"
    echo "=========================================="
    echo ""

    parse_conf
    check_deps
    check_mysql_connection
    check_redis_connection
    create_database
    create_tables
    verify

    echo ""
    info "Initialization complete! You can now start the service."
}

main "$@"
