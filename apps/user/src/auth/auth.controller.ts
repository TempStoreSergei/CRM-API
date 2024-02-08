import { Controller } from '@nestjs/common';
import { AuthService } from './auth.service';
import {
  AuthServiceControllerMethods,
  AuthServiceController,
  SignOutRequest,
  SignOutResponse,
  SignInRequest,
  JwtToken,
  RestorePasswordRequest,
  RestorePasswordResponse,
  VerifyTokenRequest,
  VerifyTokenResponse,
} from '@app/common';
import { Observable } from 'rxjs';

@Controller()
@AuthServiceControllerMethods()
export class AuthController implements AuthServiceController {
  constructor(private readonly authService: AuthService) {}

  signOut(
    request: SignOutRequest,
  ): Promise<SignOutResponse> | Observable<SignOutResponse> | SignOutResponse {
    return this.authService.signOut(request);
  }

  signIn(
    request: SignInRequest,
  ): Promise<JwtToken> | Observable<JwtToken> | JwtToken {
    return this.authService.signIn(request);
  }

  restorePassword(
    request: RestorePasswordRequest,
  ):
    | Promise<RestorePasswordResponse>
    | Observable<RestorePasswordResponse>
    | RestorePasswordResponse {
    return this.authService.restorePassword(request);
  }

  verifyToken(
    request: VerifyTokenRequest,
  ):
    | Promise<VerifyTokenResponse>
    | Observable<VerifyTokenResponse>
    | VerifyTokenResponse {
    return this.authService.verifyToken(request);
  }
}
