# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

packer {
  required_plugins {
    zstack = {
      version = ">=v0.0.1"
      source  = "github.com/kubezstack/zstack"
    }
  }
}

variable "access_key" {
  type = string
}
variable "secret_key" {
  type = string
}
variable "host" {
  type = string
}

source "zstack" "example" {
  access_key_id        = var.access_key
  access_secret_key    = var.secret_key
  auth_type            = "access_key"
  communicator         = "ssh"
  host                 = var.host
  image_name           = "test-packer"
  image_backup_storage = "833a2d77ecfe44f591d8673f516e7970"
  instance_offering    = "9192fc8c8f434f06b6abe9e4cff85da3"
  instance_network     = "35df90cff65244fbb2601d60b98ed3e3"
  source_image         = "9eb237dee6ef401fa643810f743de489"
  ssh_username         = "ubuntu"
  ssh_port             = 22
}

build {
  sources = ["source.zstack.example"]
  provisioner "shell" {
    inline = [
      "echo Hello From ${source.type} ${source.name}"
    ]
  }
}
