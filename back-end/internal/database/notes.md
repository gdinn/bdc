# Notes for cognito setup
## Manager and board member user validation
### On BDC Database
- Users table has specific properties for that. BDC services wont write on this properties carelessly.
### On Cognito (token and front-end validation)
- Custom attribute named `role` (custom:role) **MUST** be created in the user pool
- BDC front-end Cognito *App client* **MUST** have read-only access to the `role` property. This property will be written machine-to-machine with BDC back-end *App client* when `manager` or `board_member` roles are given to a certain User.
- BDC front-end auth.config.ts **MUST** have the `config.scope` property unspecified