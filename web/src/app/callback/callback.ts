import { Component } from '@angular/core';
import {AuthenticatedResult, OidcSecurityService} from 'angular-auth-oidc-client'
import { inject } from '@angular/core';

@Component({
  selector: 'app-callback',
  imports: [],
  templateUrl: './callback.html',
  styleUrl: './callback.scss'
})
export class Callback {
  private readonly oidcSecurityService = inject(OidcSecurityService);

  ngOnInit() {
    const token = this.oidcSecurityService.getAccessToken().subscribe((token) => {
      console.log("Token: ", token)
    });
    this.oidcSecurityService.checkAuth().subscribe(res => {
      console.log("Aqui oh", res)
    })
  }
}
