# dig-inv
*digital inventory*
### A system for keeping track of assets like domains, servers, and more with as little friction as possible.

TODO: Screenshots / Demo

> **Note:** This is a learning project to explore go, grpc, and other technologies. Don't expect it to be production ready.
> This is pretty overengineered for a simple inventory system, but it serves as a good exercise in building a full-stack application with some 
> of the chosen technologies.

## Features

- **Domain Management**: Keep track of your domains, their expiration dates, and associated DNS records.
- **Server Inventory**: Manage your servers, their IP addresses, and associated metadata.
- **Asset Tracking**: Record and manage various assets, including hardware and software.
- **User Management**: Create and manage users with different roles and permissions.
- **API Access**: Access your inventory programmatically via a gRPC API.
- **Web Interface**: A simple web interface to view and manage your inventory.
- **Exports**: Export your inventory data in various formats for reporting or backup purposes.

## Supported Providers
- **DNS Providers**
    - Cloudflare
    - Namecheap
- **Server Providers**
    - Hetzner
- **Asset Providers**
    - Custom asset management


## Built with
- **Go**: The primary language for the backend services.
- **gRPC**: For efficient communication between services.
- **Just**: For task automation and scripting.
- **Docker**: For containerization and development.
- **SvelteKit**: For the web interface.
- **Typescript**: For type safety in the web interface.
- **Tailwind CSS**: For styling the web interface.
