import {
  ConflictException,
  Injectable,
  Logger,
  NotFoundException,
} from '@nestjs/common';
import { PaginationDto, UpdateUserDto, User, Users } from '@app/common';
import { Observable, Subject } from 'rxjs';
import { PrismaService } from '../../../user/src/prisma/prisma.service';

@Injectable()
export class UsersService {
  private readonly logger = new Logger(UsersService.name);
  private readonly users = [];
  constructor(private readonly prismaService: PrismaService) {}

  async createProfile(dto: any) {
    // Assuming dto includes userId and optionally other fields like displayName, avatar, etc.
    const existingProfile = await this.prismaService.profile
      .findUnique({
        where: { userId: dto.userId },
      })
      .catch((err) => {
        this.logger.error(
          `Error finding existing profile for user ${dto.userId}: ${err}`,
        );
        throw new NotFoundException('Error finding existing profile');
      });

    if (existingProfile) {
      throw new ConflictException('Profile already exists for this user');
    }

    return this.prismaService.profile
      .create({
        data: {
          userId: dto.userId,
          displayName: dto.displayName,
          avatar: dto.avatar,
          statusMessage: dto.statusMessage,
          phoneNumber: dto.phoneNumber,
        },
      })
      .catch((err) => {
        this.logger.error(`Error creating profile: ${err}`);
        return null;
      });
  }

  findAll(): Users {
    return { users: this.users };
  }

  findOne(id: string): User {
    return this.users.find((user) => user.id === id);
  }

  update(id: string, updateUserDto: UpdateUserDto): User {
    const userIndex = this.users.findIndex((user) => user.id === id);
    if (userIndex !== -1) {
      const userInfo = {
        ...this.users[userIndex],
        ...updateUserDto,
      };
      this.users[userIndex] = userInfo;
      return userInfo;
    }
    throw new NotFoundException(`User not found by id: ${id}`);
  }

  remove(id: string): User {
    const userIndex = this.users.findIndex((user) => user.id === id);
    if (userIndex !== -1) {
      const infoUser = this.users[userIndex];
      this.users.splice(userIndex, 1);
      return infoUser;
    }
    throw new NotFoundException(`This action removes a #${id} user`);
  }

  queryUsers(request: Observable<PaginationDto>): Observable<Users> {
    const subject = new Subject<Users>();

    const onNext = (pagination: PaginationDto) => {
      const start = pagination.page * pagination.skip;
      const users = this.users.slice(start, start + pagination.skip);
      subject.next({ users });
    };

    const onCompleted = () => subject.complete();
    request.subscribe({ next: onNext, complete: onCompleted });

    return subject.asObservable();
  }
}
