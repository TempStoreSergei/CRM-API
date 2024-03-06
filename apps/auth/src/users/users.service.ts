import { Injectable, NotFoundException, OnModuleInit } from '@nestjs/common';
import {
  User,
  UpdateUserDto,
  CreateUserDto,
  Users,
  PaginationDto,
} from '@app/common';
import { randomUUID } from 'crypto';
import { Observable, Subject } from 'rxjs';

@Injectable()
export class UsersService implements OnModuleInit {
  private readonly users: User[] = [];

  onModuleInit() {
    for (let i = 0; i < 1000; i++) {
      this.create({ age: 0, username: randomUUID(), password: `password-one` });
    }
  }

  create(createUserDto: CreateUserDto): User {
    const user: User = {
      ...createUserDto,
      subscribed: false,
      socialMedia: {},
      id: randomUUID(),
    };
    this.users.push(user);
    return user;
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
