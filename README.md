# Virtual File System (VFS) 
A Virtual File System Implementation for IsCoolLab by Terence Liu 

## Introduction

This project introduces a Virtual File System (VFS), implemented in Go, that emulates a Unix-like environment for managing digital files and directories. Designed with simplicity and efficiency in mind, it facilitates basic operations such as user registration, and creating, updating, deleting, and listing files and folders. Leveraging in-memory data structures for storage, it provides a non-persistent solution that resets upon reboot, offering a streamlined approach for file system management.

## Design Principles

The VFS is grounded in several core design principles aimed at enhancing its modularity, extensibility, and overall user experience:

### 1. Domain-Driven Design (DDD)

- **Modularity & Extensibility**: The system architecture promotes a modular design, allowing for easy extension and modification of its capabilities.
- **Testability & Maintainability**: Emphasizes ease of testing and maintenance, ensuring long-term reliability and ease of updates.
- **Performance & Scalability**: Engineered for high performance and scalability to accommodate growth in user demand and data volume.
- **User-Friendliness**: Focuses on delivering an intuitive and seamless user experience, mimicking familiar Unix-like file operations.

### 2. Structure of the File System

The VFS architecture comprises three primary models: User, File, and Folder. Each is represented as a Go struct, with associated repository and service layers facilitating data management and business logic, respectively.

#### 2.1 Repository Layer

- **Data Storage & Retrieval**: Handles the storage, retrieval, and management of model data, abstracted through interfaces to support diverse storage mechanisms.
- **Flexibility in Storage**: Implemented primarily as an in-memory store, with temporary text files serving as the storage medium, ensuring rapid access and manipulation of file system data.
- **Storage Mechanism Independence**: The interface-driven design permits easy substitution of storage backends, enhancing the system's adaptability to future storage requirements.

#### 2.2 Service Layer

- **Business Logic**: Encapsulates the core business logic of each model, abstracting the complexities of data storage and retrieval from the application layer.
- **Repository Interface Interaction**: The service layer interacts exclusively with the repository interface, decoupling business logic from specific storage implementations.
- **Domain & Storage Separation**: This distinct separation ensures that domain logic remains unaffected by changes in the storage layer, facilitating seamless transitions to alternative storage solutions.

## Conclusion

The VFS project exemplifies a thoughtful application of DDD principles to create a virtual file management system that is both robust and user-centric. Its architecture not only prioritizes performance and ease of use but also ensures that the system is well-positioned for future expansion and adaptation to new technologies or storage paradigms.




    