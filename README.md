# Project
Cloud Computing Project

## Phase 1 - Datasets, business capabilities and use cases

### Datasets
- ![Reddit Top 2.5 Million](https://github.com/umbrae/reddit-top-2.5-million) (1.66GB) (Aug. 2013)
- ![Twitter Friends](https://www.kaggle.com/datasets/hwassner/TwitterFriends) (448.48MB) (Sep. 2016)

### Business capabilities
- User Management
  - Authentication (email+password, Google and Apple)
  - Profile customization (avatar, bio, preferences, etc.)
  - Karma system
- Subreddit Management
  - Creation of subreddits
  - Subreddit moderation (administration, rules, bans, etc.)
  - Roles (moderators, admins, members, etc.)
  - Privacy settings (public and private)
- Content Management
  - Posts (text, image, video, links, etc.)
  - Comments
  - Upvotes and downvotes
  - Editing and deleting content
  - NSFW and spoiler tags
- Engagement & Interaction
  - Notifications
  - Direct messages
  - Follows (users and subreddits)
  - Rewards and badges
- Search & Discovery
  - Global and subreddit search
  - Filtering (relevance, date, popularity, etc.)
  - Tags
  - Personalized recommendations
- Monetization & Premium Features
  - Advertisements
  - Sponsored posts
  - Premium subscription plans (remove ads, exclusive themes, etc.)
  - Virtual currency (rewards and special features)
- Moderation & Safety
  - Reports (post and comments)
  - Automated moderation (rules, etc.)
  - Filtering settings (keywords, nfsw, etc.)
  - Bans

### Use cases

#### UC-1: User registration
**Actor:** Unregistered user

**Preconditions:** The user is not logged in and has not registered before.

**Flow:**
1. The user navigates to the registration page.
2. The user chooses a registration method:
    - Provides a unique username, email, and password.
    - Uses Google account authentication.
    - Uses Apple account authentication.
3. The system validates the provided details.
4. The system stores the user details securely.
5. The system confirms registration and logs in the user.

**Exceptions:**
- If the username or email is already in use, an error message is displayed.
- If validation fails (e.g., weak password), an error message is displayed.
- If third-party authentication fails, an error message is displayed.

#### UC-2: User login
**Actor:** Registered user

**Preconditions:** The user has an existing account.

**Flow:**
1. The user navigates to the login page.
2. The user chooses a login method:
    - Enters their email/username and password.
    - Uses Google account authentication.
    - Uses Apple account authentication.
3. The system verifies the credentials.
4. If valid, the system logs in the user and redirects to the homepage.

**Exceptions:**
- If the username or email is already in use, an error message is displayed.
- If validation fails (e.g., weak password), an error message is displayed.
- If third-party authentication fails, an error message is displayed.

#### UC-3: Create a post
**Actor:** Registered user

**Preconditions:** The user is logged in.

**Flow:**
1. The user clicks the "Create Post" button.
2. The user selects a subreddit.
3. The user selects a post type (text, image, link, poll, etc.).
4. The user enters the post details (title, content, attachments, etc.).
5. The user submits the post.
6. The system validates and stores the post.
7. If valid, the system displays the created post.

**Exceptions:**
- If validation fails (e.g. required fields missing, banned from subreddit), an error message is displayed.

#### UC-4: Edit a post
**Actor:** Registered user (Post owner)

**Preconditions:** The user is logged in and owns the post.

**Flow:**
1. The user navigates to their post.
2. The user clicks the "Edit Post" button.
3. The user modifies the post content.
4. The user submits the changes.
5. The system validates and updates thhe post.
6. If valid, the system displays the updated post.

**Exceptions:**
- If the user is not the owner of the post, an error message is displayed.
- If validation fails (e.g. required fields missing), an error message is displayed.

#### UC-5: Delete a post
**Actor:** Registered user (Post owner)

**Preconditions:** The user is logged in and owns the post.

**Flow:**
1. The user navigates to their post.
2. The user clicks the "Delete Post" button.
3. The system prompts for confirmation.
4. The user confirms deletion.
5. The system removes the post details (content, attachments, etc.).

**Exceptions:**
- If the user is not the owner of the post, an error message is displayed.

#### UC-6: Comment on a post
**Actor:** Registered user

**Preconditions:** The user is logged in.

**Flow:**
1. The user navigates to a post.
2. The user enters a comment.
3. The user submits the comment.
4. The system validates and stores the comment.
5. The system displays the comment under the post.

**Exceptions:**
- If the user lacks permissions, an error message is displayed.
- If validation fails (e.g. required fields missing), an error message is displayed.

#### UC-7: Vote on a post or comment
**Actor:** Registered user

**Preconditions:** The user is logged in.

**Flow:**
1. The user navigates to a post or comment.
2. The user clicks the upvote or downvote button.
3. The system registers the vote.
4. The system updates the post or comment score.

**Exceptions:**
- If the user has already voted, the system toggles or removes the vote.

#### UC-8: Report a post or comment
**Actor:** Registered user

**Preconditions:** The user is logged in.

**Flow:**
1. TODO

**Exceptions:**
- TODO

#### UC-9: Join a subreddit
**Actor:** Registered user

**Preconditions:** The user is logged in.

**Flow:**
1. TODO

**Exceptions:**
- TODO

#### UC-10: Create a subreddit
**Actor:** Registered user

**Preconditions:** The user is logged in.

**Flow:**
1. TODO

**Exceptions:**
- TODO

#### UC-11. Manage subreddit join request
**Actor:** Registered user (Subreddit moderator)

**Preconditions:** The user is logged in and is a moderator of the subreddit.

**Flow:**
1. TODO

**Exceptions:**
- TODO

#### UC-12: Ban a user
**Actor:** Registered user (Subreddit moderator)

**Preconditions:** The user is logged in and is a moderator of the subreddit.

**Flow:**
1. TODO

**Exceptions:**
- TODO

#### UC-13: Remove another user's post
**Actor:** Registered user (Subreddit moderator)

**Preconditions:** The user is logged in and is a moderator of the subreddit.

**Flow:**
1. TODO

**Exceptions:**
- TODO
