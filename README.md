# RidesApp

A ride-sharing application that allows users to share rides and request on-demand rides with cash payment options.

## Features

- **User Authentication**: Register, login, and manage user profiles
- **Ride Sharing**: Create and join shared rides with specified seats
- **On-Demand Rides**: Request rides similar to Uber/Lyft with cash payment option
- **Real-time Location Tracking**: Track ride locations in real-time
- **Payment Options**: Support for cash, card, and wallet payments
- **User Ratings**: Rate drivers and passengers after rides

## Tech Stack

### Backend
- Go with Gin framework
- PostgreSQL database with GORM
- JWT authentication
- RESTful API

### Frontend (Planned)
- React Native for mobile app
- React for web app

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login a user

### User Management
- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update current user profile

### Ride Management
- `POST /api/v1/rides` - Create a new ride
- `GET /api/v1/rides/:id` - Get a specific ride
- `GET /api/v1/rides/my` - Get my rides (as rider or driver)
- `GET /api/v1/rides/shared/available` - Get available shared rides
- `GET /api/v1/rides/shared/upcoming` - Get upcoming shared rides
- `PUT /api/v1/rides/:id/status` - Update ride status
- `POST /api/v1/rides/:id/join` - Join a shared ride
- `DELETE /api/v1/rides/:id/passengers/:passengerId` - Leave a shared ride
- `GET /api/v1/rides/:id/passengers` - Get passengers for a ride

## Getting Started

### Prerequisites
- Go 1.16 or higher
- PostgreSQL 12 or higher
- Node.js and npm (for frontend development)

### Installation

1. Clone the repository
```
git clone https://github.com/yourusername/ridesapp.git
cd ridesapp
```

2. Set up the database
```
# Create a PostgreSQL database named 'ridesapp'
# Update the database configuration in backend/pkg/config/config.go
```

3. Run the backend
```
cd backend
go mod download
go run cmd/api/main.go
```

4. Run the frontend (when available)
```
cd frontend
npm install
npm start
```

## Environment Variables

Create a `.env` file in the backend directory with the following variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ridesapp
SERVER_PORT=8080
JWT_SECRET_KEY=your-secret-key
```

## License

This project is licensed under the MIT License - see the LICENSE file for details. 