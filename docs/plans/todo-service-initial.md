# TodoService Initial Implementation Plan

To implement the TodoService as a Connect buf service in `services/todo` with all constraints, here's my complete final revised step-by-step plan:

1. **Adapt the protobuf schema**: Confirm `Task.id` is `string` for UUID7. Regenerate code.

2. **Define internal models and requests/responses**: In `services/todo/internal/models/`, create structs like `Task`, `TodoList`, and internal request/response types (e.g., `CreateTaskRequest { Title string }`, `CreateTaskResponse { Task Task }`) mirroring protobuf but using internal models.

3. **Implement in-memory KV storage**: In `services/todo/internal/database/kv.go`, define `KVStore` with methods `Set`, `Get`, `Delete`, `List`.

4. **Add repository layer**: In `services/todo/internal/database/repository/`, create `TaskRepository` interface with methods like `CreateTask(ctx context.Context, req CreateTaskRequest) (CreateTaskResponse, error)`, implemented using `KVStore`, JSON serialization, and UUID7 generation.

5. **Add service layer**: In `services/todo/internal/service/todo/`, define `TodoService` struct with `TaskRepository` field. Implement methods with the same signature format (e.g., `CreateTask(ctx context.Context, req CreateTaskRequest) (CreateTaskResponse, error)`), adding business logic like validation.

6. **Create the gRPC handler**: In `services/todo/internal/grpc/handler.go`, define a struct with `TodoService` field, implementing `todov1connect.TodoServiceHandler`. Convert protobuf requests to internal, call service, convert responses back.

7. **Set up the Connect server**: In `services/todo/cmd/connect/main.go`, instantiate layers (KVStore -> Repository -> TodoService -> Handler), set up Connect mux, start HTTP server on port 4000.

8. **Testing**: Add unit tests for each layer, verifying the signature format and flow.

This ensures consistent method signatures across layers. If you approve, I can implement it.