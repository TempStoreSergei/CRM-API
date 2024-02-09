import { ApiProperty } from '@nestjs/swagger';
import { LoginDto, RegisterDto } from '@app/common';

export class LogInDtoClass implements LoginDto {
  @ApiProperty()
  readonly email: string;

  @ApiProperty()
  readonly password: string;

  @ApiProperty()
  readonly agent: string;
}
