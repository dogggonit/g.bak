import {Pipe, PipeTransform} from "@angular/core";
import {Email, Message} from "../definitions/email";

@Pipe({name: 'emailMime'})
export class EmailMimePipe implements PipeTransform {
  mimeType = '';

  transform = (email: Email, mimeType = this.mimeType): string => {
    return email.Message ? this.getBody(email.Message, mimeType) : '';
  }

  getBody = (msg: Message, mimeType: string): string => {
    if (mimeType === msg.MimeType) {
      const body = msg.Body?.Data || '';
      if (body !== '') {
        return decodeURIComponent(escape(atob(body.replace(/\-/g, '+').replace(/\_/g, '/'))));
      }
    }

    for (const p of msg.Parts) {
      const body = this.getBody(p, mimeType);
      if (body !== '') {
        return body;
      }
    }

    return '';
  }
}
