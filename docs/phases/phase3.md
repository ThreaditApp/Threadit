## 🔍 Phase 3 – Functional Requirements and Application Architecture

### 📌 Functional Requirements

#### 🔑 User Management & Authentication
- **FR1.1:** The system should allow users to register a new account by providing a display name, email, and password.
- **FR1.2:** The system should enforce password security policies, including a minimum length and complexity requirements.
- **FR1.3:** The system should allow users to authenticate using their email and password.
- **FR1.4:** The system should issue a token upon successful authentication.
- **FR1.5:** The system should allow users to log out by invalidating their authentication token.
- **FR1.6:** The system should allow authenticated users to delete their account.
- **FR1.7:** The system should ensure that deleted user accounts are permanently removed from the authentication database.

#### 🧑‍💻 User Profiles
- **FR2.1:** The system should allow authenticated users to view their own profile information.
- **FR2.2:** The system should allow users to view another user’s public profile information.
- **FR2.3:** The system should allow authenticated users to update their own profile information, including display name, username, and bio.
- **FR2.4:** The system should enforce uniqueness constraints on usernames.
- **FR2.5:** The system should allow authenticated users to delete their own profile, which should also trigger account deletion.
- **FR2.6:** The system should prevent unauthorized users from modifying another user's profile.

#### 🌍 Communities
- **FR3.1:** The system should allow users to retrieve a list of communities.
- **FR3.2:** The system should allow users to retrieve a list of communities they are members of.
- **FR3.3:** The system should allow users to retrieve a list of communities another user is a member of.
- **FR3.4:** The system should allow users to view the details of a specific community, including its name, description, icon, owner, and creation date.
- **FR3.5:** The system should allow users to create a new community by providing a name and description.
- **FR3.6:** The system should allow community owners to update the community name.
- **FR3.7:** The system should allow community owners to update the community description.
- **FR3.8:** The system should allow community owners to delete the community.
- **FR3.9:** The system should allow community owners to update the community icon by uploading an image file.
- **FR3.10:** The system should allow users to join a community.
- **FR3.11:** The system should allow users to leave a community.
- **FR3.12:** The system should allow users to retrieve a list of members in a specific community.

#### 📝 Threads
- **FR4.1:** The system should allow users to retrieve a list of threads.
- **FR4.2:** The system should allow users to retrieve a list of threads in a specific community.
- **FR4.3:** The system should allow users to retrieve a list of threads they have published.
- **FR4.4:** The system should allow users to retrieve a list of threads they have published in a specific community.
- **FR4.5:** The system should allow users to retrieve a list of threads published by another user.
- **FR4.6:** The system should allow users to retrieve a list of threads published by another user in a specific community.
- **FR4.7:** The system should allow users to view the details of a specific thread, including its title, content, author, creation date, update date, and associated community.
- **FR4.8:** The system should allow users to create a new thread by providing a title, content, and selecting a community.
- **FR4.9:** The system should allow users to update the title of a thread they have published.
- **FR4.10:** The system should allow users to update the content of a thread they have published.
- **FR4.11:** The system should allow users to delete a thread they have published.
- **FR4.12:** The system should allow community owners to delete any thread in their community.

#### 💬 Comments
- **FR5.1:** The system should allow users to retrieve a list of comments for a specific thread.
- **FR5.2:** The system should support pagination for comment retrieval.
- **FR5.3:** The system should allow authenticated users to create a new comment for a thread.
- **FR5.4:** The system should allow users to retrieve the details of a specific comment.
- **FR5.5:** The system should allow authenticated users to update their own comments.
- **FR5.6:** The system should enforce authorization checks to prevent users from modifying or deleting comments they do not own.
- **FR5.7:** The system should allow authenticated users to delete their own comments.
- **FR5.8:** The system should allow users to reply to existing comments.
- **FR5.9:** The system should ensure that replies are correctly linked to their parent comment.

#### ⬆️ Voting System
- **FR6.1:** The system should allow authenticated users to upvote or downvote a thread.
- **FR6.2:** The system should allow authenticated users to upvote or downvote a comment.
- **FR6.3:** The system should allow authenticated users to change or remove their vote on a thread or comment.

#### 🔗 Social Connectivity
- **FR7.1:** The system should allow authenticated users to follow or unfollow another user.
- **FR7.2:** The system should allow authenticated users to follow or unfollow a community.
- **FR7.3:** The system should allow authenticated users to retrieve a list of users and communities they are following.
- **FR7.4:** The system should allow authenticated users to retrieve a list of users following them.
- **FR7.5:** The system should allow authenticated users to retrieve a list of the users following a community.

#### 🏠 Personalized Feed
- **FR8.1:** The system should allow authenticated users to retrieve a personalized feed.
- **FR8.2:** The system should support pagination for feed retrieval.
- **FR8.3:** The system should allow users to sort their feed by newest, trending, or top items.
- **FR8.4:** The system should ensure that feed items contain both threads and comments, categorized accordingly.
- **FR8.5:** The system should provide metadata for each feed item, including community ID, timestamps, and content.
- **FR8.6:** The system should ensure that only authorized users can access personalized feeds.

#### 🔎 Search & Discovery
- **FR9.1:** The system should allow users to search for threads by keyword.
- **FR9.2:** The system should allow users to search for communities by keyword.
- **FR9.3:** The system should allow users to search for other users by keyword.
- **FR9.4:** The system should allow users to search all content by keyword.
- **FR9.5:** The system should support pagination for search results.
- **FR9.6:** The system should allow users to sort search results by date or popularity.

#### 🔗 Microservices Communication
- **FR10.1:** The system should enforce that all inter-service communication between microservices uses gRPC.

### 📦 Application Architecture

![Application architecture](../images/architecture.png)