FROM debian:13

# Install Go dari repo Debian + dev tools
RUN apt-get update && apt-get install -y \
    golang-go \
    git \
    curl \
    build-essential \
    unzip \
    ca-certificates \
    sudo \
    && rm -rf /var/lib/apt/lists/*

# Buat user non-root 'vscode' dengan akses sudo tanpa password
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME -s /bin/bash \
    && echo $USERNAME ALL=\(ALL\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

USER $USERNAME
ENV PATH=$PATH:/home/$USERNAME/go/bin

# Install Goose CLI di home directory user vscode
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /workspace