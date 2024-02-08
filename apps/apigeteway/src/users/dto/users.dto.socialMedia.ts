import { ApiProperty } from '@nestjs/swagger';
import { SocialMedia } from '@app/common';

export class SocialMediaUserDtoClass implements SocialMedia {
  @ApiProperty()
  readonly fbUri: string;

  @ApiProperty()
  readonly twitterUri: string;
}
