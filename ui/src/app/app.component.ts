import {Component} from '@angular/core';
import {WorkflowService} from './workflow.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  providers: [WorkflowService]
})
export class AppComponent {
  maxLength = 8;
  execution = null;
  executions: Execution[] = [];

  static newExecution(payload: any): Execution {
    let execution = new Execution(payload.id, payload.start);
    if (payload.status != null && payload.status >= 0) {
      execution.finish = payload.finish;
      execution.status = payload.status;
    }
    if (payload.log != null) {
      for (let index = 0; index < payload.log.length; index++) {
        execution.log.push(AppComponent.newExecutionLog(payload.log[index]));
      }
    }
    execution.calculate();
    return execution;
  }

  static newExecutionLog(payload: any): ExecutionLog {
    return new ExecutionLog(
      payload.message,
      payload.error,
      payload.timestamp
    )
  }

  constructor(private workflow: WorkflowService) {
    workflow.events.subscribe(event => {
      if (event.type == 'execution-start') {
        this.addExecution(AppComponent.newExecution(event.payload));
      } else if (event.type == 'execution-finish') {
        this.finishExecution(AppComponent.newExecution(event.payload));
      } else if (event.type == 'execution-log') {
        this.executionLog(event.payload.execution, AppComponent.newExecutionLog(event.payload.log));
      } else if (event.type == 'execution') {
        this.addExecution(AppComponent.newExecution(event.payload))
      }
    });
  }

  public show(execution: Execution) {
    this.execution = execution;
  }

  private addExecution(execution: Execution) {
    let index = 0;
    while (index < this.executions.length && this.executions[index].id > execution.id) index++;
    if (index >= this.executions.length || this.executions[index].id != execution.id) {
      this.executions.splice(index, 0, execution);
      while (this.executions.length > this.maxLength) this.executions.pop();
    }
  }

  private finishExecution(execution: Execution) {
    for (let index = 0; index < this.executions.length; index++) {
      if (this.executions[index].id == execution.id) {
        this.executions[index].finish = execution.finish;
        this.executions[index].status = execution.status;
        this.executions[index].calculate();
        return;
      }
    }
  }

  private executionLog(id: number, log: ExecutionLog) {
    for (let index = 0; index < this.executions.length; index++) {
      if (this.executions[index].id == id) {
        this.executions[index].log.push(log);
        return;
      }
    }
  }
}

class Execution {
  public finish: string;
  public status: number;
  public elapsed: number = -1;
  public log: ExecutionLog[] = [];

  constructor(public id: number, public start: string) {
  }

  public calculate(): void {
    this.elapsed = this.finish != null ? Date.parse(this.finish) - Date.parse(this.start) : -1;
  }
}

class ExecutionLog {
  public tag: string;

  constructor(public message: string, public error: boolean, public timestamp: string) {
    this.extract('API OPTIONS:', 'API OPTIONS');
    this.extract('API GET');
    this.extract('API PUT');
    this.extract('API POST');
    this.extract('API DELETE');
    this.extract('HTTP REQUEST \\[\\d+\\]');
    this.extract('HTTP RESPONSE \\[\\d+\\]');
    this.extract('ELASTICSEARCH COUNT');
    this.extract('ELASTICSEARCH AVERAGE');
  }

  private extract(prefix: string, tag?: string) {
    let match = new RegExp('^(' + prefix + ')(.*)', 'i').exec(this.message);
    if (match && match.length > 2) {
      this.message = match[2];
      this.tag = tag ? tag : match[1];
    }
  }
}
