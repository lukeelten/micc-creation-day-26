import { computed, Injectable, Signal, signal, inject } from '@angular/core';
import PocketBase, { LocalAuthStore, RecordService } from 'pocketbase';
import { Observable } from 'rxjs';
import { Collections, TypedPocketBase, UsersResponse } from '../models';
import { ConfigService } from './config';

@Injectable({
  providedIn: 'root'
})
export class BackendService {
  private readonly configService = inject(ConfigService);
  private pocketbase: TypedPocketBase | null = null;

  public readonly isLoggedIn: Signal<boolean>;
  private readonly isLoggedIn$ = signal(false);

  private user = computed<UsersResponse>(() => {
    if (!this.pocketbase) {
      return {} as UsersResponse;
    }

    if (this.isLoggedIn()) {
      return this.pocketbase.authStore.record as UsersResponse;
    }

    return {} as UsersResponse;
  });

  constructor() {
    this.isLoggedIn = this.isLoggedIn$.asReadonly();

    this.configService.getBackendUrl().then((backendUrl) => {
      this.pocketbase = new PocketBase(backendUrl, new LocalAuthStore());
      this.isLoggedIn$.set(this.pocketbase.authStore.isValid || false);

      this.pocketbase.authStore.onChange(() => {
        this.isLoggedIn$.set(this.pocketbase?.authStore.isValid || false);
      });
    });
  }

  public login(username: string, password: string): Observable<boolean> {
    return new Observable<boolean>((observer) => {
      if (!this.pocketbase) {
        observer.error('PocketBase is not initialized yet.');
        observer.complete();
        return;
      }

      this.pocketbase
        .collection<UsersResponse>(Collections.Users)
        .authWithPassword(username, password)
        .then(
          () => {
            observer.next(this.pocketbase?.authStore.isValid || false);
            observer.complete();
          },
          (err: any) => {
            console.error('Login failed', err);
            observer.error(err);
            observer.complete();
          }
        );
    });
  }

  public getRecordService<T>(collectionName: string): RecordService<T> {
    if (!this.pocketbase) {
      throw new Error('PocketBase is not initialized yet.');
    }

    return this.pocketbase.collection<T>(collectionName);
  }

  public logout() {
    if (!this.pocketbase) {
      throw new Error('PocketBase is not initialized yet.');
    }

    this.pocketbase.authStore.clear();
    this.isLoggedIn$.set(false);
  }

  get currentUser(): UsersResponse {
    return this.user();
  }
}
