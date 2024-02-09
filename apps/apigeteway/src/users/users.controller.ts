import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
  UseGuards,
} from '@nestjs/common';
import { UsersService } from './users.service';
import { CreateUserDto, UpdateUserDto } from '@app/common';
import { ApiBody, ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { CreateUserDtoClass, UpdateUserDtoClass, UserDtoClass } from './dto';
import { JwtAuthGuard } from '../guard/premissions.guard';

@ApiTags('Users')
@Controller('users')
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  @Post()
  @UseGuards(JwtAuthGuard)
  @ApiBody({ type: CreateUserDtoClass })
  @ApiResponse({
    status: 201,
    description: 'The record has been successfully created.',
    type: UserDtoClass,
  })
  create(@Body() createUserDto: CreateUserDto) {
    return this.usersService.create(createUserDto);
  }

  @Get()
  @UseGuards(JwtAuthGuard)
  @ApiResponse({
    status: 201,
    description: 'The record has been successfully created.',
    type: [UserDtoClass],
  })
  findAll() {
    return this.usersService.findAll();
  }

  @Get(':id')
  @UseGuards(JwtAuthGuard)
  @ApiResponse({
    status: 201,
    description: 'The record has been successfully created.',
    type: UserDtoClass,
  })
  findOne(@Param('id') id: string) {
    return this.usersService.findOne(id);
  }

  @Patch(':id')
  @UseGuards(JwtAuthGuard)
  @ApiOperation({ summary: 'Update a user' })
  @ApiBody({ type: UpdateUserDtoClass })
  @ApiResponse({
    status: 200,
    description: 'The user has been successfully updated.',
    type: UserDtoClass,
  })
  update(@Param('id') id: string, @Body() updateUserDto: UpdateUserDto) {
    return this.usersService.update(id, updateUserDto);
  }

  @Delete(':id')
  @UseGuards(JwtAuthGuard)
  @ApiResponse({
    status: 200,
    description: 'The user has been successfully updated.',
    type: UserDtoClass,
  })
  remove(@Param('id') id: string) {
    return this.usersService.remove(id);
  }

  @Post('email')
  emailUser() {
    return this.usersService.emailUsers();
  }
}
