resource "google_compute_firewall" "jumpbox-to-jenkins" {
  name    = "${var.env_id}-jumpbox-to-jenkins"
  network = "${google_compute_network.bbl-network.name}"
  description = "Jumpbox to Jenkins for test jobs"

  source_tags = ["${var.env_id}-jumpbox"]

  allow {
    ports    = ["8080"]
    protocol = "tcp"
  }

  target_tags = ["${var.env_id}-internal", "${var.env_id}-bosh-director"]
}
