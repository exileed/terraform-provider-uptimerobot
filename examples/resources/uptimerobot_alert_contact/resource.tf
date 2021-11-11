# Create a alert contact
resource "uptimerobot_alert_contact" "test" {
  friendly_name = "Email"
  type          = "mail"
  url           = "test@example.com"
}
