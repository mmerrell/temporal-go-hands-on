FROM golang:1.23-bookworm

WORKDIR /opt/exercises

COPY 1_converting/ ./1_converting/
COPY 2_child_workflows/ ./2_child_workflows/
COPY 3_parallel_activities/ ./3_parallel_activities/
COPY 4_local_activities/ ./4_local_activities/

RUN for dir in \
    /opt/exercises/1_converting/practice \
    /opt/exercises/1_converting/solution \
    /opt/exercises/2_child_workflows/practice \
    /opt/exercises/2_child_workflows/solution \
    /opt/exercises/3_parallel_activities/practice \
    /opt/exercises/3_parallel_activities/solution \
    /opt/exercises/4_local_activities/practice \
    /opt/exercises/4_local_activities/solution; do \
    echo "go mod tidy: $dir" && cd "$dir" && go mod tidy; \
    done

CMD ["echo", "Temporal Go sandbox image — use docker cp to extract /opt/exercises"]
