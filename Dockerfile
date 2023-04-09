FROM alpine:latest

RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.7/main/ > /etc/apk/repositories \
 && apk add -U tzdata \
 && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ARG ARG_PROJECT_NAME

ENV PROJECT_NAME=$ARG_PROJECT_NAME

COPY ./bin/${PROJECT_NAME} /app/bin/
WORKDIR /app/bin
ENTRYPOINT ./${PROJECT_NAME} --registry=${MICRO_REGISTRY} --registry_address=${MICRO_REGISTRY_ADDRESS}