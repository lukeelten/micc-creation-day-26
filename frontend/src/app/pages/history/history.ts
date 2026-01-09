import { TableColumn, TableComponent } from '@/shared/table/table';
import { formatDate } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, OnDestroy, OnInit, resource, Signal } from '@angular/core';
import { Router } from '@angular/router';
import { UnsubscribeFunc } from 'pocketbase';
import { RecordIdString, RunsResponse } from 'src/models';
import { RunsRepository } from 'src/services/repositories';

@Component({
    selector: 'app-history',
    imports: [TableComponent],
    templateUrl: './history.html',
    standalone: true,
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class History implements OnInit, OnDestroy {
    private readonly runRepo = inject(RunsRepository);
    private readonly router = inject(Router);
    private readonly unsubscribeFuncs: Array<UnsubscribeFunc> = [];

    public readonly runResource = resource({
        loader: () => this.loadRuns()
    });

    public readonly runs: Signal<RunsResponse[]> = computed(() => {
        if (this.runResource.hasValue()) {
            return this.runResource.value();
        }

        return [];
    });


    ngOnInit(): void {
        this.runRepo.getRecordService().subscribe('*', (e) => {
            this.runResource.reload();
        }).then((unsubscribe) => {
            this.unsubscribeFuncs.push(unsubscribe);
        });
    }

    ngOnDestroy(): void {
        this.unsubscribeFuncs.forEach((func) => func());
    }

    public click(recordId: RecordIdString): void {
        console.log(recordId);
        this.router.navigate(['/run', recordId]);
    }

    private loadRuns(): Promise<RunsResponse[]> {
        return this.runRepo.getList();
    }

    public readonly columns: Array<TableColumn> = [
        {
            field: 'id',
            label: 'ID',
            sortable: true
        },
        {
            field: 'created',
            label: 'Created At',
            transform: (value) => formatDate(value, 'dd.MM.yyyy HH:mm', 'de-DE'),
            sortable: true
        },
        {
            field: 'message',
            label: 'Description',
            searchable: true,
            sortable: true
        },
        {
            field: 'status',
            label: 'Status',
            searchable: true
        },
        // @todo Author
    ];
}
