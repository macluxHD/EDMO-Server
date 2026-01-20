# EDMO IDE

A visual programming environment for controlling EDMO modular robots. This project is part of Project 3-1 (BCS3300) at Maastricht University.

This server acts as a middleman between the browser based interface and the local connection to the robot to relay commands from the browser to the robot.

## Run Locally

### Prerequisites

- [Go](https://go.dev/)
- Git

### Setup Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/macluxHD/EDMO-Server
   ```

2. **Navigate to the project directory**
   ```bash
   cd EDMO-Server
   ```

2. **Copy exampl.env file**
   ```bash
   cp example.env .env
   ```

2. **Edit .env file**

    Set to the corresponding port the robot is connected to on your machine.
    
    On windows you can find the correct COM port number in the device manager.
    
    ```dotenv
    SERIAL_PORT = /COM3
    ```

4. **Run the server**
   ```bash
   go run .
   ```
