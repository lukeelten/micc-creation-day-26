import { inject, Injectable } from "@angular/core";
import { BackendService } from "./backend";
import { RecordService } from "pocketbase";
import { Collections, EventsResponse, RunsRecord, RunsResponse, RunsStatusOptions, UsersResponse } from "src/models";

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

  public getAuthorName(authorId: string): Promise<string> {
    return this.backendService.getRecordService<UsersResponse>(Collections.Users).getOne(authorId).then(user => user.name);
  }

  public createRun(message: string): Promise<RunsResponse> {
    let data: any = {
      message: message,
      status: RunsStatusOptions.CREATED,
      runtimeSeconds: 0,
      author: this.backendService.currentUser.id
    };

    return this.recordService.create(data);
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