import {
  ConflictException,
  ForbiddenException,
  HttpException,
  HttpStatus,
  Inject,
  Injectable,
  Logger,
  UnauthorizedException,
} from '@nestjs/common';
import { add } from 'date-fns';
import { v4 } from 'uuid';
import { JwtService } from '@nestjs/jwt';
import { Role, Token, User } from '@prisma/client';
import {
  JwtPayload,
  LoginDto,
  Provider,
  RegisterDto,
  Tokens,
} from '@app/common';
import { genSaltSync, hashSync, compareSync } from 'bcrypt';
import { Cache } from 'cache-manager';
import { CACHE_MANAGER } from '@nestjs/cache-manager';
import { PrismaService } from '../prisma/prisma.service';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class AuthService {
  private readonly logger = new Logger(AuthService.name);
  constructor(
    private readonly prismaService: PrismaService,
    @Inject(CACHE_MANAGER) private cacheManager: Cache,
    private readonly configService: ConfigService,
    private readonly jwtService: JwtService,
  ) {}

  async signIn(dto: RegisterDto) {
    const user: User = await this.findOne(dto.email).catch((err) => {
      this.logger.error(err);
      return null;
    });
    if (user) {
      throw new ConflictException(
        'Пользователь с таким email уже зарегистрирован',
      );
    }
    return this.save(dto).catch((err) => {
      this.logger.error(err);
      return null;
    });
  }

  async logIn(dto: LoginDto, agent: string): Promise<Tokens> {
    const user: User = await this.findOne(dto.email, true).catch((err) => {
      this.logger.error(err);
      return null;
    });
    if (!user || !compareSync(dto.password, user.password)) {
      throw new UnauthorizedException('Не верный логин или пароль');
    }
    return this.generateTokens(user, agent);
  }

  async deleteAccount(id: string, user: JwtPayload) {
    if (user.id !== id && !user.roles.includes(Role.ADMIN)) {
      throw new ForbiddenException();
    }
    await Promise.all([
      this.cacheManager.del(id),
      this.cacheManager.del(user.email),
    ]);
    return this.prismaService.user.delete({
      where: { id },
      select: { id: true },
    });
  }

  async logOut(id: string, user: JwtPayload) {
    if (user.id !== id) {
      throw new ForbiddenException();
    }
    await Promise.all([
      this.cacheManager.del(id),
      this.cacheManager.del(user.email),
    ]);
    return this.deleteRefreshToken(id);
  }

  async refreshTokens(refreshToken: string, agent: string): Promise<Tokens> {
    const token = await this.prismaService.token.delete({
      where: { token: refreshToken },
    });
    if (!token || new Date(token.exp) < new Date()) {
      throw new UnauthorizedException('Invalid or expired refresh token');
    }
    const user = await this.findOne(token.userId);
    return this.generateTokens(user, agent);
  }

  async providerAuth(email: string, agent: string, provider: Provider) {
    const userExists = await this.findOne(email);
    if (userExists) {
      const user = await this.save({ email, provider }).catch((err) => {
        this.logger.error(err);
        return null;
      });
      return this.generateTokens(user, agent);
    }
    const user = await this.save({ email, provider }).catch((err) => {
      this.logger.error(err);
      return null;
    });
    if (!user) {
      throw new HttpException(
        `Не получилось создать пользователя с email ${email} в Google auth`,
        HttpStatus.BAD_REQUEST,
      );
    }

    return this.generateTokens(user, agent);
  }

  private deleteRefreshToken(token: string) {
    return this.prismaService.token.delete({ where: { token } });
  }

  private async generateTokens(user: User, agent: string): Promise<Tokens> {
    const accessToken =
      'Bearer ' +
      this.jwtService.sign({
        id: user.id,
        roles: user.roles,
      });
    const refreshToken = await this.getRefreshToken(user.id, agent);
    return { accessToken, refreshToken };
  }

  private hashPassword(password: string) {
    return hashSync(password, genSaltSync(10));
  }

  private async getRefreshToken(userId: string, agent: string): Promise<Token> {
    const existingToken = await this.prismaService.token.findFirst({
      where: {
        userId,
        userAgent: agent,
      },
    });

    if (existingToken) {
      return this.prismaService.token.update({
        where: {
          token: existingToken.token,
        },
        data: {
          token: v4(),
          exp: add(new Date(), { months: 1 }).toISOString(),
        },
      });
    }

    return this.prismaService.token.create({
      data: {
        token: v4(),
        exp: add(new Date(), { months: 1 }).toISOString(),
        userId,
        userAgent: agent,
      },
    });
  }

  private async save(user: Partial<User>) {
    const hashedPassword = user?.password
      ? this.hashPassword(user.password)
      : null;
    const savedUser = await this.prismaService.user.upsert({
      where: {
        email: user.email,
      },
      update: {
        password: hashedPassword ?? undefined,
        provider: user?.provider ?? undefined,
        rolesId: "1",
        isBlocked: user?.isBlocked ?? undefined,
      },
      create: {
        email: user.email,
        password: hashedPassword,
        provider: user?.provider,
        rolesId: "1",
      },
    });
    await this.cacheManager.set(savedUser.id, savedUser);
    await this.cacheManager.set(savedUser.email, savedUser);
    return savedUser;
  }

  private async findOne(idOrEmail: string, isReset = false): Promise<User> {
    if (isReset) {
      await this.cacheManager.del(idOrEmail);
    }
    const user = await this.cacheManager.get<User>(idOrEmail);
    if (!user) {
      const user = await this.prismaService.user.findFirst({
        where: {
          OR: [{ id: idOrEmail }, { email: idOrEmail }],
        },
      });
      if (!user) {
        return null;
      }
      await this.cacheManager.set(
        idOrEmail,
        user,
        this.configService.get('JWT_EXP'),
      );
      return user;
    }
    return user;
  }
}
