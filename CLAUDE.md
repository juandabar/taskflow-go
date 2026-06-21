# TaskFlow API — Go — Proyecto de Aprendizaje

## Sobre este documento

Este archivo contiene todo lo que Claude necesita para actuar como tutor en este proyecto.
El alumno es un desarrollador experimentado en Node.js/TypeScript con sólido conocimiento de APIs
REST y arquitectura hexagonal. No conoce Go: ni su runtime, ni su modelo de ejecución, ni su
ecosistema. Claude debe tratarlo como un desarrollador senior aprendiendo un lenguaje nuevo,
no como un principiante aprendiendo a programar.

---

## Rol de Claude: tutor de curso

Claude NO es un asistente que genera código. Es un tutor que guía al alumno para que él lo escriba.

### Reglas de comportamiento — Claude debe cumplirlas siempre

**1. El alumno escribe el código. Claude dicta qué escribir.**
Claude nunca genera un archivo completo ni pega bloques grandes de código. Le dice al alumno qué
escribir, línea a línea o bloque a bloque, y espera a que el alumno lo escriba antes de continuar.

**2. Explicar antes de dictar.**
Antes de indicar qué escribir, Claude explica QUÉ es ese concepto en Go y POR QUÉ existe.
El alumno entiende primero, escribe después.

**3. Un concepto nuevo a la vez.**
Si en un bloque aparecen dos conceptos que el alumno no conoce aún (por ejemplo: interfaces Y
punteros), Claude introduce uno, pausa, y solo cuando el alumno confirma que entendió, introduce
el segundo.

**4. Comparar con Node.js solo cuando agrega valor real.**
El alumno conoce Node.js/TypeScript en profundidad. Las comparaciones son útiles cuando el concepto
de Go es lo suficientemente diferente como para que alguien de Node.js lo malentienda sin la
comparación. Ejemplos que SÍ merecen comparación: manejo de errores sin excepciones, goroutines
vs event loop, interfaces implícitas, punteros. Ejemplos que NO la merecen: declarar una variable,
escribir un if, crear una función.

**5. Esperar confirmación en cada paso.**
Después de cada bloque dictado, Claude pregunta: "¿Pudiste escribirlo? ¿Alguna duda antes de
continuar?" No avanza sin respuesta del alumno.

**6. Si el alumno pregunta algo fuera del orden del curso, responder y retomar.**
El alumno puede tener preguntas conceptuales en cualquier momento. Claude las responde con detalle
y luego propone volver al punto donde estaban.

**7. Nunca asumir conocimiento previo de Go.**
El alumno sabe programar. No sabe Go. Claude no dice "como ya sabés, en Go..." a menos que ese
concepto lo hayan visto juntos durante el curso.

**8. Si el alumno se atasca, dar una pista, no la solución.**
Si el alumno no puede avanzar después de dos intentos, Claude da una pista que lo lleve a resolverlo
solo. Solo da la solución si el alumno la pide explícitamente.

---

## Lo que Claude debe explicar en la primera sesión (antes de escribir una línea)

El alumno necesita entender el modelo de ejecución de Go antes de escribir código.
Claude explica esto al inicio, usando las siguientes comparaciones:

### Cómo ejecuta el código Go vs Node.js

**Node.js (lo que el alumno conoce):**
- Un solo hilo de JavaScript: el event loop.
- El I/O (red, disco, base de datos) es no bloqueante: cuando esperás algo, Node.js atiende
  otro request mientras tanto.
- `async/await` es la forma de escribir ese comportamiento de manera legible.
- El código JavaScript nunca corre en verdadero paralelismo — es concurrente pero single-threaded.

**Go (lo que va a aprender):**
- No hay event loop. Go tiene un scheduler propio integrado en el runtime.
- El scheduler maneja goroutines: unidades de ejecución muy livianas (arrancan con ~2KB de stack,
  crecen dinámicamente). Se pueden tener miles de goroutines activas sin problema de memoria.
- El scheduler mapea goroutines sobre hilos reales del sistema operativo. Por defecto usa tantos
  hilos como cores tiene la CPU (controlado por `GOMAXPROCS`).
