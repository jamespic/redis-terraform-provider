provider "redis" {
  addr = "localhost:6379"
  # password = "secret"   # or set REDIS_PASSWORD
  # username = "default"  # or set REDIS_USERNAME (Redis 6+ ACL)
  # db       = 0

  # TLS — enable when connecting to Redis over an encrypted channel
  # tls = true

  # Disable certificate verification — only for dev/test, never in production
  # tls_insecure_skip_verify = true
}
