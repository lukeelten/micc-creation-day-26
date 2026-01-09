import { ChangeDetectionStrategy, Component, ElementRef } from '@angular/core';
import { AppMenu } from './app.menu';

@Component({
    selector: 'app-sidebar',
    standalone: true,
    changeDetection: ChangeDetectionStrategy.OnPush,
    imports: [AppMenu],
    template: ` <div class="layout-sidebar">
        <app-menu></app-menu>
    </div>`
})
export class AppSidebar {
    constructor(public el: ElementRef) {}
}