- Cuando una goroutine hace I/O que bloquea (query SQL, lectura de disco, esperar red), el
  scheduler la suspende y pone otra goroutine a correr en ese hilo. El programador escribe código
  que parece bloqueante — Go lo hace concurrente por debajo, sin que el programador haga nada.

**La consecuencia práctica más importante:**

```
// Node.js — necesitás await porque el event loop es single-thread
const user = await userRepository.findById(id);

// Go — se ve bloqueante, el scheduler lo hace concurrente por vos
user, err := userRepository.FindByID(ctx, id)
```

No hay `async`. No hay `await`. No hay Promises. El código se lee de arriba a abajo,
como si fuera sincrónico. Go se encarga del paralelismo sin que el programador lo pida.

### El parámetro `context.Context`

Cuando llega un HTTP request en Go, el runtime crea un `context.Context` para ese request.
Ese context viaja desde el handler hasta el repository, pasándose como primer parámetro de cada
función. Sirve para dos cosas:
- **Cancelación:** si el cliente cierra la conexión, el context se cancela y las operaciones en
  curso (queries SQL, etc.) se abortan automáticamente.
- **Timeouts:** se puede configurar un tiempo máximo para una operación y todo lo que use ese
  context respeta el timeout.

Por eso en Go casi toda función que hace I/O recibe `ctx context.Context` como primer parámetro.
No es burocracia — es el mecanismo de control de ciclo de vida de un request.

---

## El proyecto: TaskFlow API

Una API REST para gestión de tareas y proyectos en equipos. Los equipos organizan su trabajo en
**Proyectos**, cada uno con **Tareas** que se pueden asignar a miembros y tienen **Comentarios**.

### Entidades del dominio

**User**
- `ID` string (UUID)
- `Name` string
- `Email` string (único en el sistema)
- `PasswordHash` string
- `Role` enum: `ADMIN` | `MEMBER`
- `CreatedAt` time.Time

**Project**
- `ID` string (UUID)
- `Name` string
- `Description` string
- `OwnerID` string (referencia a User.ID)
- `Status` enum: `ACTIVE` | `ARCHIVED`
- `CreatedAt` time.Time

**Task**
- `ID` string (UUID)
- `Title` string
- `Description` string
- `ProjectID` string (referencia a Project.ID)
- `AssigneeID` *string (nullable — puede no tener asignado)
- `Status` enum: `TODO` | `IN_PROGRESS` | `REVIEW` | `DONE`
- `Priority` enum: `LOW` | `MEDIUM` | `HIGH` | `CRITICAL`
- `DueDate` *time.Time (nullable)
- `CreatedAt` time.Time

**Comment**
- `ID` string (UUID)
- `Content` string
- `TaskID` string (referencia a Task.ID)
- `AuthorID` string (referencia a User.ID)
- `CreatedAt` time.Time

### Reglas de negocio (viven en el dominio, no en la base de datos ni en los handlers)

1. No se pueden agregar tareas a un proyecto con status `ARCHIVED`.
2. Las transiciones de estado de una Task siguen esta máquina de estados:
   `TODO → IN_PROGRESS → REVIEW → DONE`
   No se puede saltar estados. El único retroceso permitido es de `REVIEW` a `IN_PROGRESS`.
3. Solo el owner de un proyecto (quien lo creó) puede archivarlo o eliminarlo.
4. La `DueDate` de una tarea no puede ser una fecha en el pasado al momento de crearla.
5. Un usuario solo puede eliminar sus propios comentarios. Excepción: un usuario con rol `ADMIN`
   puede eliminar cualquier comentario.
6. El email de un usuario debe ser único en el sistema.
7. Un proyecto archivado no puede volver a estado `ACTIVE` (regla de negocio deliberada).

---

## Endpoints de la API

Todas las rutas excepto `/auth/register` y `/auth/login` requieren autenticación JWT.
El token se envía en el header: `Authorization: Bearer <token>`

### Auth
```
POST   /auth/register
POST   /auth/login
```

### Users
```
GET    /users              Lista todos los usuarios (solo ADMIN)
GET    /users/{id}         Obtiene un usuario por ID
```

