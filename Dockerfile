FROM ubuntu:rolling
RUN apt-get update && apt-get install -y \
    htop \
    tmux \
    nano \
 && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /gotag/
ADD gotag /gotag/
RUN chmod +x /gotag
ADD frontend/dist /gotag/frontend
ADD migrations /gotag/migrations
ENTRYPOINT ["/gotag/gotag","serve"]
EXPOSE 4000