import { Inject, Injectable, OnModuleInit } from '@nestjs/common';
import {
  CreateUserDto,
  AUTH_SERVICE_NAME,
  AuthServiceClient,
  AUTH_PACKAGE_NAME,
  RegisterDto,
  LoginDto,
  LogOutRequest,
  RefreshTokensRequest,
  ProviderAuthRequest,
  DeleteAccountRequest
} from "@app/common";
import { ClientGrpc } from '@nestjs/microservices';

@Injectable()
export class AuthService implements OnModuleInit {
  private authService: AuthServiceClient;

  constructor(@Inject(AUTH_PACKAGE_NAME) private client: ClientGrpc) {}
  onModuleInit() {
    this.authService =
      this.client.getService<AuthServiceClient>(AUTH_SERVICE_NAME);
  }

  SignUp(registerDto: RegisterDto) {
    return this.authService.signIn(registerDto);
  }

  LogIn(loginDto: LoginDto) {
    return this.authService.logIn(loginDto);
  }

  LogOut(logOutRequest: LogOutRequest) {
    return this.authService.logOut(logOutRequest);
  }

  RefreshTokens(refreshTokensRequest: RefreshTokensRequest) {
    return this.authService.refreshTokens(refreshTokensRequest);
  }

  ProviderAuth(providerAuthRequest: ProviderAuthRequest) {
    return this.authService.providerAuth(providerAuthRequest);
  }

  DeleteAccount(deleteAccountRequest: DeleteAccountRequest) {
    return this.authService.deleteAccount(deleteAccountRequest);
  }
}
