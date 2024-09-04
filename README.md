# Message Queue Service

## Overview

This project implements a simple message queue service using Go. It provides endpoints for subscribing, unsubscribing, publishing messages, and retrieving messages.

## Endpoints

### 1. Subscribe

**Endpoint:** `POST /v1/subscribe`

**Description:** Allows a subscriber to subscribe to one or more topics.

**Request Parameters:**

```json
{
  "endpoint": "http://subscriber-endpoint.com",
  "topics": ["topic1", "topic2"]
}
```

**Request Body:**

- `endpoint` (string): The URL where the subscriber is located.
- `topics` (array of strings): List of topics to subscribe to.

**Response:**

- **Success:**

  ```json
  {
    "data": {
      "subscriber_id": "unique-subscriber-id"
    },
    "message": "Subscribed successfully"
  }
  ```

- **Error:**

  ```json
  {
    "error": "error-message"
  }
  ```

### 2. Unsubscribe

**Endpoint:** `POST /v1/unsubscribe`

**Description:** Allows a subscriber to unsubscribe from one or more topics.

**Request Parameters:**

```json
{
  "subscriber_id": "unique-subscriber-id",
  "topics": ["topic1", "topic2"]
}
```

**Request Body:**

- `subscriber_id` (string): The ID of the subscriber to be removed.
- `topics` (array of strings): List of topics to unsubscribe from.

**Response:**

- **Success:**

  ```json
  {
    "message": "Unsubscribed successfully from the specified topics"
  }
  ```

- **Error:**

  ```json
  {
    "error": "error-message"
  }
  ```

### 3. Publish

**Endpoint:** `POST /v1/publish`

**Description:** Publishes a message to a specified topic.

**Request Parameters:**

```json
{
  "topic": "topic1",
  "message": "This is a test message"
}
```

**Request Body:**

- `topic` (string): The topic to which the message will be published.
- `message` (string): The content of the message.

**Response:**

- **Success:**

  ```json
  {
    "message": "Message published successfully to all topic subscribers"
  }
  ```

- **Error:**

  ```json
  {
    "error": "error-message"
  }
  ```

### 4. Retrieve

**Endpoint:** `GET /v1/retrieve`

**Description:** Retrieves messages for a specific subscriber from a specified topic starting from a given date.

**Request Parameters:**

```json
{
  "topic": "topic1",
  "subId": "unique-subscriber-id",
  "startDate": "2024-01-01T00:00:00Z"
}
```

**Query Parameters:**

- `topic` (string): The topic from which messages will be retrieved.
- `subId` (string): The ID of the subscriber requesting the messages.
- `startDate` (string): The start date in RFC3339 format to filter messages from.

**Response:**

- **Success:**

  ```json
  {
    "data": [
      {
        "ID": "message-id",
        "Body": "Message content",
        "PUBLISHED_AT": "2024-01-01T00:00:00Z",
        "CONSUMED_AT": "2024-01-01T01:00:00Z"
      }
    ]
  }
  ```

- **Error:**

  ```json
  {
    "error": "error-message"
  }
  ```

## Running the Server

1. **Clone the repository:**

   ```sh
   git clone https://github.com/your-repo/Message-Queue.git
   ```

2. **Navigate to the project directory:**

   ```sh
   cd Message-Queue
   ```

3. **Build and run the application:**

   ```sh
   go run .
   ```

4. **The server will start at port 3000.**

