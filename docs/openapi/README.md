## 🌐 API Description

#### `GET /communities`

Retrieves a list of communities. Supports optional filtering and pagination.

**Query Parameters**:
- `name` (string, optional): Filter communities by name.
- `offset` (int32, optional): Number of items to skip (for pagination).
- `limit` (int32, optional): Maximum number of communities to return.

---

#### `POST /communities`

Creates a new community with the given name.

**Request Body** (JSON):
- `name` (string): Name of the new community.

---

#### `GET /communities/{id}`

Retrieves details of a specific community by ID.

**Path Parameters**:
- `id` (string, required): ID of the community.

---

#### `DELETE /communities/{id}`

Deletes a community by ID.

**Path Parameters**:
- `id` (string, required): ID of the community.

---

#### `PATCH /communities/{id}`

Updates a community's name or thread count offset.

**Path Parameters**:
- `id` (string, required): ID of the community.

**Request Body** (JSON):
- `name` (string, optional): New name of the community.
- `numThreadsOffset` (int32, optional): Change in number of threads.

---

#### `GET /threads`

Retrieves a list of threads. Supports filtering and pagination.

**Query Parameters**:
- `communityId` (string, optional): Filter threads by community ID.
- `title` (string, optional): Filter threads by title.
- `offset` (int32, optional): Number of items to skip.
- `limit` (int32, optional): Maximum number of threads to return.
- `sortBy` (string, optional): Sorting criteria.

---

#### `POST /threads`

Creates a new thread.

**Request Body** (JSON):
- `communityId` (string): ID of the community the thread belongs to.
- `title` (string): Title of the thread.
- `content` (string): Content of the thread.

---

#### `GET /threads/{id}`

Retrieves details of a specific thread by ID.

**Path Parameters**:
- `id` (string, required): ID of the thread.

---

#### `DELETE /threads/{id}`

Deletes a thread by ID.

**Path Parameters**:
- `id` (string, required): ID of the thread.

---

#### `PATCH /threads/{id}`

Updates fields of a thread.

**Path Parameters**:
- `id` (string, required): ID of the thread.

**Request Body** (JSON):
- `title` (string, optional): New title.
- `content` (string, optional): New content.
- `voteOffset` (int32, optional): Change in up/down votes.
- `numCommentsOffset` (int32, optional): Change in number of comments.

---

#### `GET /comments`

Retrieve a list of comments filtered by optional query parameters.

**Query Parameters:**

- `threadId` (string, optional): Filter comments by thread ID.
- `offset` (integer, optional): Pagination offset.
- `limit` (integer, optional): Pagination limit.
- `sortBy` (string, optional): Sort order or field.

---

#### `POST /comments`

Create a new comment.

**Request Body** (JSON):

- `content` (string): The text content of the comment.
- `parentId` (string, optional): The ID of the parent comment or thread.
- `parentType` (enum: `THREAD`, `COMMENT`): Type of the parent entity.

---

#### `GET /comments/{id}`

Retrieve a specific comment by its ID.

**Path Parameters:**

- `id` (string, required): The comment ID.

---

#### `DELETE /comments/{id}`

Delete a specific comment by its ID.

**Path Parameters:**

- `id` (string, required): The comment ID.

---

#### `PATCH /comments/{id}`

Update fields of a specific comment.

**Path Parameters:**

- `id` (string, required): The comment ID.

**Request Body:**

- `content` (string, optional): New content for the comment.
- `voteOffset` (integer, optional): Change in up/down votes.
- `numCommentsOffset` (integer, optional): Change in number of comments.

---

#### `POST /votes/comment/{commentId}/down`

Downvote a comment by its ID.

**Path Parameters:**

- `commentId` (string, required): The ID of the comment to downvote.

**Request Body** (JSON):

- Empty object (no fields required).

---

#### `POST /votes/comment/{commentId}/up`

Upvote a comment by its ID.

**Path Parameters:**

- `commentId` (string, required): The ID of the comment to upvote.

**Request Body** (JSON):

- Empty object (no fields required).

---

#### `POST /votes/thread/{threadId}/down`

Downvote a thread by its ID.

**Path Parameters:**

- `threadId` (string, required): The ID of the thread to downvote.

**Request Body** (JSON):

- Empty object (no fields required).

---

#### `POST /votes/thread/{threadId}/up`

Upvote a thread by its ID.

**Path Parameters:**

- `threadId` (string, required): The ID of the thread to upvote.

**Request Body** (JSON):

- Empty object (no fields required).

---

#### `GET /search`

Search across threads and communities globally.

**Query Parameters:**

- `query` (string, optional): Search keyword or phrase.
- `offset` (integer, optional): Pagination offset.
- `limit` (integer, optional): Maximum number of results to return.

---

#### `GET /search/community`

Search for communities matching the query.

**Query Parameters:**

- `query` (string, optional): Search keyword or phrase for communities.
- `offset` (integer, optional): Pagination offset.
- `limit` (integer, optional): Maximum number of results to return.

---

#### `GET /search/thread`

Search for threads matching the query.

**Query Parameters:**

- `query` (string, optional): Search keyword or phrase for threads.
- `offset` (integer, optional): Pagination offset.
- `limit` (integer, optional): Maximum number of results to return.

---

#### `GET /popular/comments`

Retrieve a list of popular comments.

**Query Parameters:**

- `offset` (integer, optional): Pagination offset for the results.
- `limit` (integer, optional): Maximum number of comments to return.

---

#### `GET /popular/threads`

Retrieve a list of popular threads.

**Query Parameters:**

- `offset` (integer, optional): Pagination offset for the results.
- `limit` (integer, optional): Maximum number of threads to return.