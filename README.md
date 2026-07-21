# API de Productos con Go y GraphQL

API GraphQL para la gestión de productos, desarrollada como parte de una prueba técnica.

## Estado

Proyecto en desarrollo.

## Tecnologías

- Go
- GraphQL
- gqlgen
- Clean Architecture

## Arquitectura

El proyecto separa las responsabilidades en las siguientes capas:

- Dominio: entidades, errores y contratos del repositorio.
- Casos de uso: reglas de negocio de productos.
- Repositorio: almacenamiento de productos en memoria.
- Delivery: servidor y resolvers GraphQL.

## Funcionalidades previstas

- Listar productos.
- Consultar un producto por ID.
- Crear productos.
- Actualizar nombre o precio.
- Eliminar productos.