import { CommonModule } from "@angular/common";
import { Component, ChangeDetectionStrategy, input, computed, Signal } from "@angular/core";
import { StepperModule } from "primeng/stepper";


@Component({
    selector: 'app-stepper',
    imports: [CommonModule, StepperModule],
    templateUrl: './stepper.html',
    standalone: true,
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class StepperComponent {
  public readonly status = input.required<string>();

  public readonly currentStepIndex: Signal<number> = computed(() => {
    const step = this.steps.find(s => s.key === this.status());
    return step ? step.value : 0;
  });

  public readonly steps = [
    { key: 'CREATED', value: 1 },
    { key: 'SCHEDULED', value: 2 },
    { key: 'PROCESSING', value: 3 },
    { key: 'COMPLETED', value: 4 },
    { key: 'FAILED', value: 4 },
  ];
}