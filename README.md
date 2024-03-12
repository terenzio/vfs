# Virtual File System (VFS) 
A Virtual File System Implementation for IsCoolLab by Terence Liu 

## Introduction

This project introduces a Virtual File System (VFS), implemented in Go, that emulates a Unix-like environment for managing digital files and directories. Designed with simplicity and efficiency in mind, it facilitates basic operations such as user registration, and creating, updating, deleting, and listing files and folders. Leveraging in-memory data structures for storage, it provides a non-persistent solution that resets upon reboot, offering a streamlined approach for file system management.

## Setup

1) Clone or unzip the repository.
2) Make sure that Golang 1.20+ is installed on your system.
3) Enter into the /cmd directory by using the following command on a Unix-like system:
    ```
    cd cmd
    ```
4) Run the following command to start the VFS program: 
    ```
    go run main.go
    ```

## Available Commands
   ```
   ❯ cd cmd
   ❯ go run main.go

      ==== IsCoolLab: Virtual File System CLI ====
      The current time is: 2024-03-12 03:26:27
      Type 'help' to see available commands.
      
      # help
      Available commands:
      > register [username]
      > create-folder [username] [foldername] [description]?
      > delete-folder [username] [foldername]
      > list-folders [username] [--sort-name|--sort-created] [asc|desc]
      > rename-folder [username] [foldername] [new-folder-name]
      > create-file [username] [foldername] [filename] [description]?
      > delete-file [username] [foldername] [filename]
      > list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]
      > exit
   ```
      

## Sample Usage

   ```
   ❯ cd cmd
   ❯ go run main.go
    
      ==== IsCoolLab: Virtual File System CLI ====
      The current time is: 2024-03-12 03:19:29
      Type 'help' to see available commands.
      
      # register user1
      Add 'user1' successfully.
      
      # create-folder user1 folder1
      Create 'folder1' successfully.
      
      # create-folder user1 folder2 this-is-folder-2
      Create 'folder2' successfully.
      
      # list-folders user1 --sort-name asc
      Name    | Description      | Created At          | User Name
      -------------------------------------------------------------------
      folder1 |                  | 2024-03-12 03:19:50 | user1
      folder2 | this-is-folder-2 | 2024-03-12 03:20:01 | user1
      
      # create-file user1 folder1 config a-config-file
      Create 'config' in user1/folder1 successfully.
      
      # list-files user1 folder1 --sort-name desc
      Name   | Description   | Created At          | Folder  | User Name
      ----------------------------------------------------------------
      config | a-config-file | 2024-03-12 03:20:41 | folder1 | user1
      
      # exit
      Removing file users.txt ...
      Removing file folders.txt ...
      Removing file files.txt ...
      Removed all temp files.
      Exiting program.
      See you next time!

   ```   

## Input Validation
- All input validation is done at the Service Layer, ensuring that the VFS is robust and secure against invalid or malicious inputs.
  - All names (user / folder / file) must contain only alphabets (uppercase and lowercase) and numbers with no spaces.
  - The length of the folder name must be less than or equal to 30 characters.
    ```
    ❯ cd cmd
    ❯ go run main.go
    
        ==== IsCoolLab: Virtual File System CLI ====
        The current time is: 2024-03-12 16:57:25
        Type 'help' to see available commands.
        
        # register User12#$%
        Error: The name [User12#$%] contains invalid chars. Only alphabets and numbers are allowed.
        
        # register user123456789012345678901234567890
        Error: The name [user123456789012345678901234567890] is too long. The maximum length is 30 characters.
    ```

## Unit Tests

- All tests are done on the Service Layer, which contains the core directory logic of the VFS. 
  - The tests are written in the `*_service_test.go` files, and can be run using the following command:
    ```
    ❯ cd service
    ❯ go test -v
    
     === RUN   TestCreateFile
     === RUN   TestCreateFile/ValidFileCreation
     === RUN   TestCreateFile/UserDoesNotExist
     === RUN   TestCreateFile/FolderDoesNotExist
     === RUN   TestCreateFile/InvalidFileName
     --- PASS: TestCreateFile (0.00s)
     --- PASS: TestCreateFile/ValidFileCreation (0.00s)
     --- PASS: TestCreateFile/UserDoesNotExist (0.00s)
     --- PASS: TestCreateFile/FolderDoesNotExist (0.00s)
     --- PASS: TestCreateFile/InvalidFileName (0.00s)
     === RUN   TestCreateFolder
     === RUN   TestCreateFolder/ValidFolderCreation
     === RUN   TestCreateFolder/UserDoesNotExist
     === RUN   TestCreateFolder/InvalidFolderName
     --- PASS: TestCreateFolder (0.00s)
     --- PASS: TestCreateFolder/ValidFolderCreation (0.00s)
     --- PASS: TestCreateFolder/UserDoesNotExist (0.00s)
     --- PASS: TestCreateFolder/InvalidFolderName (0.00s)
     === RUN   TestFolderService
     === RUN   TestFolderService/RenameExistingFolder
     === RUN   TestFolderService/ListFolders
     --- PASS: TestFolderService (0.00s)
     --- PASS: TestFolderService/RenameExistingFolder (0.00s)
     --- PASS: TestFolderService/ListFolders (0.00s)
     === RUN   TestRegister
     === RUN   TestRegister/ErrorInvalidUsername
     === RUN   TestRegister/ErrorUserExists
     === RUN   TestRegister/Success
     --- PASS: TestRegister (0.00s)
     --- PASS: TestRegister/ErrorInvalidUsername (0.00s)
     --- PASS: TestRegister/ErrorUserExists (0.00s)
     --- PASS: TestRegister/Success (0.00s)
     PASS
     ok
    ```

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
