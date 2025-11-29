# 1. Use the official GenAI Toolbox image as the base
FROM us-central1-docker.pkg.dev/database-toolbox/toolbox/toolbox:latest

# 2. Set the working directory inside the container
WORKDIR /app

# 3. Copy your local tools.yaml file into the container
#    (We rename it to config.yaml inside, but you can keep the name tools.yaml)
COPY tools.yaml /app/tools.yaml

# 4. Expose the default port (5000)
EXPOSE 8080

# 5. Define the command to run the toolbox pointing to your config file
#    The binary is globally available as 'toolbox' in the base image.
ENTRYPOINT ["/toolbox"]
CMD ["--tools-file", "/app/tools.yaml", "--address", "0.0.0.0", "--port", "8080"]
