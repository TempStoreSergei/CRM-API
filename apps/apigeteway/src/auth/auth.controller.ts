import {
  Controller,
  Post,
  Body,
  HttpException,
  HttpStatus,
} from '@nestjs/common';
import { AuthService } from './auth.service';
import { LoginDto, RefreshTokensRequest, RegisterDto } from "@app/common";
import { ApiBody, ApiTags } from '@nestjs/swagger';
import { LogInDtoClass, SignUpDtoClass } from './dto';
import { catchError, map, throwError } from 'rxjs';

@ApiTags('Auth')
@Controller('auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @Post('signUp')
  @ApiBody({ type: SignUpDtoClass })
  signUp(@Body() registerDto: RegisterDto) {
    return this.authService.SignUp(registerDto);
  }

  @Post('logIn')
  @ApiBody({ type: LogInDtoClass })
  logIn(@Body() loginDto: LoginDto) {
    return this.authService.LogIn(loginDto).pipe(
      map((result) => {
        const accessToken = result.accessToken;
        const refreshToken = result.refreshToken.token;
        return { accessToken, refreshToken };
      }),
      catchError((error) => {
        console.error('Error logging in:', error);

        if (error instanceof HttpException) {
          return throwError(() => error);
        }

        throw new HttpException(
          'Internal Server Error',
          HttpStatus.INTERNAL_SERVER_ERROR,
        );
      }),
    );
  }

  @Post('refresh')
  @ApiBody({ type: SignUpDtoClass })
  refreshToken(@Body() refreshTokensRequest: RefreshTokensRequest) {
    return this.authService.RefreshTokens(refreshTokensRequest);
  }
}
