import { ChangeDetectionStrategy, Component } from '@angular/core';

import { RouterModule } from '@angular/router';
import { MenuItem } from 'primeng/api';
import { AppMenuitem } from './app.menuitem';

@Component({
    selector: 'app-menu',
    standalone: true,
    imports: [AppMenuitem, RouterModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    template: `<ul class="layout-menu">
          @for (item of model; track item; let i = $index) {
            @if (!item.separator) {
              <li app-menuitem [item]="item" [index]="i" [root]="true"></li>
            }
            @if (item.separator) {
              <li class="menu-separator"></li>
            }
          }
        </ul>`
})
export class AppMenu {
    model: MenuItem[] = [];

    ngOnInit() {
        this.model = [
            {
                label: 'Processing',
                items: [
                    { label: 'Home', icon: 'pi pi-fw pi-home', routerLink: ['/'] },
//                    { label: 'Start Run', icon: 'pi pi-fw pi-play', routerLink: ['/start'] },
                    { label: 'History', icon: 'pi pi-fw pi-history', routerLink: ['/history'] },
                ]
            },
        ];
    }
}
