import {Injectable} from "@angular/core";
import {BehaviorSubject, Observable} from "rxjs";

@Injectable()
export class ToolbarService {
  private _title = new BehaviorSubject<string>('');

  setTitle(title: string) {
    this._title.next(title);
  }

  getTitle(): Observable<string> {
    return this._title.asObservable();
  }
}
