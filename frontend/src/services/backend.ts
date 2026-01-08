import { computed, Injectable, signal } from '@angular/core';
import PocketBase, { LocalAuthStore, RecordService } from 'pocketbase';
import { environment } from '../environments/environment';
import { Observable } from 'rxjs';
import { Collections, TypedPocketBase, UsersResponse } from '../models';

@Injectable({
  providedIn: 'root'
})
export class BackendService {
  private readonly pocketbase: TypedPocketBase;

  public readonly isLoggedIn;
  private readonly isLoggedIn$ = signal(false);

  private user = computed<UsersResponse>(() => {
    if (this.isLoggedIn()) {
      return this.pocketbase.authStore.record as UsersResponse;
    }

    return {} as UsersResponse;
  });

  constructor() {
    this.isLoggedIn = this.isLoggedIn$.asReadonly();

    this.pocketbase = new PocketBase(environment.backendUrl, new LocalAuthStore());
    this.isLoggedIn$.set(this.pocketbase.authStore.isValid || false);

    this.pocketbase.authStore.onChange(() => {
      this.isLoggedIn$.set(this.pocketbase.authStore.isValid);
    });
  }

  public login(username: string, password: string): Observable<boolean> {
    return new Observable<boolean>((observer) => {
      this.pocketbase
        .collection<UsersResponse>(Collections.Users)
        .authWithPassword(username, password)
        .then(
          () => {
            observer.next(this.pocketbase.authStore.isValid);
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

  public logout() {
    this.pocketbase.authStore.clear();
    this.isLoggedIn$.set(false);
  }

  get currentUser(): UsersResponse {
    return this.user();
  }
}
