FROM alpine:3.19

# Install build essentials and dev tools
RUN apk add --no-cache \
    make \
    git \
    bash \
    curl \
    unzip \
    build-base \
    vim \
    openssh \
    zsh \
    coreutils

# Install goenv
RUN git clone --depth=1 https://github.com/syndbg/goenv.git ~/.goenv \
    && echo 'export GOENV_ROOT="$HOME/.goenv"' >> ~/.zshrc \
    && echo 'export PATH="$GOENV_ROOT/bin:$PATH"' >> ~/.zshrc \
    && echo 'eval "$(goenv init -)"' >> ~/.zshrc

# Install tfenv
RUN git clone --depth=1 https://github.com/tfutils/tfenv.git ~/.tfenv \
    && ln -s ~/.tfenv/bin/* /usr/local/bin

# Install specific versions
RUN source ~/.zshrc \
    && goenv install 1.22.4 \
    && goenv global 1.22.4 \
    && tfenv install 1.7.4 \
    && tfenv use 1.7.4

# Set up zsh
RUN sh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)" "" --unattended

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Set environment variables
ENV GOENV_ROOT="/root/.goenv"
ENV PATH="/root/.goenv/shims:/root/.goenv/bin:$PATH"
ENV TF_CLI_CONFIG_FILE=/app/.terraformrc
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin

# Download dependencies
RUN source ~/.zshrc && go mod download

# Copy the rest of the code
COPY . .

# Build the provider
RUN source ~/.zshrc && make build

# Create terraform config directory and provider directory
RUN mkdir -p ~/.terraform.d/plugins

# Set zsh as default shell
SHELL ["/bin/zsh", "-c"]
ENTRYPOINT ["/bin/zsh"]