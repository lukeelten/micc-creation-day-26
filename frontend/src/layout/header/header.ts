import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, HostBinding, inject } from '@angular/core';
import { ThemeOptions } from '../../services/theme';

@Component({
  selector: 'app-header',
  imports: [
    CommonModule,
  ],
  templateUrl: './header.html',
  styles: [],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class Header {

  public globals = inject(ThemeOptions);


  @HostBinding('class.isActive')
  get isActiveAsGetter() {
    return this.isActive;
  }

  isActive = false;


  toggleSidebar() {
    this.globals.toggleSidebar = !this.globals.toggleSidebar;
    // Clear hover state when toggling
    if (this.globals.toggleSidebar) {
      this.globals.sidebarHover = false;
    }
  }

  toggleSidebarMobile() {
    this.globals.toggleSidebarMobile = !this.globals.toggleSidebarMobile;
  }

  toggleHeaderMobile() {
    this.globals.toggleHeaderMobile = !this.globals.toggleHeaderMobile;
  }
}