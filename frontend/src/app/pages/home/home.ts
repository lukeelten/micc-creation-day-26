import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
    selector: 'app-home',
    imports: [],
    templateUrl: './home.html',
    standalone: true,
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class Home {
    
}
