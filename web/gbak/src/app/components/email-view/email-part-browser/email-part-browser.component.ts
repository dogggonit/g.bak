import {
  AfterContentInit,
  AfterViewChecked,
  AfterViewInit,
  Component,
  EventEmitter,
  Input,
  OnInit,
  Output
} from "@angular/core";
import {Attachment, Message, MessageBody} from "../../../definitions/email";
import {EmailViewDisplay} from "../../../definitions/internal";
import {max, Observable, Subject} from "rxjs";

@Component({
  selector: 'email-part-browser',
  templateUrl: 'email-part-browser.component.html',
  styleUrls: [
    '../../../app.component.css',
    'email-part-browser.component.css',
  ]
})
export class EmailPartBrowserComponent {
  @Input() message!: Message;
  @Input() set display(d: EmailViewDisplay | undefined) {
    this._display = d;
    this.displayChange.emit(d);
  }
  @Output() displayChange = new EventEmitter<EmailViewDisplay>();
  _display?: EmailViewDisplay;
  get display(): EmailViewDisplay | undefined {
    return this._display;
  }

  get nextId(): number {
    return this.getMaxAttId(this.message);
  }

  getMaxAttId(msg: Message): number {
    let a = 1;
    msg.Parts.forEach(p => a += this.getMaxAttId(p))
    return a;
  }
}
