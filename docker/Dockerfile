FROM golang:1.22-bookworm

WORKDIR /opt/exercises

# Copy all exercise code
COPY lab1-child-workflows/ ./lab1-child-workflows/
COPY lab2-local-and-parallel/ ./lab2-local-and-parallel/

# Pre-download all Go module dependencies so the sandbox VM
# doesn't need internet access during the lab
RUN cd /opt/exercises/lab1-child-workflows/practice && go mod download && \
    cd /opt/exercises/lab1-child-workflows/solution && go mod download && \
    cd /opt/exercises/lab2-local-and-parallel/practice && go mod download && \
    cd /opt/exercises/lab2-local-and-parallel/solution && go mod download

CMD ["echo", "Temporal Go sandbox image — use docker cp to extract /opt/exercises"]