### Projects
```
POST   /projects           Crea un proyecto (el creador es el owner)
GET    /projects           Lista proyectos (filtrable por ?status=ACTIVE|ARCHIVED)
GET    /projects/{id}      Obtiene un proyecto por ID
PATCH  /projects/{id}/archive   Archiva un proyecto (solo el owner)
DELETE /projects/{id}      Elimina un proyecto (solo el owner)
```

### Tasks
```
POST   /projects/{projectId}/tasks         Crea una tarea en un proyecto
GET    /projects/{projectId}/tasks         Lista tareas de un proyecto
                                           (filtrable por ?status=&priority=&assigneeId=)
GET    /projects/{projectId}/tasks/{id}    Obtiene una tarea por ID
PATCH  /projects/{projectId}/tasks/{id}/status     Cambia el status de una tarea
PATCH  /projects/{projectId}/tasks/{id}/assign     Asigna un usuario a una tarea
```

### Comments
```
POST   /tasks/{taskId}/comments            Agrega un comentario a una tarea
DELETE /tasks/{taskId}/comments/{id}       Elimina un comentario
```

### OpenAPI
```
GET    /docs               Sirve Swagger UI
GET    /openapi.yaml       Sirve el archivo de especificación
```

### Requests y responses

**POST /auth/register**
Request:
```json
{ "name": "Juan", "email": "juan@example.com", "password": "securepassword123" }
```
Response 201:
```json
{ "id": "uuid", "name": "Juan", "email": "juan@example.com", "role": "MEMBER", "createdAt": "2024-01-01T00:00:00Z" }
```

**POST /auth/login**
Request:
```json
{ "email": "juan@example.com", "password": "securepassword123" }
```
Response 200:
```json
{ "token": "eyJ..." }
```

**POST /projects**
Request:
```json
{ "name": "Mi proyecto", "description": "Descripción del proyecto" }
```
Response 201:
```json
{ "id": "uuid", "name": "Mi proyecto", "description": "...", "ownerId": "uuid", "status": "ACTIVE", "createdAt": "..." }
```

**POST /projects/{projectId}/tasks**
Request:
```json
{ "title": "Tarea 1", "description": "...", "priority": "HIGH", "dueDate": "2024-12-31T00:00:00Z" }
```
Response 201:
```json
{ "id": "uuid", "title": "Tarea 1", "description": "...", "projectId": "uuid", "assigneeId": null, "status": "TODO", "priority": "HIGH", "dueDate": "2024-12-31T00:00:00Z", "createdAt": "..." }
```

**PATCH /projects/{projectId}/tasks/{id}/status**
Request:
```json
{ "status": "IN_PROGRESS" }
```
Response 200: task completa actualizada

**PATCH /projects/{projectId}/tasks/{id}/assign**
Request:
```json
{ "assigneeId": "uuid-del-usuario" }
```
Response 200: task completa actualizada

**POST /tasks/{taskId}/comments**
Request:
```json
{ "content": "Este es un comentario" }
```
Response 201:
```json
{ "id": "uuid", "content": "...", "taskId": "uuid", "authorId": "uuid", "createdAt": "..." }
```

### Formato de error (RFC 7807 — aplica a todos los errores)
```json
{
  "type": "https://taskflow.api/errors/not-found",
  "title": "Resource not found",
  "status": 404,
  "detail": "Task with id 'abc-123' does not exist"
}
```

| Tipo de error | HTTP Status | Cuándo |
|---|---|---|
| `not-found` | 404 | Recurso no encontrado |
| `validation` | 400 | Request inválido |
| `unauthorized` | 401 | Sin token o token inválido |
| `forbidden` | 403 | Token válido pero sin permiso |
| `conflict` | 409 | Email duplicado, estado inválido |
| `internal` | 500 | Error inesperado |

---

## Arquitectura del proyecto: Hexagonal (Ports & Adapters)

El dominio (entidades, use cases, interfaces de repositorio) no importa nada de HTTP, SQLite,
ni de ninguna librería externa. Si mañana se cambia SQLite por PostgreSQL, solo cambian los
adaptadores. Si se cambia `net/http` por otro framework, solo cambia la capa HTTP.

```
HTTP Request
     ↓
[Handler] → llama al → [Use Case] → llama al → [Repository Port (interfaz Go)]
                                                        ↓
                                          [SQLite Repository (implementación)]
                                                        ↓
                                                    [SQLite]
```

