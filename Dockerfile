FROM golang:1.17-alpine AS server
RUN apk update && apk add ca-certificates && apk add tzdata
ADD ./server /go/src/github.com/riversy/chat-room/server
WORKDIR /go/src/github.com/riversy/chat-room/server
RUN CGO_ENABLED=0 go build -o ./chat-room main.go


FROM node:16 AS client
ARG REACT_APP_WS_PATH
ADD ./client /go/src/github.com/riversy/chat-room/client
WORKDIR /go/src/github.com/riversy/chat-room/client
RUN npm install
RUN npm run build


FROM scratch AS bin
ENV STATIC_FILES_PATH=/build
COPY --from=server /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=server /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=server /go/src/github.com/riversy/chat-room/server/chat-room /
COPY --from=client /go/src/github.com/riversy/chat-room/client/build /build
CMD [ "./chat-room" ]

