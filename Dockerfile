# Nix builder
FROM nixos/nix:latest AS builder

# Copy our source and setup our working dir.
COPY . /tmp/build
WORKDIR /tmp/build/bsf

# Build runtime package dependencies
RUN nix \
    --extra-experimental-features "nix-command flakes" \
    --option filter-syscalls false \
    build

# Build additional packages we need for runtime
RUN nix \
    --extra-experimental-features "nix-command flakes" \
    --option filter-syscalls false \
    build .#runtimeEnvs.aarch64-linux.runtime -o runtimeEnv



# Copy the Nix store closure into a directory. The Nix store closure is the
# entire set of Nix store values that we need for our build and custom environment.
RUN mkdir /tmp/nix-store-closure
RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure
RUN cp -R $(nix-store -qR runtimeEnv/) /tmp/nix-store-closure


# # Final image is based on scratch. We copy a bunch of Nix dependencies
# # but they're fully self-contained so we don't need Nix anymore.
FROM scratch

WORKDIR /result

# Copy /nix/store
COPY --from=builder /tmp/nix-store-closure /nix/store
# Add symlink to result
COPY --from=builder /tmp/build/bsf/result /result
COPY --from=builder /tmp/build/bsf/runtimeEnv /result/env
# Add /result/env to the PATH
ENV SSL_CERT_FILE="/result/env/etc/ssl/certs/ca-bundle.crt"
ENV PATH="/result/env/bin:${PATH}"

CMD [  "/result/bin/bsf \n" ]

