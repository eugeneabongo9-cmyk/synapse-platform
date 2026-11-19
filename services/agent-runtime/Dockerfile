FROM rust:1.78-alpine AS builder

WORKDIR /app
COPY Cargo.toml Cargo.lock ./
COPY src ./src

RUN cargo build --release

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/target/release/agent-runtime /agent-runtime

USER nobody
ENTRYPOINT ["/agent-runtime"]
