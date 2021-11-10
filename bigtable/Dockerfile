FROM alpine:latest

RUN apk --no-cache add ca-certificates wget python3 bash openjdk8-jre

RUN wget https://dl.google.com/dl/cloudsdk/channels/rapid/google-cloud-sdk.zip \
  && unzip google-cloud-sdk.zip && rm google-cloud-sdk.zip

RUN /google-cloud-sdk/install.sh \
  --usage-reporting=true --path-update=true --bash-completion=true \
  --rc-path=/.bashrc --additional-components beta bigtable

EXPOSE 8086

ENTRYPOINT ["/google-cloud-sdk/bin/gcloud"]

CMD ["beta", "emulators", "bigtable", "start", "--host-port=0.0.0.0:8086"]
