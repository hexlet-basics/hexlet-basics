resource "twc_vpc" "hexlet_basics" {
  name = "Hexlet basics network"
  subnet_v4 = "10.3.1.0/24"
  location = var.location
}
