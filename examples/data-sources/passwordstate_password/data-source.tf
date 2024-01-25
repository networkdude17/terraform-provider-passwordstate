# Get Username and Password from PasswordState using the Password ID (PID)
data "passwordstate_password" "example" {
  passwordid = 00000
}

# Output the password, sensitive must be set to true
output "password" {
  value     = data.passwordstate_password.example.password
  sensitive = true
}

# Output the username
output "username" {
  value = data.passwordstate_password.example.username
}
