import {AuthenticatedResult, OidcSecurityService, LogoutAuthOptions} from 'angular-auth-oidc-client'
import { inject } from '@angular/core';
import { concat, concatMap, iif, map, Observable, of } from 'rxjs';
import { UserClaims } from '../models/auth.model';
import { cognitoLogoutUrl } from '../app/auth/auth.config';

export class AuthService {
  private readonly oidcSecurityService = inject(OidcSecurityService);

  getUserRole(): Observable<string> {
    return this.oidcSecurityService.getPayloadFromAccessToken(false).pipe(
      map(res => res.role)
    )
  }

  getUserClaims(): Observable<UserClaims | undefined> {
    const userClaims$ = this.oidcSecurityService.getPayloadFromAccessToken(false).pipe(
      map(res => {
        return {
          role: res.role,
          username: res.username,
          email: res.email
        } as UserClaims
      })
    )

    return this.isAuthenticated().pipe(
      concatMap( res =>
        iif(
          () => res == true,
          userClaims$,
          of(undefined)
        )
      )
    )
  }

  getToken(): Observable<string | undefined> {
    return this.isAuthenticated().pipe(
      concatMap( res =>
        iif(
          () => res == true,
          this.oidcSecurityService.getAccessToken(),
          of(undefined)
        )
      )
    )    
  }

  restoreSession(): Observable<UserClaims | undefined> {
    return this.oidcSecurityService.checkAuth().pipe(
      concatMap(res => 
        iif(
          () => res != undefined,
          this.getUserClaims(),
          of(undefined)
        )
      )
    )
  }

  isAuthenticated(): Observable<boolean> {
    return this.oidcSecurityService.isAuthenticated$.pipe(
      map(res => res.isAuthenticated)
    )
  }

  login() {
    this.oidcSecurityService.authorize();      
  }

  logout() {
      // Clear session storage
    if (window.sessionStorage) {
      window.sessionStorage.clear();
    }

    window.location.href = cognitoLogoutUrl
  }
}
