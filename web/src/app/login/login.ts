import { Component, OnInit } from '@angular/core';
import {AuthenticatedResult, OidcSecurityService, LogoutAuthOptions} from 'angular-auth-oidc-client'
import { inject } from '@angular/core';
import { AsyncPipe, JsonPipe } from '@angular/common';

@Component({
  selector: 'app-login',
  imports: [AsyncPipe, JsonPipe],
  templateUrl: './login.html',
  styleUrl: './login.scss',
})
export class Login {
  private readonly oidcSecurityService = inject(OidcSecurityService);

  configuration$ = this.oidcSecurityService.getConfiguration();

  userData$ = this.oidcSecurityService.userData$;

  isAuthenticated: AuthenticatedResult | undefined = undefined;


  ngOnInit(): void {
    this.oidcSecurityService.isAuthenticated$.subscribe(
      (isAuthenticated: AuthenticatedResult) => {
        this.isAuthenticated = isAuthenticated;
        console.warn('authenticated: ', isAuthenticated);
      }
    );

    this.oidcSecurityService.checkAuth().subscribe(({ isAuthenticated, accessToken }) => {
      console.log('app authenticated', isAuthenticated);
      console.log(`Current access token is '${accessToken}'`);
    });    
  }


  logout(): void {
    // // Clear session storage
    this.oidcSecurityService
      .logoff()
      .subscribe((result) => console.log(result));
  }  

  login(): void {
    this.oidcSecurityService.authorize();   
  }
}
