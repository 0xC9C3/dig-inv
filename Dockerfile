ARG BASE_IMAGE=golang:1.24.3-bookworm

FROM bufbuild/buf:1 AS buf

FROM $BASE_IMAGE AS dev

COPY --from=buf /usr/local/bin/buf /usr/local/bin/buf

# just
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin


# watchexec
RUN apt update && apt install -y cargo && \
    curl -L --proto '=https' --tlsv1.2 -sSf https://raw.githubusercontent.com/cargo-bins/cargo-binstall/main/install-from-binstall-release.sh | bash && \
    cargo binstall watchexec-cli --no-confirm && \
    apt remove -y cargo && \
    apt autoremove -y && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    mv /root/.cargo/bin/watchexec /usr/local/bin/watchexec


# golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6

WORKDIR /app

FROM $BASE_IMAGE AS build

COPY . /app
WORKDIR /app
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/app

FROM scratch AS prod

COPY --from=build /app/bin/app /app/bin/app

WORKDIR /app

CMD ["/app/bin/app"]
