# Social Network — Full Stack Social Media Platform

A modern, responsive, and fully-featured social media platform built with a high-performance Go backend and a Next.js (React) frontend. The system is containerized with Docker for easy orchestration and deployment.

---

## 🚀 Features

### 1. Authentication & Session Management
- **Security:** Standard email/password registration and login with encrypted credentials.
- **Session Control:** Secure, cookie-based session verification middleware.
- **OAuth:** Direct integration with Google Login/OAuth.

### 2. User Profiles
- **Privacy:** Toggle account between Public and Private status.
- **Information:** Customizable profile details: First name, Last name, Nickname, Date of birth, About me.
- **Avatars:** High-quality avatar image uploads.
- **Stats:** Live count display for Posts, Followers, and Following.

### 3. Follower Network
- **Bypass Flow:** Auto-follow public users directly without confirmation.
- **Requests:** Send, accept, or decline follow requests for private profiles.
- **Manage Feed:** Unfollow users or view the full Followers/Following list under the Network dashboard.

### 4. Posts & Feed
- **Create Post:** Post text messages with support for images and video uploads.
- **Privacy Levels:** Post visibility can be set to:
  - **Public:** Visible to all logged-in users.
  - **Followers Only (Almost Private):** Visible only to accepted followers.
  - **Selected Followers (Private):** Fine-grained visibility toggle allowing only specific followers to view the post.
- **Feed:** Real-time feed fetching posts based on following status and visibility filters.
- **Interactions:** Fully functional comment sections for each post.

### 5. Groups
- **Create Group:** Start new groups with title and description.
- **Join requests:** Join groups directly, or send join requests that require approval.
- **Invites:** Group members can invite followers to join the group.
- **Group Posts:** Share group-specific posts.
- **Group Events:** Create events with details (date/time) and RSVP (Going vs. Not Going).
- **Group Chat:** Real-time group messaging and chat history.

### 6. Private Messaging (Real-Time)
- **Live Connection:** WebSocket-based messaging hub for immediate communication.
- **History:** Automatically loads previous chat history on chat initialization.
- **Emoji Picker:** Integrated emoji shortcut panel for rich chat messages.
- **Conversations List:** Access and search ongoing conversations from the messaging sidebar.

### 7. Global Notifications
- **Status Indicator:** Persistent header bell icon displaying live unread notification counts.
- **Interactive List:** Mark notifications as read or view details directly.

---

## 🛠️ Tech Stack

- **Backend:** Go (Golang) 1.25, SQLite3 (CGO_ENABLED=1), Go Migrate.
- **Frontend:** Next.js 16 (React 19), Redux Toolkit, Tailwind CSS, Lucide React.
- **Containerization:** Docker, Docker Compose.

---

## 📦 How to Run

### Option A: Using Docker (Recommended)
Make sure you have Docker and Docker Compose installed, then run the helper script:
```bash
./start.sh
```
This will build and spin up the containers:
- **Frontend:** [http://localhost:3000](http://localhost:3000)
- **Backend API:** [http://localhost:8080](http://localhost:8080)

### Option B: Local Manual Development
If you prefer running directly on your system:

#### 1. Backend
Ensure you have Go installed.
```bash
cd backend
go run cmd/api/main.go
```
The server will run on port `8080`.

#### 2. Frontend
Ensure you have Node.js installed.
```bash
cd front-end
npm install
npm run dev
```
Open [http://localhost:3000](http://localhost:3000) in your browser.