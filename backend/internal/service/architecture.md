# Storage Decision

Recommended production split for this app:

## PostgreSQL

Use PostgreSQL as the transactional source of truth for:

- users and auth/session metadata
- wallet balances and Sparks transactions
- referrals and rewards
- likes, matches, purchases
- chat room metadata
- catalog tables like gifts and Spark packages

Why:

- strong consistency for money-like balances and rewards
- safer relational constraints for referrals, matches, and ownership
- easier reporting and analytics over transactional data

## MongoDB

Use MongoDB for document-heavy or flexible-shape data:

- rich profile documents
- questionnaire answers
- discovery/search projection documents
- chat messages
- activity feed documents if they become high-volume

Why:

- questionnaire/profile structure will evolve often
- chat messages are append-heavy and document-friendly
- denormalized read models are easier for discovery and feed use cases

## Code Shape

The code is now arranged so:

- `internal/http` acts as controllers
- `internal/service` holds business logic
- `internal/repository` defines domain repository contracts
- `internal/repository/memory` is the current adapter

Next production step:

- add `internal/repository/postgres` for transactional domains
- add `internal/repository/mongo` for profile/message document domains
- compose both in the app container behind the existing service interfaces
