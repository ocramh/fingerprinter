FROM golang:1.15-alpine
RUN apk update
RUN apk add --no-cache --update \
  build-base \
  nasm \
  tar \
  bzip2 \
  bash \
  cmake \
  make \
  g++ \
  curl \
  yasm \
  pkgconfig \
  git \
  ffmpeg-dev \
  fftw-dev \
  ffmpeg

ARG CHROMAPRINT_VERSION=1.4.3
ARG CHROMAPRINT_URL=https://github.com/acoustid/chromaprint/archive/v${CHROMAPRINT_VERSION}.tar.gz

WORKDIR /tmp
RUN mkdir chromaprint && \
    curl -# -L ${CHROMAPRINT_URL} | tar xz --strip 1 -C chromaprint && \
    cd chromaprint && \
    cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_TOOLS=ON . && \
    make && make install

WORKDIR /app/fingerprinter
COPY . .
RUN make install

CMD ["sh"]
