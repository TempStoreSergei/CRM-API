import {
  Injectable,
  OnModuleInit,
  UnauthorizedException,
} from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import {
  JwtToken,
  RestorePasswordRequest,
  RestorePasswordResponse,
  SignInRequest,
  SignOutRequest,
  SignOutResponse,
  User,
  VerifyTokenResponse,
} from '@app/common';
import { randomUUID } from 'crypto';

@Injectable()
export class AuthService implements OnModuleInit {
  private readonly users: User[] = [];

  onModuleInit(): any {
    for (let i = 0; i < 1; i++) {
      this.users.push({
        age: 0,
        username: 'test-id',
        password: `password-one`,
        id: randomUUID(),
        socialMedia: {},
        subscribed: false,
      });
    }
  }

  constructor(private jwtService: JwtService) {}

  findOne(id: string): User {
    return this.users.find((user) => user.id === id);
  }

  async signIn(request: SignInRequest): Promise<JwtToken> {
/*    const user = this.findOne('test-id');
    console.log(user);
    if (!user) {
      throw new UnauthorizedException('Invalid credentials');
    }*/

    const payload = { username: 'test', sub: 'test-id'};
    const access_token = await this.jwtService.signAsync(payload);
    return { token: access_token };
  }

  async signOut(username: SignOutRequest): Promise<SignOutResponse> {
    if (username.username === 'valid') return { success: true };
    return { success: false };
  }

  async verifyToken(token: JwtToken): Promise<VerifyTokenResponse> {
    if (token.token === 'valid') return { valid: true };
    return { valid: false };
  }

  async restorePassword(
    userInfo: RestorePasswordRequest,
  ): Promise<RestorePasswordResponse> {
    const user = this.findOne(userInfo.username);

    if (!user || user.username !== 'ldskfsld') {
      throw new UnauthorizedException('Invalid credentials');
    }
    return { success: false };
  }
}