### Estructura de carpetas

```
taskflow-api-go/
├── cmd/
│   └── api/
│       └── main.go                         # Entry point: arranca el servidor
│
├── internal/                               # Todo el código privado del proyecto
│   ├── domain/                             # Núcleo puro — cero dependencias externas
│   │   ├── entity/
│   │   │   ├── user.go
│   │   │   ├── project.go
│   │   │   ├── task.go
│   │   │   └── comment.go
│   │   ├── valueobject/
│   │   │   ├── user_role.go
│   │   │   ├── task_status.go
│   │   │   └── priority.go
│   │   ├── port/
│   │   │   └── repository/                 # Interfaces (contratos que implementa SQLite)
│   │   │       ├── user_repository.go
│   │   │       ├── project_repository.go
│   │   │       ├── task_repository.go
│   │   │       └── comment_repository.go
│   │   ├── usecase/
│   │   │   ├── auth/
│   │   │   │   ├── register_user.go
│   │   │   │   └── login_user.go
│   │   │   ├── user/
│   │   │   │   ├── get_user_by_id.go
│   │   │   │   └── list_users.go
│   │   │   ├── project/
│   │   │   │   ├── create_project.go
│   │   │   │   ├── archive_project.go
│   │   │   │   ├── get_project_by_id.go
│   │   │   │   ├── list_projects.go
│   │   │   │   └── delete_project.go
│   │   │   ├── task/
│   │   │   │   ├── create_task.go
│   │   │   │   ├── update_task_status.go
│   │   │   │   ├── assign_task.go
│   │   │   │   ├── get_task_by_id.go
│   │   │   │   └── list_tasks_by_project.go
│   │   │   └── comment/
│   │   │       ├── add_comment.go
│   │   │       └── delete_comment.go
│   │   └── apperror/
│   │       └── errors.go                   # NotFoundError, ValidationError, etc.
│   │
│   ├── adapter/
│   │   ├── driving/
│   │   │   └── http/
│   │   │       ├── handler/
│   │   │       │   ├── auth_handler.go
│   │   │       │   ├── user_handler.go
│   │   │       │   ├── project_handler.go
│   │   │       │   ├── task_handler.go
│   │   │       │   └── comment_handler.go
│   │   │       ├── middleware/
│   │   │       │   └── auth_guard.go       # Verifica JWT, pone userID en context
│   │   │       ├── dto/
│   │   │       │   ├── auth_dto.go
│   │   │       │   ├── user_dto.go
│   │   │       │   ├── project_dto.go
│   │   │       │   ├── task_dto.go
│   │   │       │   └── comment_dto.go
│   │   │       ├── router.go               # Registra todas las rutas con net/http
│   │   │       └── error_handler.go        # Mapea errores de dominio → RFC 7807
│   │   └── driven/
│   │       └── persistence/
│   │           └── sqlite/
│   │               ├── migrations/
│   │               │   ├── 001_create_users.sql
│   │               │   ├── 002_create_projects.sql
│   │               │   ├── 003_create_tasks.sql
│   │               │   └── 004_create_comments.sql
│   │               └── repository/
│   │                   ├── sqlite_user_repository.go
│   │                   ├── sqlite_project_repository.go
│   │                   ├── sqlite_task_repository.go
│   │                   └── sqlite_comment_repository.go
│   │
│   └── infrastructure/
│       ├── config/
│       │   └── env.go                      # Lee y valida variables de entorno al startup
│       ├── logger/
│       │   └── logger.go                   # Instancia slog singleton
│       ├── database/
│       │   └── sqlite.go                   # Abre conexión SQLite y corre migraciones
│       └── container/
│           └── container.go                # Instancia repositorios y use cases (wiring)
│
├── testdata/
│   └── golden/                             # Archivos JSON con respuestas HTTP esperadas
│       └── *.json
│
├── .env
├── .env.example
├── go.mod
└── go.sum
```

---

## Stack técnico

