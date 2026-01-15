import { ChangeDetectionStrategy, Component, effect, inject, OnInit, signal } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { CheckboxModule } from 'primeng/checkbox';
import { InputTextModule } from 'primeng/inputtext';
import { PasswordModule } from 'primeng/password';
import { RippleModule } from 'primeng/ripple';
import { AppFloatingConfigurator } from '../../layout/component/app.floatingconfigurator';
import { BackendService } from 'src/services/backend';

@Component({
    selector: 'app-login',
    standalone: true,
    changeDetection: ChangeDetectionStrategy.OnPush,
    imports: [ButtonModule, CheckboxModule, InputTextModule, PasswordModule, FormsModule, RouterModule, RippleModule, AppFloatingConfigurator, ReactiveFormsModule],
    templateUrl: './login.html',
})
export class Login implements OnInit {
  public readonly loginForm = new FormGroup({
    username: new FormControl('', [Validators.required]),
    password: new FormControl('', [Validators.required, Validators.minLength(6)]),
    rememberme: new FormControl(false)
  });

  public readonly errorMessage = signal('');

  private readonly backendService: BackendService = inject(BackendService);
  private readonly router: Router = inject(Router);

  constructor() {
    effect(() => {
      if (this.backendService.isLoggedIn()) {
        this.router.navigate(['/']).catch(console.error);
      }
    });
  }

  ngOnInit(): void {
    if (this.backendService.isLoggedIn()) {
      this.router.navigate(['/']).catch(console.error);
    }
  }

  login() {
    if (this.loginForm.valid && this.loginForm.value.username && this.loginForm.value.password) {
      this.backendService.login(this.loginForm.value.username, this.loginForm.value.password).subscribe((isLoggedIn) => {
        if (isLoggedIn) {
          this.router.navigate(['/']).catch(console.error);
        } else {
          console.log('error');
          this.errorMessage.set('Login failed. Please check your credentials.');
        }
      });
    } else {
      this.errorMessage.set('Please enter a valid username and password.');
    }
  }
}
