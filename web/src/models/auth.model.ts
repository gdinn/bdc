import { UserRole } from "./user.model"

export interface UserClaims {
  username: string
  email: string
  role: UserRole
}