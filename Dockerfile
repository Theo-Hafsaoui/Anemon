FROM golang:1.23 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build
RUN ./anemon generate

FROM debian:latest
COPY --from=build /app/assets/latex/output/ /internal_output
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    texlive \
    texlive-latex-extra \
    texlive-fonts-extra \
    texlive-xetex \
    texlive-font-utils \
    fonts-font-awesome \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

CMD mkdir -p /tmp_output && cd /internal_output && for i in *.tex; do pdflatex $i -output-directory=/tmp_output || true; done && ls && pwd && ls /tmp_output && cp /internal_output/*.pdf /output/
