import { CanActivateFn, Router } from '@angular/router';
import { inject } from '@angular/core';
import { BackendService } from './backend';

export const isAuthenticatedRouteGuard: CanActivateFn = () => {
  const backendService = inject(BackendService);
  if (backendService.isLoggedIn()) {
    return true;
  }

  const router = inject(Router);
  return router.parseUrl('/login');
};
