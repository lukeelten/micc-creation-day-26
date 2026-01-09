import { inject, Injectable } from "@angular/core";
import { BackendService } from "./backend";
import { RecordService } from "pocketbase";
import { Collections, EventsResponse, RunsResponse } from "src/models";

@Injectable({
  providedIn: 'root'
})
export class RunsRepository {

  private readonly backendService = inject(BackendService);

  private readonly recordService: RecordService<RunsResponse>;


  constructor() {
    this.recordService = this.backendService.getRecordService<RunsResponse>(Collections.Runs);
  }

  public getList() {
    return this.recordService.getFullList({
      sort: '-created',
    });
  }

  public getRecordService() {
    return this.recordService;
  }

  public getById(id: string) {
    return this.recordService.getOne(id);
  }
}

@Injectable({
  providedIn: 'root'
})
export class EventsRepository {

  private readonly backendService = inject(BackendService);
  private readonly recordService: RecordService<EventsResponse>;


  constructor() {
    this.recordService = this.backendService.getRecordService<EventsResponse>(Collections.Events);
  }

  public getForRun(runId: string) {
    return this.recordService.getFullList({
      filter: `run = "${runId}"`,
      sort: '-created',
    });
  }

  public getRecordService() {
    return this.recordService;
  }


}