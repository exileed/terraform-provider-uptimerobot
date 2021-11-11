# Create a monitor
resource "uptimerobot_monitor" "web" {
  friendly_name = "My Monitor"
  type          = "http"
  url           = "http://example.com"
}