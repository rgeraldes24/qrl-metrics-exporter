FROM gcr.io/distroless/static-debian11:latest
COPY qrl-metrics-exporter* /qrl-metrics-exporter
ENTRYPOINT ["/qrl-metrics-exporter"]
