import {
  Injectable,
  CanActivate,
  ExecutionContext,
  UnauthorizedException,
} from '@nestjs/common';
import { Observable } from 'rxjs';

@Injectable()
export class PermissionsGuard implements CanActivate {
  canActivate(
    context: ExecutionContext,
  ): boolean | Promise<boolean> | Observable<boolean> {
    const request = context.switchToHttp().getRequest();

    // Add the authorization logic here with Permit.io.
    // If the user has the necessary permissions, it will return true.
    // If the user does not have the necessary permissions,
    // throw an UnauthorizedException.

    const userHasPermission = false; // <- replace this line with your Permit.io logic

    if (!userHasPermission) {
      throw new UnauthorizedException(
        'You do not have the necessary permissions.',
      );
    }
    return true;
  }
}
