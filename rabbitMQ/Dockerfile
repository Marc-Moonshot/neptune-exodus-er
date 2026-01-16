FROM rabbitmq:4.2.2-management-alpine

# Install dedupe plugin
RUN apk add --no-cache curl unzip && \
  curl -L -o /tmp/rabbitmq_message_deduplication-0.7.0.zip \
  https://github.com/noxdafox/rabbitmq-message-deduplication/releases/download/0.7.3/plugins-rmqv4.2.x-erl27-elx1.18.zip && \
  unzip /tmp/rabbitmq_message_deduplication-0.7.0.zip -d /plugins && \
  rm /tmp/rabbitmq_message_deduplication-0.7.0.zip && \
  rabbitmq-plugins enable --offline rabbitmq_message_deduplication && \
  apk del curl unzip

# Expose ports (default)
EXPOSE 5672 15672
