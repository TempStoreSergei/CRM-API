import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';

@Injectable()
export class MessangerService {
  constructor(private readonly prisma: PrismaService) {}

  async createMessage(data: any) {
    return this.prisma.chat.create({ data });
  }

  async getMessages() {
    return { chat: await this.prisma.chat.findMany() };
  }
}
