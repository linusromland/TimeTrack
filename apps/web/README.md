# TimeTrack Web UI

Modern React 19 web application with Material-UI 9 for the TimeTrack time tracking system.

## Features

- **Modern React 19** with TypeScript for type safety
- **Material-UI 9** for consistent, accessible UI components
- **Responsive design** supporting mobile, tablet, and desktop
- **HTTP-only cookie authentication** for security
- **Real-time statistics** and charts using MUI X Charts
- **Date/time pickers** with timezone handling
- **Project and time entry management**
- **Jira integration** support

## Development

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Running TimeTrack API server on `localhost:8080`

### Getting Started

1. Install dependencies:
```bash
npm install
```

2. Start development server:
```bash
npm run dev
```

3. Open [http://localhost:3000](http://localhost:3000) in your browser

The development server includes:
- Hot reloading for instant updates
- API proxy to `localhost:8080` for backend calls
- TypeScript type checking
- ESLint for code quality

### Build for Production

```bash
npm run build
```

Build output will be in the `dist/` directory, ready to be served by the Go API server.

### Linting

```bash
npm run lint
```

## Project Structure

```
src/
├── components/       # Reusable UI components
├── contexts/         # React contexts (auth, theme)
├── hooks/           # Custom React hooks
├── pages/           # Main application pages
├── services/        # API service layer
├── types/           # TypeScript type definitions
├── utils/           # Utility functions
├── App.tsx          # Main app component
└── main.tsx         # App entry point
```

## API Integration

The web app communicates with the TimeTrack Go API using:

- **Authentication**: HTTP-only cookies with JWT tokens
- **CRUD operations**: Full support for projects and time entries
- **Statistics**: Charts and analytics for time tracking data
- **Real-time updates**: Manual refresh with refresh buttons

## Authentication Flow

1. User logs in via login form
2. API sets HTTP-only cookie with JWT token
3. Subsequent API calls include cookie automatically
4. Protected routes redirect to login if not authenticated
5. Logout clears the authentication cookie

## Responsive Design

The application is designed to work seamlessly across:

- **Mobile** (320px+): Touch-friendly interface with drawer navigation
- **Tablet** (768px+): Optimized layouts with sidebar navigation
- **Desktop** (1024px+): Full-featured interface with persistent sidebar

## Security Features

- HTTP-only cookies prevent XSS attacks on tokens
- CSRF protection through SameSite cookie attributes
- Input validation and sanitization
- Secure API communication

## Contributing

1. Follow TypeScript strict mode guidelines
2. Use Material-UI components consistently
3. Write responsive CSS using MUI breakpoints
4. Add proper error handling and loading states
5. Test across different screen sizes