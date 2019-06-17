FROM debian:latest

RUN apt-get update && apt-get install -y curl
COPY ./web_app /web_app

CMD [ "/web_app" ]