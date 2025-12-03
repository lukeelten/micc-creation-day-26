import { CommonModule } from '@angular/common';
import { afterNextRender, ChangeDetectionStrategy, Component, HostListener, inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Observable } from 'rxjs';
import { ThemeOptions } from '../../services/theme';

@Component({
  selector: 'app-sidebar',
  imports: [
    CommonModule
  ],
  templateUrl: './sidebar.html',
  styleUrls: ['./sidebar.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class Sidebar {
  public extraParameter: any;
  public openMenus: string[] = [];
  
  // Supported menu types: dashboardsMenu, pagesMenu, elementsMenu, componentsMenu, 
  // tablesMenu, formsMenu, chartsMenu, widgetsMenu

  public readonly globals = inject(ThemeOptions);

  constructor(
    private activatedRoute: ActivatedRoute,
  ) {

    // Use afterNextRender for zoneless compatibility (replaces setTimeout)
    afterNextRender(() => {
      this.innerWidth = window.innerWidth;
      if (this.innerWidth < 1200) {
        this.globals.toggleSidebar = true;
      }
    });
  }

  private newInnerWidth = 0;
  private innerWidth = 0;
  activeId = 'dashboardsMenu';

  toggleSidebar() {
    this.globals.toggleSidebar = !this.globals.toggleSidebar;
    // If we're closing the sidebar, also clear the hover state
    if (this.globals.toggleSidebar) {
      this.globals.sidebarHover = false;
    }
  }

  onSidebarMouseEnter() {
    // Only allow hover to open sidebar if it's in collapsed state
    if (this.globals.toggleSidebar) {
      this.globals.sidebarHover = true;
    }
  }

  onSidebarMouseLeave() {
    // Only remove hover state if sidebar is in collapsed state
    if (this.globals.toggleSidebar) {
      this.globals.sidebarHover = false;
    }
  }

  ngOnInit() {
    // Get the extraParameter from the route to determine which menu should be open
    this.extraParameter = this.activatedRoute.snapshot.firstChild?.data['extraParameter'];

    // Initialize open menus based on current route
    if (this.extraParameter) {
      this.openMenus = [this.extraParameter];
    }
  }

  toggleSubmenu(menuId: string) {
    // Toggle submenu: close if open, open if closed (and close all others)
    const index = this.openMenus.indexOf(menuId);
    if (index > -1) {
      this.openMenus.splice(index, 1);
    } else {
      this.openMenus = [menuId]; // Close others and open this one
    }
  }

  onNavigate() {
    // Close sidebar on mobile when navigating
    if (window.innerWidth < 1200) {
      this.globals.toggleSidebarMobile = true;
      this.globals.sidebarHover = false;
    }
  }

  @HostListener('window:resize', ['$event'])
  onResize(event: Event) {
    this.newInnerWidth = (event.target as Window).innerWidth;

    if (this.newInnerWidth < 1200) {
      this.globals.toggleSidebar = true;
    } else {
      this.globals.toggleSidebar = false;
    }

  }
}