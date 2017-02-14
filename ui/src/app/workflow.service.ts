import {Injectable, EventEmitter} from '@angular/core';
import {environment} from '../environments/environment';
import {$WebSocket, WebSocketSendMode} from './websocket.service';

@Injectable()
export class WorkflowService {

  private websocket;

  public events: EventEmitter<WorkflowEvent>;

  constructor() {
    this.events = new EventEmitter();
    this.websocket = new $WebSocket(WorkflowService.url(), null, {
      initialTimeout: 100,
      maxTimeout: 3000,
      reconnectIfNotNormalClose: true
    });
    this.websocket.getDataStream().subscribe((message: MessageEvent) => {
      this.process(JSON.parse(message.data));
    });
    this.websocket.onOpen(() => {
      this.command({command: 'execution-history'});
    });
  }

  private command(cmd: WorkflowCommand) {
    this.websocket.send(JSON.stringify(cmd), WebSocketSendMode.Direct);
  }

  private process(message: any) {
    this.events.emit(message);
  }

  private static url() {
    let path = 'websocket';
    if (environment.wsUrl != null) {
      return 'ws://' + environment.wsUrl + '/' + path;
    } else {
      let url = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
      url += window.location.host;
      url += window.location.pathname.charAt(window.location.pathname.length - 1) == '/' ?
        window.location.pathname + path :
        window.location.pathname + '/' + path;
      return url;
    }
  }
}

export class WorkflowCommand {
  public command: string;
}

export class WorkflowEvent {
  public type: string;
  public payload: any;
}
