# Code Mafia - Code Execution Platform

This project utilizes Judge0 via the RapidAPI Cloud Endpoint for scalable and resilient code execution.

## Getting Started

### Configuration

The backend is configured to submit code securely to the Judge0 Cloud API. Provide the RapidAPI credentials in your `.env` file:

```env
RAPIDAPI_URL=https://judge0-ce.p.rapidapi.com
RAPIDAPI_HOST=judge0-ce.p.rapidapi.com
RAPIDAPI_KEY=your_rapidapi_key_here
```

### Starting the Infrastructure

We utilize `docker-compose` to spin up the local Backend, MongoDB, and Redis components smoothly. Run:
```bash
docker-compose up -d
```

Your system is now successfully set to execute code efficiently via Judge0 CE Cloud endpoints!
