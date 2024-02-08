import { ApiProperty } from '@nestjs/swagger';
import { UpdateUserDto } from '@app/common';
import { SocialMediaUserDtoClass } from './users.dto.socialMedia';

type UpdateUserDtoWithoutId = Omit<UpdateUserDto, 'id'>;
export class UpdateUserDtoClass implements UpdateUserDtoWithoutId {
  @ApiProperty()
  readonly socialMedia: SocialMediaUserDtoClass | undefined;
}
