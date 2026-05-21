resource "redis_hash" "example" {
  key = "user:42"
  fields = {
    name  = "Alice"
    email = "alice@example.com"
    role  = "admin"
  }
}
