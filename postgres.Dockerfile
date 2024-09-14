FROM postgres:latest

# Cài đặt các phụ thuộc cần thiết
RUN apt-get update && \
    apt-get install -y git make gcc postgresql-server-dev-all

# Cài đặt wal2json
RUN git clone https://github.com/eulerto/wal2json.git && \
    cd wal2json && \
    make && \
    make install

# Clean up
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
