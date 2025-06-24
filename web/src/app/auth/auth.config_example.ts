import { PassedInitialConfig } from 'angular-auth-oidc-client';
import { provideAuth, LogLevel } from 'angular-auth-oidc-client';

export const authConfig: PassedInitialConfig = {
  config: {
        authority: '',
        redirectUrl: 'http://localhost:4200/callback',
        postLogoutRedirectUri: 'http://localhost:4200',
        clientId: '',
        scope: 'email openid phone',
        responseType: 'code',
        silentRenew: true,
        useRefreshToken: true,
        logLevel: LogLevel.Debug,        
  }
}