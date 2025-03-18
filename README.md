## Overview  

This project is a **modular monolith** REST API for a social media platform where users can create posts and comment on them. The architecture is designed to be modular, allowing for easy scaling and future migration to microservices if needed.

### Key Features:  
- **Modular Design** – Each module encapsulates its own logic, making the system maintainable and scalable.  
- **RESTful API** – Provides endpoints for managing posts, comments, and user interactions.  
- **WebSocket Support** – Enables real-time updates and notifications.  
- **Flexible Communication** – Modules interact via interfaces (synchronous) and an event bus using Go channels (asynchronous).  
- **Scalability** – Designed to handle growing business requirements while maintaining performance.  


