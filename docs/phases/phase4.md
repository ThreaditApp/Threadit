## 🔍 Phase 4 - System Implementation

The system implementation can be found [here](../../code), divided into:
- [Services](../../code/services) - The microservices implementation that make up the backend of the system.
- [Protos](../../code/proto) - The protocol buffers used for gRPC communication between the microservices.
- [Gen](../../code/gen) - The generated code from the protocol buffers.
- [Traefik](../../code/traefik) - The configuration for the Traefik reverse proxy used for routing requests to the microservices.
- [gRPC Gateway](../../code/grpc-gateway) - The configuration for the gRPC Gateway used for exposing the gRPC services as REST APIs.

### Considerations

- **Reduced Microservices:** in this phase, we made some changes to both the application architecture and the functional requirements by reducing the number of microservices to implement. We removed the `AuthService`, the `UserService`, the `SocialService` and the `FeedService` (replaced by `PopularService`, which instead of providing threads based on the user’s followed communities and users, provides the most upvoted threads and comments). By doing this, we reduced the complexity of the system and the number of microservices to implement, which helped us achieve the goals of this phase.
- **Database:** we used a local MongoDB instance with three collections: `communities`, `threads` and `comments`.
- **Datasets:** the friends dataset was discarded as it was no longer needed and the reddit posts dataset was modified to fit our needs by fitting our three collections. Since the original dataset didn’t include comments, the `comments` collection is initially populated with an empty array, ready to store future comments.
