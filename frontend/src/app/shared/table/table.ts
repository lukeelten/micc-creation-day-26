import { ChangeDetectionStrategy, Component, computed, input, model, output, viewChild } from '@angular/core';
import { SortableColumn, Table, TableModule } from 'primeng/table';
import { InputTextModule } from 'primeng/inputtext';
import { FormsModule } from '@angular/forms';
import { InputIcon } from 'primeng/inputicon';
import { IconField } from 'primeng/iconfield';
import { RecordIdString } from 'src/models';

@Component({
  selector: 'app-table',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [TableModule, InputTextModule, FormsModule, InputIcon, IconField],
  templateUrl: './table.html',
})
export class TableComponent {
  // Inputs
  data = input.required<any>();
  dataKey = input<string>('id');
  rowsPerPageOptions = input<number[]>([5, 10, 20]);
  columns = input.required<TableColumn[]>();
  searchable = input<boolean>(false);
  title = input.required<string>();
  rows = input<number>(15);
  loading = input(false);

  // Outputs
  click = output<RecordIdString>();

  // Models
  searchString = model('');

  selectedItem: any = null;

  // Viewchilds
  table = viewChild.required(Table);

  // attributes
  totalRecords = computed(() => this.data().length);
  searchableFields = computed(() => {
    let columns = this.columns();
    let result: Array<string> = [];
    for (let column of columns) {
      if (column.searchable) {
        result.push(column.field);
      }
    }

    return result;
  });

  onRowClick(event: PointerEvent, item: any) {
    event.preventDefault();
    event.stopPropagation();

    if (!event || !item) {
      return;
    }

    if (item === '') {
      return;
    }

    const id = item[this.dataKey()] as RecordIdString;
    this.click.emit(id);
  }

  onSearchChange() {
    this.table().filterGlobal(this.searchString(), 'contains');
  }
}

export interface TableColumn {
  field: string;
  label: string;

  transform?: (value: any) => string | null;
  sortable?: boolean;
  searchable?: boolean;
}