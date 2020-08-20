resource "digitalocean_droplet" "hexlet-basics" {
    image = "ubuntu-18-04-x64"
    name = "hexlet-basics"
    region = "fra1"
    size = "s-1vcpu-1gb"
    private_networking = true
}
