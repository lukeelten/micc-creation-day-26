import { Routes } from '@angular/router';
import { AppLayout } from './app/layout/component/app.layout';
import { Notfound } from './app/pages/notfound/notfound';
import { Login } from './app/pages/login/login';
import { Home } from '@/pages/home/home';
import { ViewRun } from '@/pages/view/view';
import { History } from '@/pages/history/history';
import { isAuthenticatedRouteGuard } from './services/guards';
import { StartComponent } from '@/pages/start/start';

export const appRoutes: Routes = [
    {
        path: '',
        component: AppLayout,
        canActivate: [isAuthenticatedRouteGuard],
        children: [
            { path: '', component: Home },
            { path: 'start', component: StartComponent },
            { path: 'run/:id', component: ViewRun },
            { path: 'history', component: History }
        ]
    },
    { path: 'notfound', component: Notfound },
    { path: 'login', component: Login },
    { path: '**', redirectTo: '/notfound' }
];
