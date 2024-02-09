import { Controller, Post, Body } from '@nestjs/common';
import { AuthService } from './auth.service';
import { LoginDto, RegisterDto } from '@app/common';
import { ApiBody, ApiTags } from '@nestjs/swagger';
import { LogInDtoClass, SignUpDtoClass } from './dto';

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
    return this.authService.LogIn(loginDto);
  }
}
