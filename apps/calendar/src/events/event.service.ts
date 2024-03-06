import { Injectable, OnModuleInit } from '@nestjs/common';
import { Event, ResponseEvents } from '@app/common'
import { randomUUID } from 'crypto';

@Injectable()
export class EventService implements OnModuleInit {
  private readonly events: Event[] = [];

  onModuleInit() {
    for (let i = 0; i < 100; i++) {
      this.create();
    }
  }

  create(): Event {
    const user: Event = {
      user: 'serega',
      body: 'teste',
      endtime: 'tommorow',
      title: 'test',
      starttime: 'now',
      id: randomUUID(),
    };
    this.events.push(user);
    return user;
  }

  findAll({}): ResponseEvents {
    return { result: this.events };
  }
}
