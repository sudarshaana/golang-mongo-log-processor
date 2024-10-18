# Use the official MongoDB image from the Docker Hub
FROM mongo:8.0.0-rc11-jammy


# Environment variables for MongoDB
#ENV MONGO_INITDB_ROOT_USERNAME=admin
#ENV MONGO_INITDB_ROOT_PASSWORD=secret

# Expose the default MongoDB port
EXPOSE 27017

# Command to run MongoDB
CMD ["mongod"]
