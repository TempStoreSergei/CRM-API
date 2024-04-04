import { Controller, OnModuleInit } from '@nestjs/common';
import { GrpcMethod, GrpcStreamMethod } from '@nestjs/microservices';
import { MessageManager } from '../cassandra/messageManager';
import { ConversationManager } from '../cassandra/conversationManager';
import {
  ChatMessage,
  GetHistoryRequest,
  GetHistoryResponse,
  MessageServiceController,
  ReadHistoryRequest,
} from '@app/common/types/messenger';
import { catchError, Observable, Subject, takeUntil } from 'rxjs';

@Controller()
export class MessengerController
  implements MessageServiceController, OnModuleInit
{
  onModuleInit(): any {
    this.messageManager.createSchema('messages_new');
  }

  constructor(
    private readonly messageManager: MessageManager,
    private readonly conversationManager: ConversationManager,
  ) {}

  @GrpcMethod('MessageService', 'GetHistory')
  getHistory(
    request: GetHistoryRequest,
  ):
    | Promise<GetHistoryResponse>
    | Observable<GetHistoryResponse>
    | GetHistoryResponse {
    return this.messageManager.getMessages(
      request.conversationId,
      request.limit,
      request.offset,
    );
  }

  @GrpcMethod('MessageService', 'ReadHistory')
  readHistory(request: ReadHistoryRequest) {
    return 'Not implemented';
  }

  @GrpcStreamMethod('MessageService', 'Communicate')
  communicate(upstream: Observable<ChatMessage>): Observable<ChatMessage> {
    const downstream = new Subject<ChatMessage>();
    const end$ = new Subject<void>();

    upstream
      .pipe(
        takeUntil(end$),
        catchError((error, caught) => {
          // Log the error or handle it
          console.error('Error in communication:', error);
          return caught;
        }),
      )
      .subscribe({
        next: async (message: ChatMessage) => {
          try {
            const userId = '79bf07a8-07e7-4355-b692-331f54fc8868';
            console.log('Received message:', message);
            await this.messageManager.saveMessage(message, userId);

            const members = await this.conversationManager.getByID(
              message.conversationId,
            );
            members.forEach((memberId) => {
              if (memberId !== userId) {
                downstream.next(message);
              }
            });
          } catch (error) {
            console.error('Failed to process message:', error);
          }
        },
        complete: () => {
          console.log('Stream completed');
          end$.next();
          end$.complete();
        },
        error: (error) => console.error('Stream error:', error),
      });

    return downstream.asObservable();
  }
}
