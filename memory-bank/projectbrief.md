# Project Brief

**dockertestx** is a Go testing library that simplifies integration testing with databases and other services using Docker containers. It is built on top of [dockertest](https://github.com/ory/dockertest) and was originally forked from [sqltest](https://github.com/vvatanabe/sqltest). While initially focused on SQL databases (MySQL, PostgreSQL), dockertestx has expanded to support a broader range of data stores and services.

## Objectives

1. Abstract the complexity of container-based integration testing and provide a simple, user-friendly API.  
2. Expand support for various data stores and services.  
3. Automate test environment setup to improve developer productivity.  
4. Ensure a consistent pattern for writing test code.  

## Currently Supported Services

- **SQL Databases**: MySQL, PostgreSQL  
- **Cache Services**: Memcached, Redis  
- **Object Storage**: MinIO (S3-compatible)  
- **NoSQL Databases**: DynamoDB Local  

## Future Support Plans

- MongoDB  
- Additional data stores  
- Enhanced extensibility for custom service containers  

## Core Features

1. Container lifecycle management (start, stop, cleanup)  
2. Health checks and connection establishment  
3. Helper functions for test data preparation and validation  
4. Flexible customization options  

## Project Maintainer

- **[vvatanabe](https://github.com/vvatanabe/)** - Lead Developer  

## License

Apache License 2.0
