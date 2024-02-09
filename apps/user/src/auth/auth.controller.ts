import { Controller } from '@nestjs/common';
import { AuthService } from './auth.service';
import {
  AuthServiceControllerMethods,
  AuthServiceController,
  ProviderAuthRequest,
  Tokens,
  RefreshTokensRequest,
  Empty,
  LogOutRequest,
  DeleteAccountRequest,
  DeleteAccountResponse,
  RegisterDto,
  UserAuth,
  LoginDto,
} from '@app/common';
import { Observable } from 'rxjs';

@Controller()
@AuthServiceControllerMethods()
export class AuthController implements AuthServiceController {
  constructor(private readonly authService: AuthService) {}

  providerAuth(
    request: ProviderAuthRequest,
  ): Promise<Tokens> | Observable<Tokens> | Tokens {
    return this.authService.providerAuth(
      request.email,
      request.agent,
      request.provider,
    );
  }
  refreshTokens(
    request: RefreshTokensRequest,
  ): Promise<Tokens> | Observable<Tokens> | Tokens {
    return this.authService.refreshTokens(request.refreshToken, request.agent);
  }

  logOut(request: LogOutRequest): Promise<Empty> | Observable<Empty> | Empty {
    return this.authService.logOut(request.id, request.user);
  }

  deleteAccount(
    request: DeleteAccountRequest,
  ):
    | Promise<DeleteAccountResponse>
    | Observable<DeleteAccountResponse>
    | DeleteAccountResponse {
    return this.authService.deleteAccount(request.id, request.user);
  }

  signIn(
    request: RegisterDto,
  ): Promise<UserAuth> | Observable<UserAuth> | UserAuth {
    return this.authService.signIn(request);
  }

  logIn(request: LoginDto): Promise<Tokens> | Observable<Tokens> | Tokens {
    return this.authService.logIn(request, request.agent);
  }
}
