# Statement FROM digunakan untuk inisialisasi base image
# golang:alpine => image golang dengan tag alpine
# alpine adalah linux distribution image yang minimalist dan lightweight 
FROM golang:alpine

# Statement WORKDIR akan menentukan working direktori pada direktori app
# Statement ini membuat statement RUN dibawah akan dieksekusi pada working direktory yang telah ditentukan.
WORKDIR /app

# Statement COPY digunakan untuk mengcopy file dari source [argument peratama] ke destination [argument kedua]
COPY go.mod go.sum ./
COPY . .

RUN go mod download

# Build aplikasi golang menjadi binary dengan nama go-learn-mongodb
RUN CGO_ENABLED=0 GOOS=linux go build -o ./cmd/go-learn-mongodb ./cmd/main.go

# Statement EXPOSE digunakan untuk expose port container
EXPOSE 9999

# Run
CMD ["./cmd/go-learn-mongodb"]