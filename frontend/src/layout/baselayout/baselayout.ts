import { CommonModule } from '@angular/common';
import {ChangeDetectionStrategy, Component, afterNextRender, inject} from '@angular/core';
import { Footer } from '../footer/footer';
import { Header } from '../header/header';
import { Sidebar } from '../sidebar/sidebar';
import { RouterOutlet } from '@angular/router';
import { ThemeOptions } from '../../services/theme';

@Component({
  selector: 'app-base-layout',
  templateUrl: './baselayout.html',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  // Temporarily disable animations to fix jumping issue
  animations: [],
  imports: [
    CommonModule,
    RouterOutlet,
    Header,
    Sidebar,
    Footer
  ]
})
export class BaseLayoutComponent {

  public readonly globals = inject(ThemeOptions);

  constructor() {
    // Use afterNextRender for zoneless compatibility (replaces setTimeout + detectChanges)
    afterNextRender(() => {
      if (typeof window !== 'undefined' && (window as any).bootstrap) {
        // Initialize any Bootstrap tooltips, popovers, etc. that might cause layout shifts
        const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
        tooltipTriggerList.map(function (tooltipTriggerEl) {
          return new (window as any).bootstrap.Tooltip(tooltipTriggerEl);
        });
      }

      // Re-enable animations after layout is stable
      document.body.classList.add('animations-ready');
    });
  }

  toggleSidebarMobile() {
  }
}