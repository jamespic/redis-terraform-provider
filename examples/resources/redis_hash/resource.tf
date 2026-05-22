resource "redis_hash" "example" {
  key = "user:42"
  fields = {
    name  = "Ozzy"
    email = "ozzy@example.com"
    role  = "admin"
  }
}
