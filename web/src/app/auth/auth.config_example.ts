import { PassedInitialConfig } from 'angular-auth-oidc-client';
import { provideAuth, LogLevel } from 'angular-auth-oidc-client';

export const authConfig: PassedInitialConfig = {
  config: {
        authority: 'PROVIDED_ON_EXAMPLE',
        redirectUrl: 'http://localhost:4200/callback',
        clientId: 'PROVIDED_ON_EXAMPLE', 
        postLogoutRedirectUri: "http://localhost:4200/logout",
        scope: 'email openid phone',
        responseType: 'code',
        silentRenew: true,
        useRefreshToken: true,
        logLevel: LogLevel.Debug,        
  }
}

export const cognitoLogoutUrl = "PROVIDED_ON_EXAMPLE";