import {Component, EventEmitter, Input, Output} from "@angular/core";
import {Message} from "../../../../definitions/email";
import {EmailViewDisplay} from "../../../../definitions/internal";

@Component({
  selector: 'app-email-attachment-list',
  templateUrl: 'email-attachment-list.component.html',
  styleUrls: [
    'email-attachment-list.component.css',
    '../../../../app.component.css',
    '../email-part-browser.component.css',
  ]
})
export class EmailAttachmentListComponent {
  @Input() message!: Message;

  @Output() displayChange = new EventEmitter<EmailViewDisplay>();
  @Input() set display(d: EmailViewDisplay | undefined) {
    this._display = d;
    this.displayChange.emit(d);
  }
  get display(): EmailViewDisplay | undefined {
    return this._display;
  }
  _display?: EmailViewDisplay;

  onClick = () => {
    this._display = {
      displayType: 'attachment',
      message: this.message,
    };
    this.displayChange.emit(this._display);
  }
}
