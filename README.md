# dig-inv
*digital inventory*
### A system for keeping track of assets like domains, servers, and more with as little friction as possible.

TODO: Screenshots / Demo

> **Note:** This is a learning project to explore go, grpc, and other technologies. Don't expect it to be production ready.
> This is pretty overengineered & under engineered for a simple inventory system, but it serves as a good exercise in building a full-stack application with some 
> of the chosen technologies.

## Features

- **Domain Management**: Keep track of your domains, their expiration dates, and associated DNS records.
- **Server Inventory**: Manage your servers, their IP addresses, and associated metadata.
- **Asset Tracking**: Record and manage various assets, including hardware and software.
- **OIDC Authentication**: Secure access to the system using OIDC for authentication.
- **API Access**: Access your inventory programmatically via a gRPC API.
- **Web Interface**: A simple web interface to view and manage your inventory.
- **Exports**: Export your inventory data in various formats for reporting or backup purposes.
- **Helm Chart**: Deploy the application easily on Kubernetes using a Helm chart.

## Supported Providers

- **Custom Providers**
  - Create your own custom asset classes to manage assets that are not supported out of the box.
- **DNS Providers**
  - Cloudflare
  - Namecheap
- **Server Providers**
  - Hetzner



## Quick Start

### docker-compose

TODO

### Helm Chart

TODO

### Standalone

TODO




## Built with

- **DX**
  - **Just**: For task automation and scripting.
  - **Docker**: For containerization and dependency management.
- **Backend**
  - **Go**: The primary language for the backend services.
  - **gRPC**: For efficient communication between services. Using grpc-gateway for HTTP/JSON support.
  - **Nats**: For the worker queue.
- **Frontend**
  - **SvelteKit**: For the web interface.
  - **Typescript**: For type safety in the web interface.
  - **Tailwind CSS**: For styling the web interface.
  - **Flowbite Svelte**: For UI components.
