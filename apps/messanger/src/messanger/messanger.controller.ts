import { Controller } from '@nestjs/common';
import {
  Chat,
  ChatServiceController,
  CreateChatRequest,
  GetChatRequest,
} from '@app/common/types/messanegr';
import { Observable } from 'rxjs';
import { MessangerService } from './messanger.service';
import { GrpcMethod, GrpcStreamMethod } from '@nestjs/microservices';

@Controller()
export class MessangerController implements ChatServiceController {
  constructor(private readonly messengerService: MessangerService) {}

  @GrpcStreamMethod('ChatService')
  createChat(requestStream: Observable<CreateChatRequest>): Observable<Chat> {
    return new Observable((observer) => {
      const onNext = async (request: CreateChatRequest) => {
        try {
          const createdMessage = await this.messengerService.createMessage({
            mediaUrl: 'ea elit Duis pariatur',
            replyTo: 'Duis',
            text: 'ut elit consequat consectetur dolor',
            type: 'VIDEO',
            userId: 'df44d910-6690-436a-ba13-cd16780e0865',
          });
          observer.next(createdMessage);
        } catch (error) {
          console.error('Error creating message:', error);
          observer.error(error); // Properly handle errors
        }
      };

      const onError = (error: any) => {
        console.error('Stream error:', error);
        observer.error(error);
      };

      const onComplete = () => {
        console.log('Stream completed');
        observer.complete();
      };

      const subscription = requestStream.subscribe({
        next: onNext,
        error: onError,
        complete: onComplete,
      });

      return () => subscription.unsubscribe();
    });
  }
  @GrpcMethod('ChatService')
  getChat(request: GetChatRequest): Observable<Chat> {
    return this.messengerService.getMessages();
  }
}
