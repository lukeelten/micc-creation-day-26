import { ApplicationConfig, provideBrowserGlobalErrorListeners, provideZonelessChangeDetection } from '@angular/core';
import { provideRouter, withComponentInputBinding, withEnabledBlockingInitialNavigation, withInMemoryScrolling, withViewTransitions } from '@angular/router';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';

import { routes } from './app.routes';
import { provideHttpClient, withFetch } from '@angular/common/http';


export const appConfig: ApplicationConfig = {
  providers: [
    provideAnimationsAsync(),
    provideBrowserGlobalErrorListeners(),
    provideZonelessChangeDetection(),
    provideRouter(routes, withViewTransitions(), withInMemoryScrolling(), withComponentInputBinding(), withEnabledBlockingInitialNavigation()),
    provideHttpClient(withFetch()),
  ]
};