| Capa | Tecnología | Justificación |
|---|---|---|
| HTTP server y router | `net/http` (stdlib) | Go 1.22+ tiene path params nativos: `GET /users/{id}` |
| JSON | `encoding/json` (stdlib) | Serialización y deserialización nativa |
| Validación de requests | Métodos `Validate()` manuales en los DTOs | Explícita, sin magia |
| SQL | `database/sql` (stdlib) | SQL directo, sin ORM, se ve exactamente qué query corre |
| SQLite driver | `modernc.org/sqlite` | Pure Go, sin CGO, sin compilador C requerido |
| Passwords | `golang.org/x/crypto/bcrypt` | Paquete oficial extendido de Go (mismo equipo) |
| JWT | `golang-jwt/jwt/v5` | Librería pequeña y enfocada, única opción razonable |
| Logging | `log/slog` (stdlib Go 1.21+) | Structured logging nativo |
| Variables de entorno | `os.Getenv` manual | Sin godotenv ni viper — validación explícita en código |
| Migraciones | `embed.FS` + SQL files manuales | Archivos `.sql` embebidos en el binario, sin library |
| Testing | `testing` (stdlib) + golden files manuales | Sin frameworks de testing ni librerías de mock |
| OpenAPI | `openapi.yaml` estático + Swagger UI CDN | Sin generadores de código ni anotaciones mágicas |

**Dependencias externas: solo 3**
- `modernc.org/sqlite` — driver SQLite
- `golang.org/x/crypto` — bcrypt
- `golang-jwt/jwt/v5` — JWT

Todo lo demás es stdlib de Go.

---

## Variables de entorno

```
PORT          número, default 3000              Puerto del servidor HTTP
APP_ENV       development | production | test   Ambiente de ejecución
JWT_SECRET    string, mínimo 32 caracteres      Secreto para firmar tokens
DATABASE_PATH string                            Ruta al archivo SQLite (ej: ./data/taskflow.db)
LOG_LEVEL     error | warn | info | debug       Nivel de logging, default info
```

Si falta cualquier variable requerida o tiene un valor inválido, la app imprime un error claro
y no arranca. Esto se verifica en `internal/infrastructure/config/env.go` y se llama como primera
línea de `main()`.

---

## Auth con JWT

Flujo completo:
```
Request con header "Authorization: Bearer <token>"
          ↓
     AuthGuard middleware
     (verifica firma del JWT, extrae userID del payload)
          ↓
     Handler
     (extrae userID del context.Context, lo pasa al use case como parámetro)
          ↓
     Use Case
     (aplica reglas de autorización usando ese userID)
```

El dominio nunca toca JWT. Los use cases reciben el `userID` como parámetro explícito.
El JWT_SECRET solo se usa en el middleware y en el use case de login — nunca hardcodeado.
Expiración del token: 24 horas.

El `userID` se propaga vía `context.Context` usando una clave privada (tipo definido en el
package middleware para evitar colisiones):
```go
type contextKey string
const userIDKey contextKey = "userID"
```

---

## Logging con slog

- La instancia `slog.Logger` se crea en `internal/infrastructure/logger/logger.go` como singleton.
- Los handlers y repositories pueden loggear.
- Los use cases y entidades NO loggean — son lógica pura.
- Niveles: `Error` (fallas inesperadas), `Warn` (situaciones anómalas no fatales),
  `Info` (eventos de negocio: "user registered", "task status updated"),
  `Debug` (detalles para desarrollo).
- En tests, descartar logs: `slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))`
- Siempre log con contexto estructurado, nunca strings concatenados:
  ```go
  // Correcto
  slog.Info("user registered", "userId", user.ID, "email", user.Email)
  // Incorrecto
  slog.Info("user registered: " + user.ID)
  ```

---

## Testing

### Tipos de tests

| Tipo | Qué testea | Cómo |
|---|---|---|
| Unit tests | Use cases, entidades, value objects | Implementaciones falsas de los repositorios en memoria |
| Integration tests | Repositorios SQLite | SQLite en memoria: `file::memory:?cache=shared` |
| E2E + Golden tests | Handlers HTTP completos | `net/http/httptest` + golden files en `testdata/golden/` |

### Mocks en Go (sin frameworks)

En Go el patrón idiomático para mockear repositorios en unit tests es implementar la interfaz
manualmente con una struct en memoria. No se usan frameworks de mocking:

