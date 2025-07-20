# Notes for cognito setup
## Manager and board member user validation
### On BDC Database
- Users table has specific properties for that. BDC services wont write on this properties carelessly.
### On Cognito (token and front-end validation)
- BDC front-end Cognito *App client* **MUST** have read-only access to the `profile` property. This property will be written machine-to-machine with another *App client* with write access by some root-user BDC service or via AWS dashboard directly.
- BDC front-end **MUST** have `profile` property within the scope on `auth.config.ts` file.