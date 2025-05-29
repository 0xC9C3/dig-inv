ARG BASE_IMAGE=node:lts

FROM $BASE_IMAGE AS dev

# install openjdk17
RUN apt update && apt install -y openjdk-17-jdk && \
    apt clean && rm -rf /var/lib/apt/lists/*

# @todo deduplicate from backend Dockerfile
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin

RUN apt update && apt install -y cargo && \
    curl -L --proto '=https' --tlsv1.2 -sSf https://raw.githubusercontent.com/cargo-bins/cargo-binstall/main/install-from-binstall-release.sh | bash && \
    cargo binstall watchexec-cli --no-confirm && \
    apt remove -y cargo && \
    apt autoremove -y && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    mv /root/.cargo/bin/watchexec /usr/local/bin/watchexec