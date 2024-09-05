## Install your Dependencies

- Docker Desktop: https://docs.docker.com/desktop/
- Postgresql (for `psql`): `brew install postgresql`

## Setup

### Database

- Start up your database - `docker compose up`
- Log into your database - `psql postgresql://postgres:postgres@localhost:10002`

### ORM

- Install project dependencies - `npm i`
- Update the [schema.prisma](prisma/schema.prisma) file
  - The documentation for prisma is [here](https://www.prisma.io/docs/orm/prisma-schema/data-model/models) 
- Run migrations - `npx prisma migrate dev`

### Docs

- Prisma - https://www.prisma.io/docs/getting-started/quickstart

