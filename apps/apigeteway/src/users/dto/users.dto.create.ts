import { ApiProperty } from '@nestjs/swagger';
import { CreateUserDto } from '@app/common';

export class CreateUserDtoClass implements CreateUserDto {
  @ApiProperty()
  readonly username: string;

  @ApiProperty()
  readonly password: string;

  @ApiProperty()
  readonly age: number;
}
