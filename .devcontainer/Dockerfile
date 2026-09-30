FROM debian:13

# Install Go dari repo Debian + dev tools + locales (fix setlocale warning) + net & sys tools
RUN apt-get update && apt-get install -y \
    golang-go \
    git \
    curl \
    wget \
    nano \
    build-essential \
    unzip \
    zip \
    7zip \
    ca-certificates \
    sudo \
    locales \
    iputils-ping \
    htop \
    btop \
    screen \
    fastfetch \
    gdu \
    && rm -rf /var/lib/apt/lists/* \
    && sed -i '/en_US.UTF-8/s/^# //g' /etc/locale.gen \
    && locale-gen en_US.UTF-8 \
    && update-locale LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8

# Install Node.js 24 dari NodeSource + global CLI: command-code, opencode, claude-code
RUN curl -fsSL https://deb.nodesource.com/setup_24.x -o /tmp/nodesource_setup.sh \
    && bash /tmp/nodesource_setup.sh \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/* /tmp/nodesource_setup.sh \
    && node -v && npm -v \
    && npm i -g --allow-scripts=@anthropic-ai/claude-code command-code @anthropic-ai/claude-code @opencode/cli \
    && npm cache clean --force \
    && command-code --version && claude --version && opencode --version

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

WORKDIR /workspace/test-devcontainer-go