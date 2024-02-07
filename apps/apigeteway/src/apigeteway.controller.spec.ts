import { Test, TestingModule } from '@nestjs/testing';
import { ApigetewayController } from './apigeteway.controller';
import { ApigetewayService } from './apigeteway.service';

describe('ApigetewayController', () => {
  let apigetewayController: ApigetewayController;

  beforeEach(async () => {
    const app: TestingModule = await Test.createTestingModule({
      controllers: [ApigetewayController],
      providers: [ApigetewayService],
    }).compile();

    apigetewayController = app.get<ApigetewayController>(ApigetewayController);
  });

  describe('root', () => {
    it('should return "Hello World!"', () => {
      expect(apigetewayController.getHello()).toBe('Hello World!');
    });
  });
});
