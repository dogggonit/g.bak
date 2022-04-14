import {Component, EventEmitter, Input, Output} from "@angular/core";
import {Message} from "../../../../definitions/email";
import {DisplayType, EmailViewDisplay} from "../../../../definitions/internal";

@Component({
  selector: 'app-email-body-list',
  styleUrls: [
    'email-body-list.component.css',
    '../../../../app.component.css',
    '../email-part-browser.component.css',
  ],
  templateUrl: 'email-body-list.component.html',
})
export class EmailBodyListComponent {
  @Input() set message(msg: Message) {
    if (this._message.ID !== msg.ID) {
      this.openedIndexes = {};
    }
    this._message = msg;
  }
  _message: Message = {} as Message;
  openedIndexes: Record<number, boolean> = {};

  @Output() displayChange = new EventEmitter<EmailViewDisplay>();
  @Input() set display(d: EmailViewDisplay | undefined) {
    this._display = d;
    this.displayChange.emit(d);
  }

  get display(): EmailViewDisplay | undefined {
    return this._display;
  }

  _display?: EmailViewDisplay;

  onClick = (msg: Message, type: DisplayType) => {
    this._display = {
      displayType: type,
      message: msg,
    };
    this.displayChange.emit(this._display);
  }

  togglePart(idx: number) {
    this.openedIndexes[idx] = !this.openedIndexes[idx];
  }
}