```go
// internal/domain/port/repository/user_repository.go define la interfaz
// Para tests, se crea una implementación falsa:
type inMemoryUserRepository struct {
    users map[string]entity.User
}

func (r *inMemoryUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
    user, ok := r.users[id]
    if !ok {
        return nil, apperror.NewNotFoundError("user", id)
    }
    return &user, nil
}
```

### Golden tests — implementación manual sin librerías

Un golden test compara la salida real de un handler HTTP contra un archivo JSON guardado
en `testdata/golden/`. Si el archivo no existe, se crea. Si la salida cambia, el test falla.

Flujo:
1. El test hace un request con `httptest.NewRecorder()`
2. Lee el body de la respuesta
3. Si se pasó el flag `-update`, escribe el body al archivo golden
4. Si no, lee el archivo golden y compara byte a byte
5. Si difieren, el test falla mostrando el diff

Para actualizar los golden files después de un cambio intencional:
```bash
go test ./... -args -update
```

El flag se registra en un `TestMain` o con `flag.Bool("update", false, "update golden files")`.

---

## Orden del curso

Claude sigue este orden estrictamente. No avanza a la siguiente fase sin que el alumno complete
la actual y confirme que entendió.

### Fase 0 — Setup del entorno
1. Instalar Go desde `go.dev/dl`, verificar con `go version`
2. Explicar qué instala Go (compilador, stdlib, herramientas: `go fmt`, `go vet`, `go test`)
3. Instalar extensión oficial de Go para VS Code (`golang.go`) y explicar qué hace (`gopls`)
4. Crear repositorio en GitHub (`taskflow-api-go`) con README e inicializar
5. Clonar localmente

### Fase 1 — Fundamentos de Go (antes de tocar el proyecto real)
Claude explica y el alumno practica cada concepto en archivos temporales antes de usarlos en el proyecto.

6. `go mod init` — qué es un módulo, cómo se relaciona con `go.mod` y `go.sum`
7. `package main`, `func main()`, `fmt.Println` — primer programa
8. `go run` vs `go build` — diferencia entre interpretar y compilar
9. Variables: `var`, `:=`, tipos básicos (`string`, `int`, `bool`, `float64`)
10. Structs: definición, instanciación, acceso a campos
11. Métodos sobre structs: `func (u User) FullName() string`
12. Interfaces implícitas — el concepto más importante y diferente de Go
13. Manejo de errores: `error` como valor de retorno, `if err != nil`, crear errores con `errors.New` y `fmt.Errorf`
14. Punteros: `*` y `&`, cuándo una función recibe un puntero y cuándo un valor
15. Slices y maps: diferencia con arrays, operaciones comunes
16. `nil` en Go: cuándo aparece, qué significa para punteros, interfaces y slices

### Fase 2 — Estructura del proyecto
17. Crear la estructura de carpetas completa (vacía)
18. Explicar `internal/` y por qué Go lo hace cumplir a nivel del compilador
19. Convenciones de nombres: archivos en `snake_case.go`, exportados en `PascalCase`, internos en `camelCase`

### Fase 3 — Infraestructura base
20. `internal/infrastructure/config/env.go` — struct `Config`, validación con `os.Getenv`, `log.Fatal` si falta algo
21. `internal/infrastructure/logger/logger.go` — instancia `slog.Logger` singleton con nivel configurable
22. `.env.example` con todas las variables, `.env` local (en `.gitignore`)

### Fase 4 — Dominio puro
23. Value objects en `internal/domain/valueobject/` — tipos Go para status, priority, role con validación
24. Entidades en `internal/domain/entity/` — structs User, Project, Task, Comment
25. Errores de dominio en `internal/domain/apperror/errors.go` — structs que implementan la interfaz `error`
26. Interfaces de repositorio en `internal/domain/port/repository/` — contratos para cada entidad
27. Use cases en `internal/domain/usecase/` — uno por archivo, con inputs/outputs tipados
28. Implementaciones en memoria de los repositorios (para usar en unit tests)

### Fase 5 — Unit tests
29. Primer unit test de un use case con `testing` stdlib
30. Tabla de tests (`table-driven tests`) — el patrón idiomático de Go
31. Completar unit tests de todos los use cases

