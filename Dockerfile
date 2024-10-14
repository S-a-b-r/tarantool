FROM tarantool/tarantool:3
LABEL authors="shamshurin_gd"

RUN apt-get update && apt-get install -y \
  zip unzip \
  git cmake iproute2
