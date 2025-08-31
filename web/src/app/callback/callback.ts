import { Component } from '@angular/core';
import {AuthenticatedResult, OidcSecurityService} from 'angular-auth-oidc-client'
import { inject } from '@angular/core';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-callback',
  imports: [],
  templateUrl: './callback.html',
  styleUrl: './callback.scss'
})
export class Callback {
  constructor(private authService: AuthService){

  }
  ngOnInit() {
    this.authService.restoreSession().subscribe()
  }
}
