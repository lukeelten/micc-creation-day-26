import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-pagelayout',
  imports: [
    RouterOutlet,
  ],
  templateUrl: './pagelayout.html',
  styles: [],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class Pagelayout {
}