# PlexNotify

PlexNotify is a webhook client that listens for events from a Plex media server. Once a payload is recieved the application sends notifications to a Discord Channel using a JSON payload in a POST request.

# Set Up

## Environment

- GO 1.14+

### Dependecies

github.com/gorilla/mux

### Environment Variables
Discord Webhook URL
Username and Password

**Example**
`export DISCORDURL=https://discordapp.com/api/webhooks/Token/Key`
`export UNAME=admin`
`export PWORD=password`


	
	# Set environment variables
	ENV DISCORDURL=https://discordapp.com/api/webhooks/Token/Key
	ENV UNAME=admin
	ENV PWORD=password
	

## Plex

- Plex Server with Premium Pass to send webhooks

### How to enable webhooks

https://support.plex.tv/articles/115002267687-webhooks/

## Discord

### Webhooks

https://support.discordapp.com/hc/en-us/articles/228383668-Intro-to-Webhooks

## Docker

### Example Dockerfile
```
# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Add the current package to Go 
ADD . /go/src/github.com/ErwinsExpertise/PlexNotify

# Get needed dependcies
RUN go get github.com/gorilla/mux

# Test the application
RUN go test handlers/*

# Build the Go app
RUN go build -o notify .

# Expose port 9000 to the outside world
EXPOSE 9000

# Set environment variables
ENV DISCORDURL=https://discordapp.com/api/webhooks/Token/Key
ENV UNAME=admin
ENV PWORD=password

# Command to run the executable
CMD ["./notify", "-port", "9000"]
```



# To Do

- [x] Accept incoming webhooks 
- [x] Create an acitvity page
- [x] Secure activity page with PP
- [ ] Add search feature for activity
- [ ] Create stand-alone login page
- [ ] Implement option to use HTTPS with self-signed SSL
- [ ] Add page for most played media
- [ ] Implement better handling of non-media related payloads
- [ ] Create test cases for all functions