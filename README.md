# Ludo Backend (Go)

A production-grade real-time Ludo backend built in Go.

## Features
- JWT-based authentication
- Wallet with ledger-based accounting
- Real-time gameplay using WebSockets
- Admin APIs for commissions, maintenance, payouts
- Rate limiting & load tested APIs
- Dockerized local setup

## Architecture
- API Gateway (separate service)
- Core Go backend (this repo)
- Legacy Node.js game engine (temporary)

## Tech Stack
- Go
- PostgreSQL
- Redis
- WebSockets
- Docker
- k6 (load testing)

## Status
🚧 Work in progress
