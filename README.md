Seat Map Application

This project provides a flight seat map application with a Golang backend and a React frontend, orchestrated using Docker Compose. It allows users to view a seat map, select seats, and persists the selected seat information.

How to Run the Application
There are two primary ways to get the application up and running: using Docker (recommended for consistency) or running the backend and frontend manually.

##
**1. Running with Docker (Recommended)**  
This method uses Docker Compose to build and run both the backend and frontend services in isolated containers, ensuring all dependencies are handled automatically.

- Prerequisites:
  - Docker Desktop (which includes Docker Engine and Docker Compose) installed and running on your system.

- Steps:
  - Navigate to the project root directory:
    - Open your terminal or command prompt and change your current directory to the seat-app folder (the one containing backend/, frontend/, and docker-compose.yml).
      
      ```bash
      cd seat-app/
      ```
  
  - Build and run the Docker containers:
    - Execute the following command. The --build flag is important to compile the Go backend and build the React frontend image.
      
      ```bash
      docker compose up --build
      ```

- This command will:
  - Build the Golang backend image.
  - Build the React frontend image and set up Nginx to serve it.
  - Start both services.

- Access the application:
  - Once the services are running, open your web browser and navigate to:
    
    ```bash
    http://localhost:3000
    ```

##
**2. Running Manually (Backend and Frontend Separately)**  
This method involves running the Golang backend and React frontend independently on your local machine.

- Prerequisites:
  - For Backend: Go (version 1.23.10 or higher recommended) installed.
  - For Frontend: Node.js (LTS version recommended) and npm installed.

- Steps:
  - Start the Backend:
    - Navigate to the backend directory:
      
      ```bash
      cd seat-app/backend 
      ```
    - Run the Go application:
    
      ```bash
      go run main.go 
      ```

    - The backend API will start and listen on port **8080**. You can try this example API.
      ```bash
      http://localhost:8080/seatmap
      ```
    
  - Start the Frontend:

    - Navigate to the frontend project directory:
      
      ```bash
      cd seat-app/frontend
      ```

    - Install Node.js dependencies _(if you haven't already)_:  
      This command will install all necessary node_modules based on package.json. You only need to run this once or if dependencies change.
        
      ```bash
      npm install
      ```

    - Start the React development server:
      ```bash
      npm start
      ```
      
    - This will compile the React application and usually open your default web browser to 
      ```bash
      http://localhost:3000
      ```

    - Or access manually, by accessing below link to your browser:
      
      ```bash
      http://localhost:3000
      ```
##
- Data Persistence  
  The application's seat selection data is stored and updated directly within the **SeatMapResponse.json** file.  
  When running via Docker, this file is bind-mounted from your local backend/ directory into the container, meaning any changes made through the application (e.g., selecting a seat) will persist in your local SeatMapResponse.json file even if the Docker containers are stopped and restarted.
