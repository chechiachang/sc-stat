variable "REGISTRY" {
  default = "docker.io"
}

variable "REPOSITORY" {
  default = "chechiachang/sc-stat"
}

variable "TAG" {
  default = "latest"
}

variable "APP_CODE_VERSION" {
  default = "dev"
}

variable "BRANCH" {
  default = "main"
}

target "default" {
  dockerfile = "Dockerfile"
  context = "."
  args = {
    APP_CODE_VERSION = "${APP_CODE_VERSION}"
  }
  labels = {
    "org.opencontainers.image.created" = timestamp()
    "org.opencontainers.image.version" = "${APP_CODE_VERSION}"
  }
  platforms = ["linux/amd64", "linux/arm64"]
  tags = [
    "${REGISTRY}/${REPOSITORY}:${TAG}"
  ]
}
