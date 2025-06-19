# SoccerManagerApi

## Get started

Before launching the program, read the [Configuration](#config) section

```bash
docker compose up --build
```

Or refer to [Makefile](#makefile) section

### Config

To run the project, you'll need to deal with environment variables first

The repository includes the examples of how `.env` files would look like

there are 3 environment files needed for running:

- `.env`- used by docker
- `.soccer-manager.env` - used by the app
- `.envrc` - used by the makefile

All values should be exactly same everywhere

When you decide to use them, just remove `example` suffix

### Makefile

I included very reach `Makefile`, it includes simple scripts as well as very specific and complex logics too.

It works as well as on windows and linux. (batch and bash)

- Build project
  - `make build`
- Build (all platforms and architectures)
  - `make build/all`
- Openapi Spec
  - `make swag-gen`

## API

### Docs

app comes with openapi and postman api specifications.

everything is stored in [/api](api) folder.

### Authentication

Access Token: Used as a Bearer `token` authorization header. shortlived, so you’ll often refresh it.
Refresh Token: Sent as an HTTP cookie. Use it to request a new access token via POST `/v1/auth/refresh`.

### Translations

For any endpoint returning translated data (e.g. GET `/v1/teams`, `/v1/teams/{team_id}`, `/v1/users/{user_id}/team`, `/v1/player-positions`), add the `Accept-Language` header to specify your locale ([ISO 639-1](https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes) code): en, es, fr, ka, etc.

### Language and Country codes

As we sait in [Translations](#translations), you'll need to pass an `Accept-Language` with locale code.

But in places where `country code` is neede, you'll need to pass [ISO 3166](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes) country codes in the request body (eg: US, UK, GE, etc.)

### Tags (Short Overview)

auth – Login, register, token refresh.
users – Manage user profiles and account info.
teams – Retrieve and update team data, translations.
player-positions – Manage position codes and translations.
players – Fetch, create, update player data.
transfers – Create and handle transfer listings.
transfer-records – Look at transfer history details.
Keep your access token valid or refresh it with the refresh token if necessary.

## Structure

The app is divided into 3 main parts:

- `delivery` acts as a transport and sends data to the client
- `domain` includes all business and domain requirements
- `data` communicates with databases and stores data

### Infrastructure

Postgresql as a storage and nginx as a reverse proxy

It's quite possible to include redis but at this moment app uses only in memory cache.

## License

This project is licensed under the **Apache Software License 2.0**.

See [LICENSE](LICENSE) for more information.
