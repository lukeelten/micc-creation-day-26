import { CommonModule, formatDate } from "@angular/common";
import { ChangeDetectionStrategy, Component, computed, input, Signal } from "@angular/core";
import { CardModule } from "primeng/card";
import { Timeline } from "primeng/timeline";
import { EventsResponse, EventsTypeOptions } from "src/models";

interface TimelineItem {
  title: string;
  created: string;
  icon: string;
  description?: string;
}


@Component({
  selector: 'app-timeline',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, Timeline, CardModule],
  templateUrl: './timeline.html',
  styleUrl: './timeline.scss'
})
export class TimelineComponent {
  public readonly events = input.required<EventsResponse[]>();

  public readonly timelineEvents: Signal<any[]> = computed(() => {
    return this.events().map(event => this.mapEventToTimelineItem(event));
  });

  private mapEventToTimelineItem(event: EventsResponse): TimelineItem {
    return {
      title: event.title,
      created: formatDate(event.created, 'dd.MM.yy HH:mm:ss', 'de-DE'),
      icon: this.calculateIcon(event),
      description: event.description
    };
  }

  private calculateIcon(event: EventsResponse): string {
    // Example logic to determine icon based on event properties
    if (event.type === EventsTypeOptions.info) {
      return 'pi pi-info-circle';
    } else if (event.type === EventsTypeOptions.warn) {
      return 'pi pi-exclamation-triangle';
    } else if (event.type === EventsTypeOptions.error) {
      return 'pi pi-times-circle';
    }

    return 'pi pi-question-circle'; // Default icon
  }


}