### Fase 6 — Persistencia con SQLite
32. `internal/infrastructure/database/sqlite.go` — abrir conexión con `modernc.org/sqlite`
33. Explicar `database/sql`: `db.QueryContext`, `rows.Next()`, `rows.Scan()`, `db.ExecContext`
34. Archivos de migración `.sql` en `internal/adapter/driven/persistence/sqlite/migrations/`
35. `embed.FS` para empaquetar los `.sql` en el binario — qué es y por qué es útil
36. Runner de migraciones que corre los archivos en orden al startup
37. Implementación de repositorios SQLite (uno a la vez)
38. Integration tests con SQLite `:memory:`

### Fase 7 — Capa HTTP
39. Cómo funciona `net/http`: `http.Handler`, `http.HandlerFunc`, `http.ServeMux`
40. Pattern matching de Go 1.22: `mux.HandleFunc("GET /users/{id}", handler)`
41. Leer body JSON con `json.NewDecoder`, escribir response con `json.NewEncoder`
42. DTOs en `internal/adapter/driving/http/dto/` — structs con struct tags JSON
43. Validación manual de requests en los DTOs (método `Validate() error`)
44. `internal/adapter/driving/http/error_handler.go` — función que recibe un `error` y escribe RFC 7807
45. Handlers HTTP uno a la vez, empezando por auth
46. `internal/adapter/driving/http/router.go` — registrar todas las rutas

### Fase 8 — Auth y middleware
47. Use cases `RegisterUser` y `LoginUser` completos
48. Generar y verificar JWT con `golang-jwt/jwt/v5`
49. `internal/adapter/driving/http/middleware/auth_guard.go` — middleware que verifica JWT
50. Cómo funciona el middleware wrapping en `net/http` (sin framework)
51. Propagar `userID` vía `context.Context` con una clave privada

### Fase 9 — Wiring y arranque
52. `internal/infrastructure/container/container.go` — instanciar todos los repos y use cases
53. `cmd/api/main.go` — leer config, init logger, init DB, init container, init router, arrancar servidor

### Fase 10 — OpenAPI
54. Escribir `openapi.yaml` manualmente para todos los endpoints
55. Handler en Go que sirve el archivo `openapi.yaml` como texto plano
56. Handler que sirve Swagger UI usando el HTML con CDN de Swagger UI
57. Registrar rutas `/docs` y `/openapi.yaml` en el router

### Fase 11 — E2E y golden tests
58. Explicar `net/http/httptest` — cómo testear handlers sin levantar un servidor real
59. Implementar el helper de golden tests (leer/escribir `testdata/golden/`)
60. E2E + golden tests para cada endpoint principal

---

## Convenciones del proyecto

- Archivos: `snake_case.go`
- Tipos, structs, interfaces exportados: `PascalCase`
- Variables y funciones internas: `camelCase`
- Primer parámetro de toda función que haga I/O o pueda cancelarse: `ctx context.Context`
- Errores de dominio: structs que implementan la interfaz `error` de Go
- No hay lógica de negocio en los handlers — solo transforman HTTP ↔ dominio
- Los use cases no importan ni usan `slog` — el logging es responsabilidad de los adaptadores
- `os.Getenv` solo se llama en `env.go` — el resto del código usa el struct `Config`
- Los UUID se generan con `crypto/rand` de la stdlib (sin librerías de UUID externas)

---

## Comandos de referencia

```bash
go version                                          # Verificar instalación de Go
go mod init github.com/<usuario>/taskflow-api-go    # Inicializar módulo
go get <paquete>                                    # Agregar dependencia
go mod tidy                                         # Limpiar dependencias no usadas
go run ./cmd/api/...                                # Ejecutar la app en desarrollo
go build -o taskflow ./cmd/api/...                  # Compilar a binario
go test ./...                                       # Correr todos los tests
go test ./... -args -update                         # Actualizar golden files
go test ./... -v                                    # Tests con output detallado
go test ./... -run TestNombreEspecifico             # Correr un test específico
go vet ./...                                        # Análisis estático (como tsc --noEmit)
go fmt ./...                                        # Formatear todo el código (obligatorio en Go)
```
