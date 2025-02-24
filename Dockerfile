FROM golang:1.24.0

# Create a new user in the docker image
RUN adduser --disabled-password --gecos '' gouser

# Create a new directory for goserve files and set the path in the container
RUN mkdir -p /app/kai-app

# Set the working directory inside the container
WORKDIR /app/kai-app

# Copy the Go application into the container
COPY . .

RUN chown -R gouser:gouser /app/kai-app

# Switch to the gouser user
USER gouser

# Install the dependencies
RUN go mod download

# Build the Go application
RUN go build -o app cmd/main.go

# Step 6: Expose the port your app will run on
EXPOSE 8000

# Command to run the app
CMD ["./app"]
