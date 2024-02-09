/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from "@nestjs/microservices";
import { Observable } from "rxjs";

export const protobufPackage = "auth";

export enum Provider {
  GOOGLE = 0,
  UNRECOGNIZED = -1,
}

export enum Role {
  ADMIN = 0,
  USER = 1,
  UNRECOGNIZED = -1,
}

export interface RegisterDto {
  email: string;
  password: string;
}

export interface LoginDto {
  email: string;
  password: string;
  agent: string;
}

export interface UserAuth {
  email: string;
  password: string;
}

export interface RefreshToken {
  token: string;
  exp: string;
  userId: string;
  userAgent: string;
}

export interface Tokens {
  accessToken: string;
  refreshToken: RefreshToken | undefined;
}

export interface DeleteAccountRequest {
  id: string;
  user: JwtPayload | undefined;
}

export interface DeleteAccountResponse {
  id: string;
}

export interface LogOutRequest {
  id: string;
  user: JwtPayload | undefined;
}

export interface RefreshTokensRequest {
  refreshToken: string;
  agent: string;
}

export interface ProviderAuthRequest {
  email: string;
  agent: string;
  provider: Provider;
}

export interface JwtPayload {
  id: string;
  roles: Role[];
  email: string;
}

export interface Empty {
}

export const AUTH_PACKAGE_NAME = "auth";

export interface AuthServiceClient {
  signIn(request: RegisterDto): Observable<UserAuth>;

  logIn(request: LoginDto): Observable<Tokens>;

  deleteAccount(request: DeleteAccountRequest): Observable<DeleteAccountResponse>;

  logOut(request: LogOutRequest): Observable<Empty>;

  refreshTokens(request: RefreshTokensRequest): Observable<Tokens>;

  providerAuth(request: ProviderAuthRequest): Observable<Tokens>;
}

export interface AuthServiceController {
  signIn(request: RegisterDto): Promise<UserAuth> | Observable<UserAuth> | UserAuth;

  logIn(request: LoginDto): Promise<Tokens> | Observable<Tokens> | Tokens;

  deleteAccount(
    request: DeleteAccountRequest,
  ): Promise<DeleteAccountResponse> | Observable<DeleteAccountResponse> | DeleteAccountResponse;

  logOut(request: LogOutRequest): Promise<Empty> | Observable<Empty> | Empty;

  refreshTokens(request: RefreshTokensRequest): Promise<Tokens> | Observable<Tokens> | Tokens;

  providerAuth(request: ProviderAuthRequest): Promise<Tokens> | Observable<Tokens> | Tokens;
}

export function AuthServiceControllerMethods() {
  return function (constructor: Function) {
    const grpcMethods: string[] = ["signIn", "logIn", "deleteAccount", "logOut", "refreshTokens", "providerAuth"];
    for (const method of grpcMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcMethod("AuthService", method)(constructor.prototype[method], method, descriptor);
    }
    const grpcStreamMethods: string[] = [];
    for (const method of grpcStreamMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcStreamMethod("AuthService", method)(constructor.prototype[method], method, descriptor);
    }
  };
}

export const AUTH_SERVICE_NAME = "AuthService";
