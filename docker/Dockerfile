FROM golang:1.23-bookworm

WORKDIR /opt/exercises

# Copy all exercise code
COPY lab1-child-workflows/ ./lab1-child-workflows/
COPY lab2-local-and-parallel/ ./lab2-local-and-parallel/

# Run go mod tidy in each module to generate go.sum and download dependencies.
# go mod download alone does not generate go.sum — tidy does both.
RUN cd /opt/exercises/lab1-child-workflows/practice && go mod tidy && \
    cd /opt/exercises/lab1-child-workflows/solution && go mod tidy && \
    cd /opt/exercises/lab2-local-and-parallel/practice && go mod tidy && \
    cd /opt/exercises/lab2-local-and-parallel/solution && go mod tidy

CMD ["echo", "Temporal Go sandbox image — use docker cp to extract /opt/exercises"]
