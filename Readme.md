# Quiz App

This is a simple quiz application that allows users to participate in quizzes, view results, and manage quizzes.

## Features
- Create, edit, and delete quizzes.
- Participate in quizzes and get scores.
- Admin panel to manage quizzes.
  
## Tech Stack
- Node.js (Backend)
- Express.js (Framework)
- PostgreSQL (Database)
- Sequelize (ORM)
- Docker (Containerization)
- Docker Compose (Service Orchestration)

## Prerequisites
- Docker and Docker Compose must be installed on your system.

## Getting Started

### 1. Clone the Repository
```bash
git clone https://github.com/shaileshhb/quiz-app.git
cd quiz-app
```

### 2. Set Up Environment Variable
```bash
# src/.env file

# JWT secret
JWT_SECRET=your_jwt_secret
```

Replace `your_jwt_secret` with actual values.

### 3. Run the Service Using Docker Compose
```bash
docker-compose up --build
```

This command will:

- Build the Docker image for the application.
- Run the quiz app on http://localhost:8080.


## Endpoints

### 1. User Registration
**POST** `/api/v1/register`

Registers a new user.

**Body Parameters:**
- `name` (string): Full name of the user.
- `email` (string): Email address of the user.
- `password` (string): Password for the account.

**Response:**
```json
{
  "id": "1d71461b-7fed-4bb7-8395-ed24f908aaa1",
  "name": "shailesh",
  "username": "shailesh",
  "token": "your_jwt_token"
}
```

### 2. User Login
**POST** `/api/v1/login`

Logs in an existing user.

**Body Parameters:**
- `email` (string): User's email.
- `password` (string): User's password.

**Response:**
```json
{
  "id": "1d71461b-7fed-4bb7-8395-ed24f908aaa1",
  "name": "shailesh",
  "username": "shailesh",
  "token": "your_jwt_token"
}
```

---
## Quizzes
### 3. Create a Quiz
**POST** `/api/v1/quizzes`

**Headers**: Requires `Authorization: Bearer <token>`

**Body Parameters:**
- `title` (string): Title of the quiz.
- `maxTime` (int): Maximum time in minutes for which quiz will be valid. Default is 2 minutes.
- `questions` (array): An array of questions with choices and correct answers.
  - `text` (string): Question text
  - `options` (array): An array of options with one correct answer. Each must contain 4 options
    - `answer` (string): Specifies the answer
    - `isCorrect` (boolean): Whether the answer is correct
Example:
```json
{
  "title": "Quiz 2",
  "questions": [{
    "text": "This is question 1",
    "options": [{
      "answer": "Answer 1",
      "isCorrect": false
    },{
      "answer": "Answer 2",
      "isCorrect": false
    },{
      "answer": "Answer 3",
      "isCorrect": false
    },{
      "answer": "Answer 4",
      "isCorrect": true
    }]
  }]
}
```

**Response:**
```json
{
  "quizID": "4e62ebdd-38df-4882-8528-f41f1ef45b3f"
}
```

### 4. Get Quiz Details
**GET** `/api/v1/quizzes/:quizID`

Fetches details of a single quiz by its ID.

**Headers**: Requires `Authorization: Bearer <token>`

**URL Parameters:**
- `quizID` (uuid): ID of the quiz.

**Response:**
```json
{
  "id": "997f06f9-89d1-4f95-9300-09caee4d6b40",
  "title": "Sample Quiz",
  "maxTime": 1,
  "questions": [
    {
      "id": "a06217ee-5688-4a24-b752-7b441985b91e",
      "text": "What is the capital of France?",
      "quizID": "997f06f9-89d1-4f95-9300-09caee4d6b40",
      "options": [
        {
          "id": "fd13c725-c295-4a02-b463-8e7fbdfccf0b",
          "questionID": "a06217ee-5688-4a24-b752-7b441985b91e",
          "answer": "Paris"
        },
        {
          "id": "89bd797c-3d9d-4162-8f8a-33c5d3d5b976",
          "questionID": "a06217ee-5688-4a24-b752-7b441985b91e",
          "answer": "Delhi"
        },
        {
          "id": "88364273-7954-440b-83b0-61b242fda467",
          "questionID": "a06217ee-5688-4a24-b752-7b441985b91e",
          "answer": "London"
        },
        {
          "id": "9e57d44f-231b-4f1b-9a5e-89e986ce2dc6",
          "questionID": "a06217ee-5688-4a24-b752-7b441985b91e",
          "answer": "New York"
        }
      ]
    }
  ]
}
```

---

## Quiz Participation
### 5. Start Quiz
**POST** `/api/v1/users/quizzes/:quizID/start`

Start specifed quiz for the logged in user

**Headers**: Requires `Authorization: Bearer <token>`

**URL Parameters:**
- `quizID` (uuid): ID of the quiz.

**Response:**
```json
{
  "id": "a8448fc1-bf25-4903-9b9e-50579be43c35",
  "userID": "46370596-1b06-4517-94a8-0fae388de28c",
  "quizID": "997f06f9-89d1-4f95-9300-09caee4d6b40",
  "startedAt": "2024-09-29T01:21:07.2553434+05:30",
  "endAt": null,
  "totalScore": 0,
  "userResponses": null
}
```

### 6. Submit Answers
**POST** `/api/v1/users/quizzes/:quizID/attempt/:userQuizAttemptID`

**Headers**: Requires `Authorization: Bearer <token>`


### 7. Get Quiz Results
**GET** `/api/v1/users/quizzes/:quizID/results`

**Headers**: Requires `Authorization: Bearer <token>`

**Response:**
```json
{
  "id": "bc26d845-f408-4b69-b389-87f38bb531c3",
  "userID": "898f9cd8-9b77-412d-ae52-b3bff0285cbb",
  "quizID": "997f06f9-89d1-4f95-9300-09caee4d6b40",
  "startedAt": "2024-09-29T01:24:05.872295+05:30",
  "endAt": null,
  "totalScore": 1,
  "userResponses": [
    {
      "id": "ed9c0222-f9ff-4074-9b0e-ad6d0f28f2c4",
      "userID": "898f9cd8-9b77-412d-ae52-b3bff0285cbb",
      "quizID": "997f06f9-89d1-4f95-9300-09caee4d6b40",
      "userQuizAttemptID": "bc26d845-f408-4b69-b389-87f38bb531c3",
      "questionID": "0724986d-2683-466f-a672-e7eab9ad7ce0",
      "selectedOptionID": "1d3b8e64-7750-4fac-b0fe-53634c05d530",
      "isCorrect": true
    }
  ]
}
```