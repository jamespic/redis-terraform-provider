resource "redis_string" "example" {
  key   = "app:version"
  value = "1.2.3"
}
