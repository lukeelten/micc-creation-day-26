import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app.config';
import { AppComponent } from './app.component';

bootstrapApplication(AppComponent, {...appConfig, providers: [...appConfig.providers]}).catch((err) => console.error(err));
