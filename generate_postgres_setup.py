#!/usr/bin/env python3
"""
PostgreSQL Database Setup Generator for Home Server Microservices

This script automatically generates PostgreSQL initialization files for a microservices
architecture by scanning service configuration files and environment variables.

The script performs the following operations:
1. Scans all top-level directories for service configuration files (config.yaml)
2. Extracts database configuration from each service's config.yaml
3. Reads database passwords from each service's .env file
4. Generates a consolidated PostgreSQL init.sql script for database/user creation
5. Creates a unified .env file for the PostgreSQL container with all required passwords

Generated Files:
- ./postgres/init.sql: SQL commands to create databases and users for each service
- ./postgres/.env: Environment variables with database passwords for Docker Compose

This automation ensures consistent database setup across all microservices and
eliminates manual configuration errors when adding new services.

Usage:
    python3 generate_postgres_setup.py
"""

import os
import yaml

# Configuration constants for file paths and directory structure
BASE_DIR = os.getcwd()
POSTGRES_DIR = os.path.join(BASE_DIR, "postgres")
INIT_SQL_PATH = os.path.join(POSTGRES_DIR, "init.sql")
ENV_PATH = os.path.join(POSTGRES_DIR, ".env")

# Ensure the postgres directory exists for output files
os.makedirs(POSTGRES_DIR, exist_ok=True)

def get_top_level_folders(base_dir):
    """
    Retrieve all top-level directories in the given base directory.
    
    This function scans the base directory and returns a list of all subdirectories
    that could potentially contain microservice configurations. It filters out
    files and only returns actual directories.
    
    Args:
        base_dir (str): The base directory path to scan for subdirectories
        
    Returns:
        list[str]: List of absolute paths to all top-level directories
        
    Example:
        >>> get_top_level_folders("/home/server")
        ['/home/server/auth', '/home/server/gateway', '/home/server/ui-service']
    """
    return [
        os.path.join(base_dir, d)
        for d in os.listdir(base_dir)
        if os.path.isdir(os.path.join(base_dir, d))
    ]

def parse_config_yaml(path):
    """
    Parse a YAML configuration file and return its contents as a Python dictionary.
    
    This function safely loads YAML configuration files used by microservices
    to define their database requirements, service names, and other settings.
    
    Args:
        path (str): Absolute path to the YAML configuration file
        
    Returns:
        dict: Parsed YAML content as a Python dictionary
    """
    with open(path, 'r') as f:
        return yaml.safe_load(f)

def extract_db_password(env_file_path, var_name):
    """
    Extract a database password from a service's .env file.
    
    This function reads an environment file and searches for a specific variable
    containing the database password. It handles quoted and unquoted values
    and returns None if the variable is not found or the file doesn't exist.
    
    Args:
        env_file_path (str): Path to the .env file to read
        var_name (str): Name of the environment variable to extract (e.g., "DB_PASSWORD")
        
    Returns:
        str or None: The password value if found, None if not found or file missing
    """
    if not os.path.exists(env_file_path):
        return None
    with open(env_file_path, 'r') as f:
        for line in f:
            if line.startswith(var_name + "="):
                return line.strip().split('=', 1)[1].strip('"').strip("'")
    return None

# Initialize containers for generated SQL commands and environment variables
sql_entries = []  # Will contain SQL commands for database/user creation
env_entries = [   # Will contain environment variables for PostgreSQL container
    "# Root user for initial setup (used by Docker only)",
    "POSTGRES_USER=postgres",
    "POSTGRES_PASSWORD=password"
]

if __name__ == "__main__":
    print("🔍 Scanning for service configurations...")
    # Main processing loop: scan all service directories for database configuration
    for folder in get_top_level_folders(BASE_DIR):
        config_path = os.path.join(folder, "config.yaml")
        env_path = os.path.join(folder, ".env")

        # Skip directories that don't have a config.yaml file (not a service)
        if not os.path.isfile(config_path):
            continue

        # Parse the service configuration to extract database requirements
        config = parse_config_yaml(config_path)

        # Extract required database configuration values from the YAML
        service_name = config.get("service", {}).get("name")
        db_name = config.get("database", {}).get("name")
        db_user = config.get("database", {}).get("user")

        # Validate that all required configuration values are present
        if not (service_name and db_name and db_user):
            print(f"Skipping {folder} due to missing values in config.yaml")
            continue

        # Define password variable names for consistent naming convention
        password_var_in_service = f"DB_PASSWORD"  # Variable name in service's .env file
        password_var_in_psql = f"{service_name.upper()}_DB_PASSWORD"  # Variable name in postgres .env
        
        # Extract the actual password value from the service's .env file
        db_password = extract_db_password(env_path, password_var_in_service)

        # Skip services that don't have a password configured
        if db_password is None or db_password == "":
            print(f"❌ Error: {password_var_in_service} not found in {env_path}")
            exit(1)

        # Generate SQL commands for database and user creation
        # Note: We use the actual password here, not a variable, because PostgreSQL's
        # docker-entrypoint-initdb.d doesn't expand environment variables in .sql files
        sql_entries.append(f"""-- Create DB and user for {service_name} service
        CREATE DATABASE {db_name};
        CREATE USER {db_user} WITH ENCRYPTED PASSWORD '{db_password}';
        GRANT ALL PRIVILEGES ON DATABASE {db_name} TO {db_user};

        -- Grant schema permissions
        \\c {db_name}
        GRANT ALL ON SCHEMA public TO {db_user};
        GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO {db_user};
        GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO {db_user};
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO {db_user};
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO {db_user};
        """)
        sql_entries.append(f"")  # Add blank line for readability

        # Add corresponding environment variables for PostgreSQL container
        env_entries.append(f"")
        env_entries.append(f"# Database configuration for {service_name} service.")
        env_entries.append(f"{password_var_in_psql}={db_password}")

    # Write the generated SQL initialization script
    with open(INIT_SQL_PATH, 'w') as sql_file:
        sql_file.write("\n".join(sql_entries))

    # Write the consolidated environment variables file for PostgreSQL
    with open(ENV_PATH, 'w') as env_file:
        env_file.write("\n".join(env_entries))

    print("✅ Generated ./postgres/init.sql and ./postgres/.env")
