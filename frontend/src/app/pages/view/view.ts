import { StepperComponent } from '@/shared/stepper/stepper';
import { TimelineComponent } from '@/shared/timeline/timeline';
import { ChangeDetectionStrategy, Component, computed, effect, inject, OnDestroy, OnInit, resource, Signal, signal } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { UnsubscribeFunc } from 'pocketbase';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { EventsResponse, RunsResponse } from 'src/models';
import { EventsRepository, RunsRepository } from 'src/services/repositories';
import { formatDate } from '@angular/common';

@Component({
    selector: 'app-view',
    imports: [StepperComponent, TimelineComponent, ProgressSpinnerModule],
    templateUrl: './view.html',
    standalone: true,
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class ViewRun implements OnInit, OnDestroy {
    private readonly unsubscribeFuncs: Array<UnsubscribeFunc> = [];

    private readonly runsRepo = inject(RunsRepository);
    private readonly eventsRepo = inject(EventsRepository);
    private readonly route = inject(ActivatedRoute);

    public readonly runId = signal<string>('');
    public readonly authorName = signal<string>('');

    public readonly runResource = resource({
        params: () => ({id: this.runId()}),
        loader: ({params}) => this.loadRun(params.id)
    });

    public readonly eventsResource = resource({
        params: () => ({id: this.runId()}),
        loader: ({params}) => this.loadEvents(params.id)
    });

    public readonly run: Signal<RunsResponse | null> = computed(() => {
        if (this.runResource.hasValue()) {
            return this.runResource.value();
        }

        return null;
    });

    public readonly runStatus = computed(() => {
        const run = this.run();
        if (!run || !run.status) {
            return 'UNKNOWN';
        }

        return run.status;
    });
    
    public readonly events: Signal<EventsResponse[]> = computed(() => {
        if (this.eventsResource.hasValue()) {
            return this.eventsResource.value();
        }

        return [];
    });

    constructor() {
        effect(() => {
            const run = this.run();
            if (run && run.author) {
                this.runsRepo.getAuthorName(run.author).then(name => {
                    this.authorName.set(name);
                });
            }
        });
    }
    

    ngOnInit(): void {
        this.route.params.subscribe({
            next: (params) => {
                const runId = params['id'];
                this.runId.set(runId);

                this.eventsRepo.getRecordService().subscribe('*', (e) => {
                    this.eventsResource.reload();
                }).then((unsubscribe) => {
                    this.unsubscribeFuncs.push(unsubscribe);
                });

                this.runsRepo.getRecordService().subscribe('*', (e) => {
                    this.runResource.reload();
                }).then((unsubscribe) => {
                    this.unsubscribeFuncs.push(unsubscribe);
                });
            }
        });
    }

    ngOnDestroy(): void {
        this.unsubscribeFuncs.forEach((func) => func());
    }

    public formatDate(dateString: string | null | undefined, formatStr: string, locale: string): string {
        if (!dateString) {
            return '';
        }

        return formatDate(dateString, formatStr, locale);
    }

    private loadRun(runId: string): Promise<RunsResponse> {
        if (!runId || runId.length === 0) {
            return Promise.reject('No run ID set');
        }
        
        return this.runsRepo.getById(runId);
    }

    private loadEvents(runId: string): Promise<EventsResponse[]> {
        if (!runId || runId.length === 0) {
            return Promise.resolve([]);
        }

        return this.eventsRepo.getForRun(runId);
    }
}
