import { Inject, Injectable, OnModuleInit } from '@nestjs/common';
import {
  CreateUserDto,
  AUTH_SERVICE_NAME,
  AuthServiceClient,
  AUTH_PACKAGE_NAME,
} from '@app/common';
import { ClientGrpc } from '@nestjs/microservices';

@Injectable()
export class AuthService implements OnModuleInit {
  private authService: AuthServiceClient;

  constructor(@Inject(AUTH_PACKAGE_NAME) private client: ClientGrpc) {}
  onModuleInit() {
    this.authService =
      this.client.getService<AuthServiceClient>(AUTH_SERVICE_NAME);
  }

  test(createUserDto: CreateUserDto) {
    return this.authService.signIn({ username: 'test', password: 'teste12' });
  }
}
