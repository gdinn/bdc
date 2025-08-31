import { Component, OnInit } from '@angular/core';
import { AsyncPipe, JsonPipe } from '@angular/common';
import { AuthService } from '../../services/auth.service';
import { Observable } from 'rxjs';
import { UserClaims } from '../../models/auth.model';

@Component({
  selector: 'app-login',
  imports: [JsonPipe],
  templateUrl: './login.html',
  styleUrl: './login.scss',
})
export class Login {
  isAuthenticated: boolean = false
  userClaims?: UserClaims

  constructor(private authService: AuthService) {
    this.authService.isAuthenticated().subscribe(res => {
      this.isAuthenticated = res
    })

    this.authService.getUserClaims().subscribe(res => {
      this.userClaims = res
    })

    this.authService.getToken().subscribe(res => {
      console.log("Token", res)
    })
  }


  logout(): void {
    this.authService.logout()
  }

  login(): void {
    this.authService.login() 
  }
}
