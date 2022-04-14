import {Component, Input} from "@angular/core";
import {EmailViewDisplay} from "../../../definitions/internal";
import { DomSanitizer } from '@angular/platform-browser';

@Component({
  selector: 'app-part-display',
  templateUrl: 'email-part-display.component.html',
  styleUrls: [
    'email-part-display.component.css',
    '../../../app.component.css',
  ],
})
export class EmailPartDisplayComponent {
  @Input() part?: EmailViewDisplay;

  imageMimeTypes = ['image/jpeg'];

  constructor(
    private readonly sanitizer: DomSanitizer,
  ) {
  }

  getAsImage() {
    let image;
    const reader = new FileReader();
    reader.onload = (e) => image = e?.target?.result;
    const data = this.part?.message.Body?.Attachment?.Data || new Uint8Array([]);
    reader.readAsDataURL(new Blob([data]));
    return image;
  }
}
