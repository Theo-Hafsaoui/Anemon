FROM ghcr.io/theo-hafsaoui/go-latex-img:c800f95

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build

CMD ["./Anemon", "-g"]
