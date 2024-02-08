import { ApiProperty } from '@nestjs/swagger';
import { User } from '@app/common';
import { SocialMediaUserDtoClass } from './users.dto.socialMedia';

export class UserDtoClass implements User {
  @ApiProperty()
  readonly id: string;

  @ApiProperty()
  readonly username: string;

  @ApiProperty()
  readonly password: string;

  @ApiProperty()
  readonly age: number;

  @ApiProperty()
  readonly subscribed: boolean;

  @ApiProperty()
  readonly socialMedia: SocialMediaUserDtoClass | undefined;
}
