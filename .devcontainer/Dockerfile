FROM golang:latest

# Update and install dependencies
RUN apt-get update && apt-get install -y curl bash xz-utils

# Set environment variable to skip Turso signup during installation
ENV TURSO_INSTALL_SKIP_SIGNUP=true

# Set the environment variables for the flyctl installation
ENV PATH="/root/.fly/bin:$PATH"

# Install Turso CLI
RUN curl -sSfL https://get.tur.so/install.sh | sh

# Install Fly.io CLI
RUN curl -L https://fly.io/install.sh | sh

COPY . /workspace

WORKDIR /workspace

# IMPORTANT: If you need to install Go dependencies or perform other setup
# that relies on the base image being fully set up, do it AFTER installing
# the CLIs.  This ensures the CLIs are available when those commands run.
# Example:
# COPY go.mod go.sum ./
# RUN go mod download

# You likely already have this in your Dockerfile, but it's crucial
# to keep it after the CLI installations if it was previously there.
# COPY . /workspace  # Or whatever your copy command is

# If you need to set a working directory, do it after the installation
# and any dependency installation.
# WORKDIR /workspace

# If you have a CMD instruction, it should also be at the end.
# CMD ["your", "command"]