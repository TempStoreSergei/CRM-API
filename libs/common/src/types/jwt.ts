/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from "@nestjs/microservices";
import { Observable } from "rxjs";

export const protobufPackageAuth = "auth";

export interface JwtToken {
  token: string;
}

export interface SignInRequest {
  username: string;
  password: string;
}

export interface SignOutRequest {
  username: string;
}

export interface SignOutResponse {
  success: boolean;
}

export interface RestorePasswordRequest {
  username: string;
}

export interface RestorePasswordResponse {
  success: boolean;
}

export interface VerifyTokenRequest {
  token: string;
}

export interface VerifyTokenResponse {
  valid: boolean;
}

export const AUTH_PACKAGE_NAME = "auth";

export interface AuthServiceClient {
  signIn(request: SignInRequest): Observable<JwtToken>;

  signOut(request: SignOutRequest): Observable<SignOutResponse>;

  restorePassword(request: RestorePasswordRequest): Observable<RestorePasswordResponse>;

  verifyToken(request: VerifyTokenRequest): Observable<VerifyTokenResponse>;
}

export interface AuthServiceController {
  signIn(request: SignInRequest): Promise<JwtToken> | Observable<JwtToken> | JwtToken;

  signOut(request: SignOutRequest): Promise<SignOutResponse> | Observable<SignOutResponse> | SignOutResponse;

  restorePassword(
    request: RestorePasswordRequest,
  ): Promise<RestorePasswordResponse> | Observable<RestorePasswordResponse> | RestorePasswordResponse;

  verifyToken(
    request: VerifyTokenRequest,
  ): Promise<VerifyTokenResponse> | Observable<VerifyTokenResponse> | VerifyTokenResponse;
}

export function AuthServiceControllerMethods() {
  return function (constructor: Function) {
    const grpcMethods: string[] = ["signIn", "signOut", "restorePassword", "verifyToken"];
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
