import { Routes } from '@angular/router';
import { Callback } from './callback/callback';
import { Login } from './login/login';
import { Logout } from './logout/logout';

export const routes: Routes = [
  {
    path: 'callback',
    component: Callback,
    title: 'Callback',
  },  
  {
    path: 'logout',
    component: Logout,
    title: 'Logout',
  },
  {
    path: '',
    component: Login,
    title: 'Login',
  },
];
