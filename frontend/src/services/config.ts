import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

export interface AppConfig {
  backendUrl: string;
  production: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class ConfigService {
  private readonly http = inject(HttpClient);

  private backendUrl = '';

  public getBackendUrl(): Promise<string> {
    return new Promise((resolve) => {
      if (this.backendUrl && this.backendUrl.length > 0) {
        resolve(this.backendUrl);
      } else {
        this.http.get<AppConfig>('/config.json').subscribe((result) => {
          this.backendUrl = result.backendUrl;
          resolve(this.backendUrl);
        });
      }
    });
  }
}
