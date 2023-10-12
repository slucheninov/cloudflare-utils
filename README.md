# Cloudflare utils

Only enterprise domain !

## Environment variable

This API token will affect the below accounts and zones, along with their respective permissions

- Account Analytics:Read
- All zones - Zone Settings:Read, Zone:Read, Analytics:Read
- All users - User Details:Read

```bash
export CF_API_TOKEN=
```

## Run examples

```bash
go run main.go -startDate=2021-08-01 -stopDate=2021-08-02
```
Test
