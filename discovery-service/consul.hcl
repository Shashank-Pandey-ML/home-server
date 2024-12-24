# consul.hcl - Consul Configuration File

# Enable Consul UI
ui = true

# Bind address for the agent to bind to
bind_addr = "0.0.0.0"

# Advertise address for the agent to advertise
advertise_addr = "localhost"  # This can be changed based on the actual setup

# Enable the DNS resolver
enable_dns = true

# Enable the Serf LAN events (useful for cluster formation)
serf_lan_bind = "0.0.0.0"

# Default data center (can be useful when running multi-datacenter setups)
datacenter = "dc1"
