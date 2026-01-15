
import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { FloatLabelModule } from 'primeng/floatlabel';
import { InputTextModule } from 'primeng/inputtext';
import { MessageModule } from 'primeng/message';
import { RunsRepository } from 'src/services/repositories';

@Component({
  selector: 'app-start',
  imports: [ReactiveFormsModule, FloatLabelModule, InputTextModule, ButtonModule, MessageModule],
  templateUrl: './start.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class StartComponent {
  runForm = new FormGroup({
    message: new FormControl('', [Validators.required, Validators.minLength(3)])
  });

  private readonly runRepo = inject(RunsRepository);
  private readonly router = inject(Router);

  onSubmit() {
    if (this.runForm.valid) {
      this.runRepo.createRun(this.runForm.value.message!).then((response) => {
        this.router.navigate(['/run', response.id]);
      });
    }
  }
}
