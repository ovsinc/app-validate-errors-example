data "docker_registry_image" "nginx_img" {
  name = "nginx:1.19"
}

resource "docker_image" "nginx" {
  name          = data.docker_registry_image.nginx_img.name
  pull_triggers = [data.docker_registry_image.nginx_img.sha256_digest]
  keep_locally  = true
}


resource "docker_volume" "html" {
  name   = "html"
  driver = "local"
  driver_opts = {
    "device" = length(var.nginx_source_volume) == 0 ? join("/", [abspath(path.root), "..", "cmd/static/dist/spa"]) : var.nginx_source_volume
    "o"      = "bind"
    "type"   = "none"
  }
}


resource "docker_service" "nginx" {
  name = "nginx"

  endpoint_spec {
    mode = "vip"

    ports {
      protocol       = "tcp"
      publish_mode   = "ingress"
      published_port = 8080
      target_port    = 80
    }
  }


  mode {
    replicated {
      replicas = 1
    }
  }

  task_spec {
    force_update = 0
    networks = [
      docker_network.default_static_network.id
    ]
    runtime = "container"

    container_spec {
      image     = docker_image.nginx.repo_digest
      isolation = "default"
      read_only = false


      mounts {
        read_only = true
        source    = docker_volume.html.name
        target    = "/usr/share/nginx/html"
        type      = "volume"
      }
    }

    placement {
      constraints  = []
      max_replicas = 5
      prefs        = []

      platforms {
        architecture = "amd64"
        os           = "linux"
      }
    }

    resources {
      limits {
        memory_bytes = 536870912
        nano_cpus    = 2000000000
      }

      reservation {
        memory_bytes = 214748364
        nano_cpus    = 500000000
      }
    }

    restart_policy {
      condition    = "any"
      delay        = "3s"
      max_attempts = 20
      window       = "10s"
    }
  }

  update_config {
    delay             = "3s"
    failure_action    = "pause"
    max_failure_ratio = "0.0"
    monitor           = "0s"
    order             = "start-first"
    parallelism       = 1
  }
}
