# Sistema de Gestión de Pedidos - DB1

Proyecto académico de bases de datos que implementa un sistema completo de gestión de pedidos para un negocio de deliveries. Desarrollado en **Go** con **PostgreSQL** como base de datos relacional y **BoltDB** como base de datos NoSQL para migración de datos.

## Autores

- Diaz
- Molina
- Pachao


## Tecnologías utilizadas

- **Go 1.19** - Lenguaje de programación principal
- **PostgreSQL** - Base de datos relacional (PL/pgSQL)
- **BoltDB** - Base de datos NoSQL (key-value store)
- **lib/pq** - Driver de PostgreSQL para Go
- **bbolt** - Driver de BoltDB para Go

## Funcionalidades

### Menú interactivo principal
- Crear base de datos
- Crear tablas
- Añadir PKs y FKs
- Eliminar PKs y FKs
- Cargar datos de prueba
- Crear stored procedures y triggers
- Iniciar pruebas automatizadas
- Cargar datos en BoltDB

### Stored Procedures (PL/pgSQL)
- `crear_pedido(id_usuarie, id_direccion)` - Crea un nuevo pedido validando cliente, dirección y tarifa de envío
- `agregar_producto(id_pedido, id_producto, cantidad)` - Agrega productos al pedido con control de stock
- `entregar_pedido(id_pedido, fecha_hora_entrega)` - Marca un pedido como entregado
- `anular_pedido(id_pedido)` - Anula un pedido y devuelve el stock

### Triggers
- `trg_pedido_iniciado` - Envía email al crear un pedido
- `trg_pedido_entregado` - Envía email al entregar un pedido
- `trg_pedido_anulado` - Envía email al anular un pedido

### Migración PostgreSQL → BoltDB
- Migración de clientes, productos, direcciones y pedidos
- Solo se migran pedidos con mínimo 3 productos en detalle

## Modelado de datos

### Tablas principales
- `cliente` - Datos de los clientes
- `direccion_entrega` - Direcciones de entrega por cliente
- `tarifa_entrega` - Costos de envío por código postal
- `producto` - Catálogo de productos con stock
- `pedido` - Pedidos con estados (ingresado/entregado/anulado)
- `pedido_detalle` - Productos incluidos en cada pedido
- `error` - Registro de errores del sistema
- `envio_email` - Registro de emails enviados
- `datos_de_prueba` - Datos de prueba para ejecutar el sistema

## Estructura del proyecto

```
.
├── main.go                    # Programa principal en Go
├── go.mod                     # Módulo Go
├── go.sum                     # Dependencias
├── sql/
│   ├── tablas.sql             # Creación de tablas
│   ├── pks.sql                # Definición de claves primarias
│   ├── fks.sql                # Definición de claves foráneas
│   ├── crear_pedido.sql       # SP para crear pedido
│   ├── agregar_producto.sql   # SP para agregar producto
│   ├── entregar_pedido.sql    # SP para entregar pedido
│   ├── anular_pedido.sql      # SP para anular pedido
│   ├── eliminar_pks.sql       # Eliminación de PKs
│   ├── eliminar_fks.sql       # Eliminación de FKs
│   ├── enviar_email_pedido_iniciado.sql   # Trigger de email
│   ├── enviar_email_pedido_entregado.sql  # Trigger de email
│   └── enviar_email_pedido_anulado.sql    # Trigger de email
└── json/
    ├── clientes.json          # Datos de clientes
    ├── direcciones.json       # Direcciones de entrega
    ├── productos.json         # Catálogo de productos
    ├── tarifas_entrega.json   # Tarifas por código postal
    └── datos_de_prueba.json   # Secuencia de operaciones de prueba
```

## Prerrequisitos

- Go 1.19 o superior
- PostgreSQL instalado y corriendo
- usuario `postgres` configurado (sin contraseña por defecto)

## Instalación y uso

1. Clonar el repositorio
   ```bash
   git clone https://github.com/DiazMathias/diaz_molina_pachao_db1.git
   cd diaz_molina_pachao_db1
   ```

2. Compilar y ejecutar
   ```bash
   go run main.go
   ```

3. Seguir el menú interactivo en orden:
   - Opción 1: Crear la base de datos
   - Opción 2: Crear las tablas
   - Opción 3: Establecer PKs y FKs
   - Opción 5: Cargar datos de prueba
   - Opción 6: Crear stored procedures y triggers
   - Opción 7: Ejecutar pruebas automatizadas

4. Opcionalmente, opción 8: Migrar datos a BoltDB

## Configuración de PostgreSQL

Las credenciales de conexión están configuradas en `main.go`:

```go
const connStr = "user=postgres host=localhost dbname=diaz_molina_pachao_db1 sslmode=disable"
```

