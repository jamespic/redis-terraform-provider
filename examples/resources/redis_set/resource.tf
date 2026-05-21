resource "redis_set" "example" {
  key     = "article:1:tags"
  members = ["terraform", "redis", "infrastructure"]
}
