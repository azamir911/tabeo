# Use a minimal Debian-based image
FROM debian:buster-slim

# Set the working directory
WORKDIR /app

# Install necessary tools
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

# Copy the pre-built 'main' executable from the host system to the container
COPY main /app/main

# Ensure the executable has the correct permissions
RUN chmod +x /app/main

# Verify the contents of /app to make sure 'main' was copied correctly
RUN echo "Verifying if 'main' was copied:" && ls -la /app

# Expose port 8080
EXPOSE 8080

# Run the executable
CMD ["/app/main"]
